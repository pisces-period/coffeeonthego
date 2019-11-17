package model

import "gopkg.in/mgo.v2/bson"

//Coffee type - exported struct and fields
type Coffee struct {
	ID               bson.ObjectId `json:"id" bson:"_id"`
	Flavor           string        `json:"flavor"`
	CoffeeMessage    string        `json:"coffee-message"`
	PreparationState string        `json:"preparationState"`
}
