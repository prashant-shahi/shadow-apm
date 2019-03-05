package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"gopkg.in/mgo.v2/bson"
	"github.com/gorilla/mux"
)

// constant for context -- transaction and metadata
type context map[string]interface{}

// for the return response of api GET /services
type ListingServices struct {
	Status   string   `json:"status"`
	Services []string `json:"services"`
}

// for the return response of api GET /services/{servicename}
type ListingUrls struct {
	Status      string   `json:"status"`
	ServiceName string   `json:"service_name"`
	Urls        []string `json:"urls"`
}

// for the return response of api GET /service/{servicename}/requests
type Method struct {
	Timestamp int64  `json:"timestamp"`
	TraceID string `json:"trace_id"`
	Result    string `json:"result"`
	Body    string        `bson:"body" json:"body"`
	Headers interface{} `bson:"headers" json:"headers"`
	Method string `json:"method"`
	Duration float64 `json:"duration"`
}
type ListRequests struct {
	Status      string `json:"status"`
	ServiceName string `json:"service_name"`
	URL         string `json:"url"`
	Methods     []Method `json:"methods"`
}

// The request body for which to simulate the request to the server
type SimulateRequest struct {
	TraceID string `json:"trace_id"`
	ServiceName    string `json:"service_name"`
}

// The response body of the simulation of the request from the server
type SimulateResponse struct {
	Status 		string `json:"status"`
	Response 	string `json:"response"`
	StatusCode 	int    `json:"status_code"`
	Headers interface{} `bson:"headers" json:"headers"`
}

func (a *App) getServices(w http.ResponseWriter, r *http.Request) {
	log.Output(0, "Function: getServices [ HTTP handler function ]")
	services, err := dao.FindDistinct("metadata.service.name", nil)
	if err != nil {
		log.Output(0, "Error while fetching services.\tReason:"+err.Error())
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	var allServices []string
	for _, t := range services{
	    allServices = append(allServices, t.(string))
	}
	ls := ListingServices{
		Status: "success",
		Services: allServices,
	}
	log.Output(0, "Successfully reached the end of getServices")
	respondWithJson(w, http.StatusOK, ls)
}

func (a *App) getServiceUrls(w http.ResponseWriter, r *http.Request) {
	log.Output(0, "Function: getServiceUrls [ HTTP handler function ]")
	serviceName, ok := mux.Vars(r)["servicename"]
	if ok == false {
		log.Output(0, "Error: service name not found")
		respondWithError(w, http.StatusNotAcceptable, "Error: invalid request")
		return
	}
	log.Output(0, "Service Name: "+serviceName)
	urls, err := dao.FindDistinct("request.url", bson.M{"metadata.service.name": serviceName})
	if err != nil {
		log.Output(0, "Error while fetching urls.\tReason:"+err.Error())
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	var allUrls []string
	for _, t := range urls{
	    allUrls = append(allUrls, t.(string))
	}
	lu := ListingUrls{
		Status: "success",
		ServiceName: serviceName,
		Urls : allUrls,
	}
	log.Output(0, "Successfully reached the end of getServiceUrls")
	respondWithJson(w, http.StatusOK, lu)
}

func (a *App) getServiceRequests(w http.ResponseWriter, r *http.Request) {
	log.Output(0, "Function: getServiceRequests [ HTTP handler function ]")
	serviceName, ok := mux.Vars(r)["servicename"]
	if ok == false {
		log.Output(0, "Error: service name not found")
		respondWithError(w, http.StatusNotAcceptable, "Error: invalid request")
		return
	}
	log.Output(0, "Service Name: "+serviceName)
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Output(0, "Error: "+err.Error())
		respondWithError(w, http.StatusInternalServerError, "Error while reading request body")
		return
	}
	bodyJSONString := string(body)
	bodyUrl := make(map[string]string)
	err = json.Unmarshal([]byte(bodyJSONString), &bodyUrl)
	if err != nil {
		log.Output(0, "Error while Unmarshalling bodyJSON.\tReason:"+err.Error())
		respondWithError(w, http.StatusInternalServerError, "Invalid request body")
		return
	}
	url := bodyUrl["url"]
	allTransactions, err := dao.FindAll(bson.M{"$and": []bson.M{ {"metadata.service.name": serviceName }, { "request.url": url }}})
	if err != nil {
		log.Output(0, "Error while fetching transactions.\tReason:"+err.Error())
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	log.Println(allTransactions)
	lr := ListRequests{
		Status 		:	"success",
		ServiceName :	serviceName,
		URL :			url,
	}
	var methods []Method
	for _, mo := range allTransactions {
		if mo.Metadata.Service.Name == serviceName && mo.Request.URL == url {
			method := Method{
				TraceID 	:	mo.TraceID,
				Timestamp 	:	mo.Timestamp,
				Result 		:	mo.Result,
				Duration 	:	mo.Duration,
				Body 		:	mo.Request.Body,
				Headers 	:	mo.Request.Headers,
				Method 		:	mo.Request.Method,
			}
			methods = append(methods, method)
		}
	}
	if len(allTransactions) <= 0 {
		lr.Methods = []Method{}
	} else {
		lr.Methods = methods
	}
	log.Output(0, "Successfully reached the end of getServiceRequests")
	respondWithJson(w, http.StatusOK, lr)
}

func (a *App) getEvents(w http.ResponseWriter, r *http.Request) {
	log.Output(0, "Function: getEvents [ HTTP handler function ]")
	contentEncoding := r.Header.Get("Content-Encoding")
	var reader io.Reader
	switch contentEncoding {
	case "gzip":
		gzipReader, err := gzipReader(r.Body)
		if err != nil {
			log.Output(0, "Error: "+err.Error())
			return
		}
		reader = gzipReader
	default:
		reader = r.Body
	}
	allBody, err := readAllString(reader)
	if err != nil {
		log.Output(0, "Error: "+err.Error())
		return
	}
	transactions := getTransactions(allBody)
	statusCode, err := insertMultipleTransactions(transactions)
	if (err != nil && statusCode != http.StatusOK && statusCode != http.StatusCreated) {
		log.Output(0, "Error while inserting transaction.\tReason:"+err.Error())
		return
	}
	if statusCode != http.StatusCreated {
		log.Output(0, "200 statusCode - Insertion was success")
		return
	}
	log.Output(0, "Successfully reached the end of getEvents")
}

func (a *App) simulateRequest(w http.ResponseWriter, r *http.Request) {
	log.Output(0, "Function: simulateRequest [ HTTP handler function ]")
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Output(0, "Error: "+err.Error())
		respondWithError(w, http.StatusInternalServerError, "Error while reading request body")
		return
	}
	bodyJSONString := string(body)
	var requestBody SimulateRequest
	err = json.Unmarshal([]byte(bodyJSONString), &requestBody)
	if err != nil {
		log.Output(0, "Error while Unmarshalling bodyJSON.\tReason: "+err.Error())
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	traceId := requestBody.TraceID
	serviceName := requestBody.ServiceName
	log.Output(0, "Trace ID: "+traceId)
	log.Output(0, "Service Name: "+serviceName)
	mongoObject, err := dao.FindOne(bson.M{"$and": []bson.M{ { "trace_id": traceId }, {"metadata.service.name": serviceName }}})
	if err != nil {
		log.Output(0, "Error while fetching the mongoObject.\tReason: "+err.Error())
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	sr, statusCode, err := httpRequest(mongoObject.Request.URL, mongoObject.Request.Method, mongoObject.Request.Headers, mongoObject.Request.Body)
	if err != nil || statusCode != http.StatusOK  {
		if statusCode == http.StatusServiceUnavailable{
			log.Output(0, "Error while requesting the Service URL.\tReason: "+err.Error())
			respondWithError(w, statusCode, err.Error())
			return
		} else {
			log.Output(0, "Error while simulating the request.\tReason: "+err.Error())
			respondWithError(w, statusCode, err.Error())
			return
		}
	}
	respondWithJson(w, http.StatusOK, sr)
	log.Output(0, "Successfully reached the end of simulateRequest")
}