package main

import (
	"fmt"
	"time"
	"labix.org/v2/mgo/bson"
	"log"
	"gopkg.in/mgo.v2"
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

var err error

func main() {
	server := &http.Server{
		Addr:    ":3001",
		Handler: initRoutes(),
	}
	server.ListenAndServe()
}

func initRoutes() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/", printToScreen)
	router.HandleFunc("/{email}", findChildByEmail).Methods("GET")
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

func findAllChildren() []Child {
	sessionCopy := getMongoSession().Copy()

	// Get our collection
	collection := sessionCopy.DB(database).C(collection)

	var children []Child
	// Run query on collection to find all. Our struct holds the result.
	err := collection.Find(bson.M{}).All(&children)
	if err != nil {
		log.Printf("findAllChildren : ERROR : %s\n", err)
	}
	return children
}

func findChildByEmail(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	child := Child{Email: vars["email"]}

	sessionCopy := getMongoSession().Copy()
	collection := sessionCopy.DB(database).C(collection)

	result := Child{}
	err = collection.Find(bson.M{"email": child.Email}).One(&result)
	if err != nil {
		log.Printf("findChildByEmail : ERROR : %s\n", err)
	}

	var c []Child
	c = append(c, result)

	t, _ := template.ParseFiles("view.html")
	t.Execute(w, c)
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
	var children = findAllChildren()
	var c []Child
	for _, child := range children {
		c = append(c, child)
	}
	t.Execute(w, c)
}
