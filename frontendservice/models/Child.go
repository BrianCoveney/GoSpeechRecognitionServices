package models

import "regexp"

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

func (c *Child) Validate() bool {

	c.Errors = make(map[string]string)

	re := regexp.MustCompile(".+@.+\\..+")
	matched := re.Match([]byte(c.Email))
	if matched == false {
		c.Errors["email"] = "Please enter a valid email address"
	}

	return len(c.Errors) == 0
}