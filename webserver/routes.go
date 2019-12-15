package main

import (
	"bytes"
	"html/template"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
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
	s.router.Handle("/", handlers.LoggingHandler(os.Stdout, s.handleRoot()))
}

func (s *server) handleRoot() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		p := defaultProduct()
		t, err := template.New("product_page").Parse(tpl)
		if err != nil {
			http.Error(w, "unable to parse template: "+err.Error(), http.StatusInternalServerError)
		}

		pageBuffer := new(bytes.Buffer)
		err = t.Execute(pageBuffer, p)
		if err != nil {
			http.Error(w, "unable to execute template: "+err.Error(), http.StatusInternalServerError)
		}
		pageBuffer.WriteTo(w)
	}
}
