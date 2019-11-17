package main

import (
	ctl "controller"
	"net/http"
)

func main() {
	cm := ctl.NewCoffeeMachine(ctl.GetSession())

	http.Handle("/assets/", http.StripPrefix("/assets",
		http.FileServer(http.Dir("/src/view/assets"))))

	http.Handle("/favicon.ico", http.NotFoundHandler())
	http.HandleFunc("/", cm.GetIndex)
	http.HandleFunc("/teapot/", cm.ImATeapot)
	http.HandleFunc("/coffee/get", cm.GetCoffee)
	http.HandleFunc("/coffee/brew", cm.BrewCoffee)

	http.ListenAndServe(":80", nil)
}
