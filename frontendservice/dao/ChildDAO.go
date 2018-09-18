package dao

import (
	. "github.com/BrianCoveney/GoSpeechRecognitionServices/frontendservice/models"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
)

type ChildDAO struct {
	Server		string
	Database	string
}

var db *mgo.Database

const (
	COLLECTION = "children"
)

// Connect to the database
func (c * ChildDAO) Connect()  {
	session, err := mgo.Dial(c.Server)
	if err != nil {
		log.Fatal(err)
	}
	db = session.DB(c.Database)
}

// Find list of children
func (c *ChildDAO) FindAll() ([]Child, error) {
	var children []Child
	err := db.C(COLLECTION).Find(bson.M{}).All(&children)
	return children, err
}

// Find child by email address
func (c *ChildDAO) FindByEmail(email string) (Child, error) {
	var child Child
	err := db.C(COLLECTION).Find(bson.M{"email": child.Email}).One(&child)
	return child, err
}