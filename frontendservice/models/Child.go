package models

type (
	Child struct {
		FirstName  string         `bson:"first_name" json:"first_name"`
		SecondName string         `bson:"second_name" json:"first_name"`
		Email      string         `bson:"email" json:"first_name"`
		Word       string         `bson:"word" json:"first_name"`
		Words      map[string]int `bson:"map_of_gliding_words" json:"first_name"`
	}
)