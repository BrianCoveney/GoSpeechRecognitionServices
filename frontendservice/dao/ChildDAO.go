package dao

import (
	. "github.com/BrianCoveney/GoSpeechRecognitionServices/frontendservice/models"
	"github.com/globalsign/mgo"
	"gopkg.in/mgo.v2/bson"
	"log"
)

type ChildDAO struct {
	Server   string
	Database string
}

var db *mgo.Database

const (
	COLLECTION = "children"
)

// Establish a connection to database
func (m *ChildDAO) Connect() {
	session, err := mgo.Dial(m.Server)
	if err != nil {
		log.Fatal(err)
	}
	db = session.DB(m.Database)
}

// Find list of children
func (c *ChildDAO) FindAll() ([]Child, error) {
	var child []Child
	err := db.C(COLLECTION).Find(bson.M{}).All(&child)
	return child, err
}

// TODO Find child by email address
func (c *ChildDAO) FindByEmail(email string) ([]Child, error) {
	child := Child{}
	err := db.C(COLLECTION).Find(bson.M{"email": email}).One(&child)
	if err != nil {
		log.Printf("FindByEmail : ERROR :d %s\n", err)
	}

	var ch []Child
	ch = append(ch, child)
	return ch, err
}
