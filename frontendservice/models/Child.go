package models

type (
	Child struct {
		FirstName  string         `bson:"first_name"`
		SecondName string         `bson:"second_name"`
		Email      string         `bson:"email"`
		Word       string         `bson:"word"`
		Words      map[string]int `bson:"map_of_gliding_words"`
	}
)