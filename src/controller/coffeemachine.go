package controller

import (
	"encoding/json"
	"fmt"
	"html/template"
	"model"
	"net/http"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var t *template.Template

//CoffeeMachine struct declares a session field that references a mongo session
type CoffeeMachine struct {
	session *mgo.Session
}

//init inflates templates and assigns them to 't'
func init() {
	t = template.Must(template.ParseGlob("/src/view/templates/*.gohtml"))
}

//NewCoffeeMachine takes a reference to a mongo session
//and returns a pointer to a CoffeeMachine struct
func NewCoffeeMachine(s *mgo.Session) *CoffeeMachine {
	return &CoffeeMachine{s}
}

//GetIndex inflates the index.gohtml template
func (cm CoffeeMachine) GetIndex(w http.ResponseWriter, r *http.Request) {
	err := t.ExecuteTemplate(w, "index.gohtml", nil)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError) // 500
		return
	}

	w.WriteHeader(http.StatusOK) // 200
}

//ImATeapot inflates the teapot.gohtml template
func (cm CoffeeMachine) ImATeapot(w http.ResponseWriter, r *http.Request) {
	err := t.ExecuteTemplate(w, "teapot.gohtml", nil)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError) // 500
		return
	}

	w.WriteHeader(http.StatusTeapot) // 418
}

//GetAllCoffees returns all coffee object IDs
func (cm CoffeeMachine) GetAllCoffees(w http.ResponseWriter, r *http.Request) {
	c := []model.Coffee{}

	if err := cm.session.DB("coffeeshop").C("coffee").Find(bson.M{}).All(&c); err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError) // 500
		return
	}

	err := t.ExecuteTemplate(w, "coffeequeue.gohtml", c)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError) // 500
		return
	}

	w.WriteHeader(http.StatusOK) // 200
}

//GetCoffee inspects the GET request parameters and returns a coffee object, if found
func (cm CoffeeMachine) GetCoffee(w http.ResponseWriter, r *http.Request) {
	params, ok := r.URL.Query()["id"]

	// if 'id' parameter is not set to the request
	if !ok {
		cm.GetAllCoffees(w, r) // get all coffees
		return                 // do not execute the rest of the code
	}

	id := params[0] // HTTP Request params are maps of slices

	// if 'id' parameter is not valid
	if !bson.IsObjectIdHex(id) {
		err := t.ExecuteTemplate(w, "coffeeerror.gohtml", http.StatusBadRequest) // 400
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError) // 500
			return
		}
		w.WriteHeader(http.StatusBadRequest) // 400
		return
	}

	// converting 'id' into a BSON object
	oid := bson.ObjectIdHex(id)

	c := model.Coffee{}

	// if BSON representation of 'id' is not found in the database
	if err := cm.session.DB("coffeeshop").C("coffee").FindId(oid).One(&c); err != nil {
		err = t.ExecuteTemplate(w, "coffeeerror.gohtml", http.StatusNotFound) // 404
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError) // 500
			return
		}
		w.WriteHeader(http.StatusNotFound) // 404
		return
	}

	err := t.ExecuteTemplate(w, "coffee.gohtml", c)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError) // 500
		return
	}

}

//BrewCoffee creates a new coffee object if none exists, or updates an existing coffee object
func (cm CoffeeMachine) BrewCoffee(w http.ResponseWriter, r *http.Request) {
	// HTCPCP accepts both POST and BREW methods - POST is deprecated
	if r.Method == http.MethodPost || r.Method == "BREW" {
		c := model.Coffee{} // create variable of type Coffee
		c.Flavor = "traditional"

		json.NewDecoder(r.Body).Decode(&c) // converts JSON to Go coffee

		//TODO: if message says 'stop' the object should exist
		switch c.CoffeeMessage {
		case "start":
			c.ID = bson.NewObjectId()        // assign a bson ID for the coffee variable
			c.PreparationState = "preparing" // start preparing the coffee
			cm.session.DB("coffeeshop").C("coffee").Insert(c)
			w.WriteHeader(http.StatusCreated) // 201

		case "stop":

			// if BSON representation of 'id' is not found in the database
			if err := cm.session.DB("coffeeshop").C("coffee").FindId(c.ID).One(&c); err != nil {
				w.WriteHeader(http.StatusNotFound) // 404
				fmt.Fprintf(w, "ccould not find this id in the database: %s\nHTTP Response: %v", c.ID, http.StatusNotFound)
				return
			}

			c.PreparationState = "ready" // coffee is ready (?)
			cm.session.DB("coffeeshop").C("coffee").UpdateId(c.ID, &c)
			w.WriteHeader(http.StatusOK) // 200
		}

		cj, _ := json.Marshal(c) // converts Go coffee to JSON

		fmt.Fprintf(w, "%s\n", cj) // write to HTTP Response
	}

}
