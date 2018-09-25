package main

import (
	"crypto/tls"
	"fmt"
	. "github.com/BrianCoveney/GoSpeechRecognitionServices/frontendservice/dao"
	"github.com/BrianCoveney/GoSpeechRecognitionServices/views"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/acme/autocert"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

var index *views.View
var contact *views.View
var search *views.View
var dao = ChildDAO{}

const (
	dev = false // Or false for production
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
		certManager := autocert.Manager{
			Prompt:     autocert.AcceptTOS,
			HostPolicy: autocert.HostWhitelist("speech.briancoveney.com"),
			Cache:      autocert.DirCache("certs"),
		}

		server := &http.Server{
			Addr:    ":https",
			Handler: initRoutes(),
			TLSConfig: &tls.Config{
				GetCertificate: certManager.GetCertificate,
			},
		}

		go http.ListenAndServe(":http", certManager.HTTPHandler(nil))

		log.Fatal(server.ListenAndServeTLS("", ""))
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
	search = views.NewView("bootstrap", "static/search.gohtml")


	router.HandleFunc("/", indexHandler).Methods("GET")
	router.HandleFunc("/contact", contactHandler).Methods("GET")
	router.HandleFunc("/search", searchFormHandler).Methods("GET", "POST")
	router.HandleFunc("/{email}", searchURLHandler).Methods("GET")


	fs := http.FileServer(http.Dir("./static"))
	router.PathPrefix("/images/").Handler(fs)
	router.PathPrefix("/css/").Handler(fs)

	return router
}

// Handler "/"
func indexHandler(w http.ResponseWriter, r *http.Request) {
	c, err := dao.FindAll()
	if err != nil {
		log.Printf("indexHandler : ERROR :d %s%v\n", err, c)
	}
	index.Render(w, c)
}

// Handler for path: "/{email}"
func searchURLHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	c, err := dao.FindByEmail(vars["email"])
	if err != nil{
		log.Printf("searchURLHandler : ERROR :d %s%v\n", err, c)
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
		search.Render(w, c)
	}
}

// Handler for "/contact"
func contactHandler(w http.ResponseWriter, r *http.Request) {
	c, _ := dao.FindAll()
	contact.Render(w, c)
}



