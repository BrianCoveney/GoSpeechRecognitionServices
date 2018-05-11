package main

import (
	"time"
	"labix.org/v2/mgo/bson"
	"log"
	"gopkg.in/mgo.v2"
	"net/http"
	"github.com/gorilla/mux"
	"html/template"
	"fmt"
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
		FirstName  string         `bson:"first_name"`
		SecondName string         `bson:"second_name"`
		Email      string         `bson:"email"`
		Word       string         `bson:"word"`
		Words      map[string]int `bson:"map_of_gliding_words"`
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
	router.HandleFunc("/", findAllChildren).Methods("GET")
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

func findAllChildren(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	sessionCopy := getMongoSession().Copy()

	// Get our collection
	collection := sessionCopy.DB(database).C(collection)

	// Run our query
	var children []Child
	err := collection.Find(bson.M{}).All(&children)
	if err != nil {
		log.Printf("findAllChildren : ERROR : %s\n", err)
	}

	// Append the bson result 'children' to our struct 'c'
	var c []Child
	for _, child := range children {
		c = append(c, child)
		fmt.Printf("Child: %+v\n", child)
	}

	t, _ := template.ParseFiles("view.html")

	// Available at:  http://localhost:3001/
	t.Execute(w, c)

}

func findChildByEmail(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	vars := mux.Vars(r)

	child := Child{Email: vars["email"]}

	sessionCopy := getMongoSession().Copy()
	collection := sessionCopy.DB(database).C(collection)

	childResult := Child{}
	err = collection.Find(bson.M{"email": child.Email}).One(&childResult)
	if err != nil {
		log.Printf("findChildByEmail : ERROR : %s\n", err)
	}

	// Append the bson childResult 'childResult' to our struct 'c'
	var c []Child
	c = append(c, childResult)

	t, _ := template.ParseFiles("view.html")

	// Available at, e.g:  http://localhost:3001/der@email.com
	t.Execute(w, c)
}

