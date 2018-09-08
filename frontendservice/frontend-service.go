package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2"
	"html/template"
	"labix.org/v2/mgo/bson"
	"log"
	"net/http"
	"time"
)

// Class constants that contain information about of DB
const (
	//hostsProd      = "mongodb-repository:27017"
	hostsDev   = "94.156.189.70:27017"
	database   = "speech"
	username   = ""
	password   = ""
	collection = "children"
)

// main() method that starts our http server
func main() {

	server := &http.Server{
		Addr:    ":443",
		Handler: initRoutes(),
	}

	// We used the certbot client on our server to issue our certificates
	err := server.ListenAndServeTLS("/etc/letsencrypt/live/speech.briancoveney.com/fullchain.pem",
		"/etc/letsencrypt/live/speech.briancoveney.com/privkey.pem")
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

}

// initRoutes() method is handler.
// Here we specify that http://<ip_address>/ is handled by the findAllChildren() method, and
// http://<ip_address>/<name@some_email.com> is handled by the findChildByEmail() method.
func initRoutes() *mux.Router {
	// We create a router that we can pass the request through so that the vars will be added to the context.
	// router.HandleFunc register URL paths and their handlers.
	router := mux.NewRouter()
	router.HandleFunc("/", findAllChildren).Methods("GET")
	router.HandleFunc("/{email}", findChildByEmail).Methods("GET")
	return router
}

// Returns a mongoDB session using the constants as needed. This is used by findAllChildren() and findChildByEmail()
func getMongoSession() *mgo.Session {
	info := &mgo.DialInfo{
		Addrs:    []string{hostsDev},
		Timeout:  60 * time.Second,
		Database: database,
		Username: username,
		Password: password,
	}
	session, err1 := mgo.DialWithInfo(info)
	if err1 != nil {
		panic(err1)
	}
	return session
}

// This struct contains the type of collection we will be receiving from the DB, i.e bson strings and a map
type (
	Child struct {
		FirstName  string         `bson:"first_name"`
		SecondName string         `bson:"second_name"`
		Email      string         `bson:"email"`
		Word       string         `bson:"word"`
		Words      map[string]int `bson:"map_of_gliding_words"`
	}
)

func findAllChildren(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close() //  Close the response body when finished with it

	// Here our sessionCopy is set equal to the session returned from our getMongoSession() method
	sessionCopy := getMongoSession().Copy()
	// Get our collection
	collection := sessionCopy.DB(database).C(collection)

	// Create an array of Child
	var children []Child
	// Run our query
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

	// Generate HTML output
	t, _ := template.ParseFiles("view.html")
	t.Execute(w, c)
}

func findChildByEmail(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	// Here mux.Vars(r) creates a map of route variables. See initRoutes() router.HandleFunc("/{email}", findChildByEmail)
	vars := mux.Vars(r)

	// This utilises our Child struct with the Email field set to the result of the mux.Vars request
	child := Child{ Email: vars["email"] }
	sessionCopy := getMongoSession().Copy()
	collection := sessionCopy.DB(database).C(collection)

	// Create an empty Child struct
	childResult := Child{}

	// We use the mgo MongoDB driver, to search for the child by their email address. The result is stored in childResult
	var err = collection.Find(bson.M{"email": child.Email}).One(&childResult)
	if err != nil {
		log.Printf("findChildByEmail : ERROR :d %s\n", err)
	}

	// Append the bson childResult 'childResult' to our struct 'c'
	var c []Child
	c = append(c, childResult)

	// We use package template (html/template) that implements data-driven templates
	// for generating HTML output safe against code injection.
	t, _ := template.ParseFiles("view.html")

	// Available at, e.g http://speech.local/grace@email.com  (This is for developing locally)
	// Available at, e.g http://94.156.189.70/grace@email.com (This is the ip of my server)r
	t.Execute(w, c)
}
