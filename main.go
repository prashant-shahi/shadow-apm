package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/rs/cors"
)

type App struct {
	Router *mux.Router
}

// Initialize app
func (a *App) Initialize() {
	a.Router = mux.NewRouter()
	a.Router.StrictSlash(true)
	a.setRouters()
}
func (a *App) setRouters() {
	// gets all the services
	// $ curl -X GET http://localhost:8200/services
	//	a.Get("/services", a.getServices)

	// gets all urls for a specific service
	// $ curl -X GET http://localhost:8200/service/my-app-01
	//	a.Get("/service/{servicename}", a.getServiceUrls)

	// gets all request objects for a specific url under a service
	// $ curl -X GET http://localhost:8200/service/my-app-01/requests
	//	a.Get("/service/{servicename}/requests", a.getServiceRequests)

	// default url where all events from all client agents are posted
	// $ curl -X POST http://localhost:8200/intake/v2/events
	a.Post("/intake/v2/events", a.getEvents)
}

func main() {
	router := new(App)
	router.Initialize()
	router.run("0.0.0.0:8200")
}

func (a *App) run(host string) {
	log.Fatal(http.ListenAndServe(host, a.Router))
}
