package main

import (
	"fmt"
	"time"
	tr "github.com/BrianCoveney/GoSpeechRecognition/transport"
	"labix.org/v2/mgo/bson"
	"log"
	"gopkg.in/mgo.v2"
	"sync"
	"github.com/nats-io/nats"
	"os"
)

//const MongoDb details
const (
	hosts      = "ec2-54-202-69-181.us-west-2.compute.amazonaws.com:8080"
	database   = "speech"
	username   = "speechUser"
	password   = "bossdog12"
	collection = "children"
)

type (

	// BuoyStation contains information for an individual station.
	Child struct {
		ID          bson.ObjectId `bson:"_id,omitempty"`
		FirstName   string        `bson:"first_name"`
		SecondName  string        `bson:"second_name"`
		Email       string        `bson:"email"`

	}
)

var nc *nats.Conn

func main() {

	// NATS
	uri := os.Getenv("NATS_URI")

	var err error

	nc, err = nats.Connect(uri)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Connected to NATS " + uri)


	// MongoDB
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
	//fmt.Println(fmt.Sprintf("Messages count: %d", count))

	// Create a wait group to manage the goroutines.
	var waitGroup sync.WaitGroup

	// Perform 10 concurrent queries against the databaseservice.
	waitGroup.Add(1)
	for query := 0; query < 1; query++ {
		go RunQuery(query, &waitGroup, session)
	}

	// Wait for all the queries to complete.
	waitGroup.Wait()
	log.Println("All Queries Completed")
}

func RunQuery(query int, waitGroup *sync.WaitGroup, mongoSession *mgo.Session) {
	// Decrement the wait group count so the program knows this
	// has been completed once the goroutine exits.
	defer waitGroup.Done()

	// Request a socket connection from the session to process our query.
	// Close the session when the goroutine exits and put the connection back
	// into the pool.
	sessionCopy := mongoSession.Copy()
	defer sessionCopy.Close()

	// Get a collection to execute the query against.
	collection := sessionCopy.DB(database).C(collection)

	//log.Printf("RunQuery : %d : Executing\n", query)

	// Retrieve the list of stations.
	//var child []Child

	child := tr.ChildUser{}
	err := collection.Find(nil).All(&child)
	if err != nil {
		log.Printf("RunQuery : ERROR : %s\n", err)
		return
	}


	//err := collection.Find(nil).All(&child)
	//if err != nil {
	//	log.Printf("RunQuery : ERROR : %s\n", err)
	//	return
	//}

	fmt.Println("Phone", child)
	//log.Printf("RunQuery : %d : Count[%d]\n", query, len(child))
}