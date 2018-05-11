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
	"html/template"
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
var children []Child
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
	err := collection.Find(bson.M{}).All(&children)
	if err != nil {
		log.Printf("RunQuery : ERROR : %s\n", err)
		return
	}

	// Loop through our struct result
	for i, child := range children {
		fmt.Printf("Child %d: %+v\n", i, child)
	}
}

func printToScreen(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	fmt.Println(r.Form)
	fmt.Println("path", r.URL.Path)
	fmt.Println("scheme", r.URL.Scheme)
	fmt.Println(r.Form["url_long"])
	for k, v := range r.Form {
		fmt.Println("key:", k)
		fmt.Println("val:", strings.Join(v, ""))
	}

	// Generate HTML template
	t, _ := template.ParseFiles("view.html")

	// Print all children's details to the screen. This will be available at http://localhost:3001/
	var c []Child
	for _, child := range children {
		c = append(c, child)
	}
	t.Execute(w, c)
}
