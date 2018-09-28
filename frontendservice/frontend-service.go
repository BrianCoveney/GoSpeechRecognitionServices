package main

import (
	"encoding/json"
	"fmt"
	. "github.com/BrianCoveney/GoSpeechRecognitionServices/frontendservice/dao"
	. "github.com/BrianCoveney/GoSpeechRecognitionServices/frontendservice/models"
	"github.com/BrianCoveney/GoSpeechRecognitionServices/views"
	"github.com/gorilla/mux"
	"io/ioutil"
	"labix.org/v2/mgo/bson"
	"log"
	"net/http"
	"strings"
)

var index *views.View
var contact *views.View
var search *views.View
var dao = ChildDAO{}

const (
	dev = true // Or false for production
)

func readConfigs() []string {
	myKeysFile, err := ioutil.ReadFile("configs")
	if err != nil {
		fmt.Println("There was a problem with the configs")
	}
	return strings.Split(string(myKeysFile), "\n")
}

// Parse the configuration file 'configs.toml', and establish a connection to DB
func init() {
	config := readConfigs()
	server, database := config[0], config[1]

	dao.Server = server
	dao.Database = database
	dao.Connect()
}

// main() method that starts our http server
func main() {

	if dev {
		server := &http.Server{
			Addr:    ":80",
			Handler: initRoutes(),
		}
		server.ListenAndServe()
	} else {

		// Uncomment when pushing to production
		//certManager := autocert.Manager{
		//	Prompt:     autocert.AcceptTOS,
		//	HostPolicy: autocert.HostWhitelist("speech.briancoveney.com"),
		//	Cache:      autocert.DirCache("certs"),
		//}
		//
		//server := &http.Server{
		//	Addr:    ":https",
		//	Handler: initRoutes(),
		//	TLSConfig: &tls.Config{
		//		GetCertificate: certManager.GetCertificate,
		//	},
		//}
		//
		//go http.ListenAndServe(":http", certManager.HTTPHandler(nil))
		//
		//log.Fatal(server.ListenAndServeTLS("", ""))
	}
}

func initRoutes() *mux.Router {
	router := mux.NewRouter()

	index = views.NewView("bootstrap", "static/index.gohtml")
	contact = views.NewView("bootstrap", "static/contact.gohtml")
	search = views.NewView("bootstrap", "static/search.gohtml")

	router.HandleFunc("/", indexHandler).Methods("GET")
	router.HandleFunc("/contact", contactHandler).Methods("GET")
	router.HandleFunc("/search", searchFormHandler).Methods("GET", "POST")
	router.HandleFunc("/children", searchAllHandler).Methods("GET")
	router.HandleFunc("/{email}", searchURLHandler).Methods("GET")
	router.HandleFunc("/{email}", removeChildHandler).Methods("DELETE")
	router.HandleFunc("/{email}", updateChildHandler).Methods("PUT")
	router.HandleFunc("/{first_name}", createChildHandler).Methods("POST")

	fs := http.FileServer(http.Dir("./static"))
	router.PathPrefix("/images/").Handler(fs)
	router.PathPrefix("/css/").Handler(fs)

	return router
}

/*
 * html template endpoints
 *
*/
// Handler "/"
func indexHandler(w http.ResponseWriter, r *http.Request) {
	c, err := dao.FindAll()
	if err != nil {
		log.Printf("indexHandler : ERROR :d %s%v\n", err, c)
	}
	index.Render(w, c)
}

// Handler for search bar form input
func searchFormHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	if r.Method == "GET" {
		search.Render(w, nil)
	} else {
		r.ParseForm()
		email := r.Form["email"][0]
		c, err := dao.FindByEmail(email)
		if err != nil{
			log.Printf("searchFormHandler : ERROR :d %s\n", err)
		}
		var childSlice []Child
		childSlice = append(childSlice, c)
		search.Render(w, childSlice)
	}
}

// Handler for "/contact"
func contactHandler(w http.ResponseWriter, r *http.Request) {
	c, _ := dao.FindAll()
	contact.Render(w, c)
}

/*
 * JSON Payload endpoints
 *
*/
func createChildHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var child Child
	vars := mux.Vars(r)
	child.ID = bson.NewObjectId()
	child.FirstName = vars["first_name"]
	er := dao.CreateChild(child)
	if er != nil {
		respondWithError(w, http.StatusInternalServerError, "Er: " + er.Error())
		return
	}
	respondWithJSON(w, http.StatusCreated, child)
}

// Handler for path: "/{email}"
func searchURLHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	c, err := dao.FindByEmail(vars["email"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Child Email")
		return
	}
	respondWithJSON(w, http.StatusOK, c)
}

// Handler for path: "/children"
func searchAllHandler(w http.ResponseWriter, r *http.Request) {
	c, err := dao.FindAll()
	if err != nil || c == nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, c)
}

func updateChildHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	c := getChild(w, r)
	err := dao.UpdateChild(c)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
	}
	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

func removeChildHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	c := getChild(w, r)
	err := dao.RemoveChild(c)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
	}
	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func getChild(w http.ResponseWriter, r *http.Request) Child {
	vars := mux.Vars(r)
	child, err := dao.FindByEmail(vars["email"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Child Email")
	}
	return child
}


