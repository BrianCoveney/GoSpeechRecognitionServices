package models

type (
	Child struct {
		FirstName  string         `bson:"first_name" json:"first_name"`
		SecondName string         `bson:"second_name" json:"second_name"`
		Email      string         `bsodocker n:"email" json:"email"`
		Word       string         `bson:"word" json:"word"`
		Words      map[string]int `bson:"map_of_gliding_words" json:"map_of_gliding_words"`
		Errors     map[string]string
	}
)
