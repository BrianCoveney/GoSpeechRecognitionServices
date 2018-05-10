package main

import (
	"fmt"
	"time"
	"labix.org/v2/mgo/bson"
	"log"
	"gopkg.in/mgo.v2"
	"github.com/nats-io/nats"
	"os"
	"strings"
	"net/http"
	"github.com/gorilla/mux"
)

const (
	hosts      = "ec2-54-202-69-181.us-west-2.compute.amazonaws.com:8080"
	database   = "speech"
	username   = "speechUser"
	password   = "bossdog12"
	collection = "children"
)

type (
	Child struct {
		FirstName  string `bson:"first_name"`
		SecondName string `bson:"second_name"`
		Email      string `bson:"email"`
		Word       string `bson:"word"`
	}
)

var nc *nats.Conn
var child []Child
var err error

func main() {

	// NATs
	uri := os.Getenv("NATS_URI")
	nc, err = nats.Connect(uri)
	if err != nil {
		fmt.Println(err)
		return
	}

	// mongoDB
	session := getMongoSession()
	go RunQuery(session)

	// http
	server := &http.Server{
		Addr:    ":3001",
		Handler: initRoutes(),
	}
	server.ListenAndServe()
}

func initRoutes() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/", printToScreen)
	return router
}

func getMongoSession() *mgo.Session {
	info := &mgo.DialInfo{
		Addrs:    []string{hosts},
		Timeout:  60 * time.Second,
		Database: database,
		Username: username,
		Password: password,
	}

	session, err1 := mgo.DialWithInfo(info)
	if err1 != nil {
		panic(err1)
	}

	col := session.DB(database).C(collection)

	count, err2 := col.Count()
	if err2 != nil {
		panic(err2)
		log.Println("Error %s %d", err, count)
	}
	return session
}


func RunQuery(mongoSession *mgo.Session) {
	sessionCopy := mongoSession.Copy()

	// Get our collection
	collection := sessionCopy.DB(database).C(collection)

	// Run query on collection to find all. Our struct holds the result.
	err := collection.Find(bson.M{}).All(&child)
	if err != nil {
		log.Printf("RunQuery : ERROR : %s\n", err)
		return
	}

	// Loop through our struct result
	for i, c := range child {
		fmt.Printf("Child %d: %+v\n", i, c)
	}
}


func printToScreen(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()       // parse arguments, you have to call this by yourself
	fmt.Println(r.Form) // print form information in server side
	fmt.Println("path", r.URL.Path)
	fmt.Println("scheme", r.URL.Scheme)
	fmt.Println(r.Form["url_long"])
	for k, v := range r.Form {
		fmt.Println("key:", k)
		fmt.Println("val:", strings.Join(v, ""))
	}

	// Print all child detail to the screen, available at http://localhost:3001/
	for i, c := range child {
		fmt.Fprintf(w, "Child %d: %+v\n", i, c.String())
	}
}

// Our toString for formatting
func (this Child) String() string {
	return this.FirstName + " | " + this.SecondName + " | " + this.Email + " | " + this.Word
}
