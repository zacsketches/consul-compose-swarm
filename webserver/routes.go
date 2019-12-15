package main

import (
	"bytes"
	"context"
	"html/template"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const tpl = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <meta http-equiv="X-UA-Compatible" content="ie=edge">
    <title>Demo Product</title>
</head>
<body>
    <h1>Demo Product</h1>
    This page loads a product from a mongo database and renders the name of
    the awesome product and its description on this page.
    <h3>{{ .Name }}</h3>
    <b>Description: </b>{{ .Description }}
</body>
</html>
`

func (s *server) routes() {
	s.router.Handle("/product", handlers.LoggingHandler(os.Stdout, s.handleProduct()))
	s.router.Handle("/", handlers.LoggingHandler(os.Stdout, s.handleRoot()))
}

func (s *server) handleProduct() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var products []product
		ctx := context.TODO()
		coll := s.client.Database(s.database).Collection(s.collection)
		fo := options.Find()
		cur, err := coll.Find(ctx, bson.D{{}}, fo)
		if err != nil {
			http.Error(w, "handleProduct find: "+err.Error(), http.StatusInternalServerError)
			return
		}
		defer cur.Close(ctx)

		for cur.Next(ctx) {
			var p product
			err := cur.Decode(&p)
			if err != nil {
				http.Error(w, "handleProduct decode: "+err.Error(), http.StatusInternalServerError)
				return
			}

			products = append(products, p)
		}

		// If there are no products then return an error
		if len(products) < 1 {
			http.Error(w, "handleProduct: no products in the database", http.StatusInternalServerError)
			return
		}

		if err := cur.Err(); err != nil {
			http.Error(w, "handleProduct cursor err: "+err.Error(), http.StatusInternalServerError)
			return
		}

		t, err := template.New("product_page").Parse(tpl)
		if err != nil {
			http.Error(w, "handleProduct template parse: "+err.Error(), http.StatusInternalServerError)
			return
		}
		pageBuffer := new(bytes.Buffer)
		err = t.Execute(pageBuffer, products[0])
		if err != nil {
			http.Error(w, "handleProduct template execute: "+err.Error(), http.StatusInternalServerError)
			return
		}
		pageBuffer.WriteTo(w)
	}
}

func (s *server) handleRoot() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		p := defaultProduct()
		t, err := template.New("product_page").Parse(tpl)
		if err != nil {
			http.Error(w, "unable to parse template: "+err.Error(), http.StatusInternalServerError)
			return
		}

		pageBuffer := new(bytes.Buffer)
		err = t.Execute(pageBuffer, p)
		if err != nil {
			http.Error(w, "unable to execute template: "+err.Error(), http.StatusInternalServerError)
			return
		}
		pageBuffer.WriteTo(w)
	}
}
