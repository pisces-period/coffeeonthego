package controller

import (
	"os"

	"gopkg.in/mgo.v2"
)

//GetSession returns a reference to a MongoDB session
func GetSession() *mgo.Session {

	s, err := mgo.Dial("mongodb://" +
		os.Getenv("DB_USER") + ":" +
		os.Getenv("DB_PASS") + "@" +
		os.Getenv("DB_CONTAINER"))

	if err != nil {
		panic(err)
	}
	return s
}
