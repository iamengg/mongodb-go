// go get github.com/gorilla/mux
// go get github.com/mongodb/mongo-go-driver
// https://www.youtube.com/watch?v=oW7PMHEYiSk&list=WL&index=12
// Developing a RESTful API with Golang and a MongoDB NoSQL Database

// Other project is
//
// https://github.com/qiangxue/go-rest-api
// Youtube mongodb with golang & rest, & postman
// https://www.youtube.com/watch?v=oW7PMHEYiSk&list=WL&index=14
// postman tutorial
// https://www.guru99.com/postman-tutorial.html#4

// postman
// https://martian-eclipse-414323.postman.co/workspace/myworkspace~30b3ac36-98d0-4de1-bb06-1e0aca36cadf/request/create?requestId=2c06b2c5-d6ca-472b-ae57-c0d690d98ee5

// docker & kubernets https://www.youtube.com/channel/UCHN0KlJD8pePI83VuH_d8sw
// running docerized app with cloud (aws, azure, gcp)

// http://localhost:12345/people
// http://localhost:12345/people
/*
// POST



http://localhost:12345/person
{
    "firstname":"pratik",
    "lastname":"shitole"
}
GET
http://localhost:12345/person
GET
http://localhost:12345/people
//
//Sample mongodb pages
https://www.mongodb.com/blog/post/mongodb-go-driver-tutorial-part-1-connecting-using-bson-and-crud-operations

*/
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//connect to mongodb
//define data model
//define endpoints

//json for webbrowser
//bson for mongodb (marshal & unmarshal)

type Person struct {
	ID        primitive.ObjectID `json:"_id, omitempty", bson:"_id, omitempty"`
	Firstname string             `json:"firstname, omitempty", bson:"firstname, omitempty"`
	Lastname  string             `json:"lastname, omitempty", bson:"lastname, omitempty"`
}

var client *mongo.Client

const (
	atlas_uri = "mongodb+srv://pratiktest:standard@cluster0.f653n.mongodb.net/myFirstDatabase?retryWrites=true&w=majority"
	local_uri = "mongodb://localhost:27017"
)

func CreatePersonEndpoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("content-type", "application/json")
	var person Person
	json.NewDecoder(request.Body).Decode(&person)
	collection := client.Database("thepolyglotdeveloper").Collection("people")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	result, _ := collection.InsertOne(ctx, person)
	json.NewEncoder(response).Encode(result)
}

func GetPeopleEndpoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("content-type", "application/json")

	var people []Person
	json.NewDecoder(request.Body).Decode(&people)
	collection := client.Database("thepolyglotdeveloper").Collection("people")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var person Person
		cursor.Decode(&person)
		people = append(people, person)
	}
	if err := cursor.Err(); err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		return
	}
	json.NewEncoder(response).Encode(people)
}

func GetPersonEndpoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("content-type", "application/json")
	params := mux.Vars(request)
	id, _ := primitive.ObjectIDFromHex(params["id"])

	var person Person
	json.NewDecoder(request.Body).Decode(&person)
	collection := client.Database("thepolyglotdeveloper").Collection("people")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	result := collection.FindOne(ctx, Person{ID: id}).Decode(&person)
	json.NewEncoder(response).Encode(result)
}

func GetPersonBasedOnName(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("content-type", "application/json")
	collection := client.Database("thepolyglotdeveloper").Collection("people")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	var people []Person

	cursor, err := collection.Find(ctx, bson.D{{"firstname", bson.D{{"$in", bson.A{"pratik"}}}}})
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var person Person
		cursor.Decode(&person)
		people = append(people, person)
	}
	if err := cursor.Err(); err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		return
	}
	json.NewEncoder(response).Encode(people)
}

func main() {
	fmt.Println("Starting the applications ...")

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, _ = mongo.Connect(ctx, options.Client().ApplyURI(atlas_uri))
	router := mux.NewRouter()

	router.HandleFunc("/person", CreatePersonEndpoint).Methods("POST")
	router.HandleFunc("/people", GetPeopleEndpoint).Methods("GET")
	router.HandleFunc("/person", GetPersonEndpoint).Methods("GET")
	router.HandleFunc("/personByName", GetPersonBasedOnName).Methods("GET")
	http.ListenAndServe(":12345", router)
	fmt.Println("Closing the application")
}
