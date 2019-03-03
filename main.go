package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

var dao = ApmDAO{}

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
	a.Get("/services", a.getServices)

	// gets all urls for a specific service
	// $ curl -X GET http://localhost:8200/service/my-app-01
	a.Get("/service/{servicename}", a.getServiceUrls)

	// gets all request objects for a specific url under a service
	// $ curl -X GET http://localhost:8200/service/my-app-01/requests -H "Content-Type: application/json" -d '{ "url": "http://localhost:5000/test" }'
	a.Post("/service/{servicename}/requests", a.getServiceRequests)

	// gets the request body to simulate the request from the server and respond back
	// $ curl -X POST "http://localhost:8200/playit" -H "Content-Type: application/json" -d '{ "trace_id": "83f0865e1c76d3da71ed5bde2274365c", "service_name": "my-app-01" }'
	a.Post("/playit", a.simulateRequest)

	// default url where all events from all client agents are posted
	// $ curl -X POST http://localhost:8200/intake/v2/events
	a.Post("/intake/v2/events", a.getEvents)
}

func init() {
	dao.Server = DBUSER+":"+DBPWD+"@"+DB_SERVER_URL+":"+strconv.Itoa(DB_SERVER_PORT)+"/"+DATABASE
	dao.Database = DATABASE
	err := dao.Connect()
	if err != nil {
		log.Fatal("MongoDB Error: "+err.Error())
	}
	log.Output(0, "Connection to database successful")
}

func main() {
	log.Output(0, "Starting the server in a while")
	router := new(App)
	router.Initialize()
	log.Output(0, "Server successfully running at http://0.0.0.0:8200")
	router.run("0.0.0.0:8200")
}

func (a *App) run(host string) {
	handler := cors.Default().Handler(a.Router)
	log.Fatal(http.ListenAndServe(host, handler))
}