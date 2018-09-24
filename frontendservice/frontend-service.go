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
var dao = ChildDAO{}

const (
	dev = false
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

	router.HandleFunc("/", indexHandler).Methods("GET")
	router.HandleFunc("/contact", contactHandler).Methods("GET")
	router.HandleFunc("/{email}", searchHandler).Methods("GET")
	router.HandleFunc("/{name}", searchNameHandler).Methods("GET")

	fs := http.FileServer(http.Dir("./static"))
	router.PathPrefix("/images/").Handler(fs)
	router.PathPrefix("/css/").Handler(fs)

	return router
}

// Handler "/"
func indexHandler(w http.ResponseWriter, r *http.Request) {
	c, _ := dao.FindAll()
	index.Render(w, c)
}

// Handler for "/contact"
func contactHandler(w http.ResponseWriter, r *http.Request) {
	c, _ := dao.FindAll()
	contact.Render(w, c)
}

// Handler for path: "/{email}"
func searchHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	c, _ := dao.FindByEmail(vars["email"])
	index.Render(w, c)
}

// Handler for path: "/{email}"
func searchNameHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	c, _ := dao.FindByName(vars["first_name"])
	index.Render(w, c)
}
