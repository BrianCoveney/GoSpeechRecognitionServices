package main

import (
	"fmt"
	"github.com/BrianCoveney/GoSpeechRecognitionServices/views"
	"github.com/globalsign/mgo"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/acme/autocert"
	"gopkg.in/mgo.v2/bson"
	"log"
	"net/http"
	"time"
)

const (
	host       = "mongodb-repository:27017"
	//host       = "94.156.189.70:27017"
	database   = "speech"
	username   = ""
	password   = ""
	collection = "children"

	dev = false
)

var index *views.View
var contact *views.View

// main() method that starts our http server
func main() {

	if dev {
		server := &http.Server{
			Addr:    ":80",
			Handler: initRoutes(),
		}
		server.ListenAndServe()
	} else {

		certManager := autocert.Manager{
			Prompt:     autocert.AcceptTOS,
			HostPolicy: autocert.HostWhitelist("speech.briancoveney.com"),
			Cache:      autocert.DirCache("certs"),
		}

		server := &http.Server{
			Addr:    ":https",
			Handler: initRoutes(),
		}

		go http.ListenAndServe(":http", certManager.HTTPHandler(nil))

		log.Fatal(server.ListenAndServeTLS("/home/brian/certs/fullchain.pem",
			"/home/brian/certs/privkey.pem"))
	}

}

// initRoutes() method is handler.
// Here we specify that http://<ip_address>/ is handled by the findAllChildren() method, and
// http://<ip_address>/<name@some_email.com> is handled by the findChildByEmail() method.
func initRoutes() *mux.Router {
	// We create a router that we can pass the request through so that the vars will be added to the context.
	// router.HandleFunc register URL paths and their handlers.
	router := mux.NewRouter()

	index = views.NewView("bootstrap", "static/index.gohtml")
	contact = views.NewView("bootstrap", "static/contact.gohtml")

	router.HandleFunc("/", indexHandler).Methods("GET")
	router.HandleFunc("/contact", contactHandler).Methods("GET")
	router.HandleFunc("/{email}", searchHandler).Methods("GET")

	fs := http.FileServer(http.Dir("./static"))
	router.PathPrefix("/images/").Handler(fs)
	router.PathPrefix("/css/").Handler(fs)

	return router
}

// Returns a mongoDB session using the constants as needed. This is used by findAllChildren() and findChildByEmail()
func getMongoSession() *mgo.Session {
	info := &mgo.DialInfo{
		Addrs:    []string{host},
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

func findAllChildren() []Child {
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
	return c
}

func findChildByEmail(r *http.Request) []Child {
	defer r.Body.Close()

	// Here mux.Vars(r) creates a map of route variables. See initRoutes() router.HandleFunc("/{email}", findChildByEmail)
	vars := mux.Vars(r)

	// This utilises our Child struct with the Email field set to the result of the mux.Vars request
	child := Child{Email: vars["email"]}
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

	return c
}

// Handler "/"
func indexHandler(w http.ResponseWriter, r *http.Request) {
	var c = findAllChildren()
	index.Render(w, c)
}

// Handler for path: "/{email}"
func searchHandler(w http.ResponseWriter, r *http.Request) {
	var c = findChildByEmail(r)
	index.Render(w, c)
}

// Handler for "/contact"
func contactHandler(w http.ResponseWriter, r *http.Request) {
	var c = findAllChildren()
	contact.Render(w, c)
}


