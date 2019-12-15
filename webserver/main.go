package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

const (
	serverPort    = "8080"
	version       = "v0.1.0"
	mongoHostname = "db"
)

type server struct {
	router     *http.ServeMux
	client     *mongo.Client
	database   string
	collection string
}

func (s *server) Connect() error {
	mongoURI := fmt.Sprintf("mongodb://%s:27017", mongoHostname)
	clientOpts := options.Client().ApplyURI(mongoURI)
	client, err := mongo.Connect(context.TODO(), clientOpts)
	if err != nil {
		return fmt.Errorf("unable to connect: %s", err)
	}

	// Call Ping to verify that the deployment is up and the Client was configured successfully.
	// As mentioned in the Ping documentation, this reduces application resiliency as the server may be
	// temporarily unavailable when Ping is called.
	if err = client.Ping(context.TODO(), readpref.Primary()); err != nil {
		log.Fatal(err)
	}
	s.client = client

	// Assing a default database
	s.database = "test"
	s.collection = "default"

	return nil
}

func (s *server) Init() error {
	p2 := product{
		Name:        "Widget",
		Description: "The cool thing that does that thing.",
	}
	bson := bson.M{
		"name":        p2.Name,
		"description": p2.Description,
	}
	_, err := s.client.Database(s.database).Collection(s.collection).InsertOne(
		context.TODO(), bson)
	return err
}

func main() {

	srv := server{router: http.DefaultServeMux}

	err := srv.Connect()
	if err != nil {
		panic(err)
	}
	log.Printf("server successfully connected and pinged\n")

	err = srv.Init()
	if err != nil {
		panic(err)
	}

	srv.routes()

	log.Printf("Launching server on port %s\n", serverPort)
	err = http.ListenAndServe(":"+serverPort, nil)
	if err != nil {
		panic(err)
	}
}
