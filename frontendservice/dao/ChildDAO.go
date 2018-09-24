package dao

import (
	. "github.com/BrianCoveney/GoSpeechRecognitionServices/frontendservice/models"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"log"
)

type ChildDAO struct {
	Server   string
	Database string
}

var db *mgo.Database
var session *mgo.Session

const (
	COLLECTION = "children"
)

// Establish a connection to database
func (c *ChildDAO) Connect() {
	session, err := mgo.Dial(c.Server)
	if err != nil {
		log.Fatal(err)
	}
	db = session.DB(c.Database)
}

// Find list of children
func (c *ChildDAO) FindAll() ([]Child, error) {
	var child []Child
	err := db.C(COLLECTION).Find(bson.M{}).All(&child)
	return child, err
}

// Find child by email
func (c *ChildDAO) FindByEmail(email string) (Child, error) {
	child := Child{}
	err := db.C(COLLECTION).Find(bson.M{"email": email}).All(&child)
	return child, err
}

// Find child by name
func (c *ChildDAO) FindByName(name string) (Child, error) {
	var child Child
	err := db.C(COLLECTION).Find(bson.M{"first_name": name}).One(&child)
	log.Print("CHILD", child)
	return child, err
}
