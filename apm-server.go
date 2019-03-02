package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	_"strconv"

	"gopkg.in/mgo.v2/bson"
	"github.com/mitchellh/mapstructure"
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
	Result    string `json:"result"`
	Body      string `json:"body"`
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

/*a.Get("/services", a.getServices)
a.Get("/service/{servicename}", a.getServiceUrls)
a.Get("/service/{servicename}/requests", a.getServiceRequests)*/

func (a *App) getServices(w http.ResponseWriter, r *http.Request) {
	log.Output(0, "Function: getServices [ HTTP handler function ]")
	services, err := dao.FindDistinct("metadata.service.name", nil)
	if err != nil {
		log.Fatal(0, "Error while fetching services.\t"+err.Error())
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	/*
	allTransactions, err := dao.FindAll()
	if err != nil {
		log.Fatal(0, "Error while fetching transactions.\t"+err.Error())
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	services := []string{}
	for _, mo := range allTransactions {
		services = AppendIfUnique(services, mo.Metadata.Service.Name)
	}
	*/
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
		log.Fatal(0, "Error while fetching urls.\t"+err.Error())
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	/*allTransactions, err := dao.FindAll()
	if err != nil {
		log.Fatal(0, "Error while fetching transactions.\t"+err.Error())
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	urls := []string{}
	for _, mo := range allTransactions {
		if mo.Metadata.Service.Name == serviceName {
			urls = AppendIfUnique(urls, mo.Request.URL)
		}
	}*/
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
		log.Output(0, "Error while Unmarshalling bodyJSON.\t"+err.Error())
		respondWithError(w, http.StatusInternalServerError, "Invalid request body")
		return
	}
	url := bodyUrl["url"]
	//allTransactions, err := dao.FindAll()
	/*queryStr := `{"metadata.service.name":"`+serviceName+`", "request.url": "`+url+`"}`*/
	/*queryStr := `{"request.url": "`+url+`"}`*/
	/*queryStr := `{"metadata.service.name":"`+serviceName+`"}`
	log.Output(0, "queryStr:\n"+queryStr)*/
	allTransactions, err := dao.FindAll(bson.M{"$and": []bson.M{ {"metadata.service.name": serviceName }, { "request.url": url }}})
	if err != nil {
		log.Fatal(0, "Error while fetching transactions.\t"+err.Error())
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	log.Println(allTransactions)
	lr := ListRequests{
		Status :		"success",
		ServiceName :	serviceName,
		URL :			url,
	}
	var methods []Method
	for _, mo := range allTransactions {
		if mo.Metadata.Service.Name == serviceName && mo.Request.URL == url {
			method := Method{
				Timestamp :		mo.Timestamp,
				Result :		mo.Result,
				Duration :		mo.Duration,
				Body :			mo.Request.Body,
				Headers :		mo.Request.Headers,
				Method :		mo.Request.Method,
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

func AppendIfUnique(slice []string, i string) []string {
    for _, ele := range slice {
        if ele == i {
            return slice
        }
    }
    return append(slice, i)
}

func (a *App) getEvents(w http.ResponseWriter, r *http.Request) {
	log.Output(0, "Function: getEvents [ HTTP handler function ]")
	contentEncoding := r.Header.Get("Content-Encoding")
	/*contentType := r.Header.Get("Content-Type")*/
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
		log.Fatal(0, "Error while inserting transaction.\t"+err.Error())
		return
	}
	/*movie.ID = bson.NewObjectId()
	if err := dao.Insert(movie); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJson(w, http.StatusCreated, movie)*/
	if statusCode != http.StatusCreated {
		log.Output(0, "No-200 statusCode")
		return
	}
	log.Output(0, "Successfully reached the end of getEvents")
}

func getTransactions(allBody string) []interface{} {
	log.Output(0, "Function: getTransaction")
	requestObjects := strings.Split(allBody, "\n")
	var message context
	var transactions []interface{}
	var m Metadata
	for _, element := range requestObjects {
		if element == "" {
			continue
		}
		log.Println("element:\n",element)
		err := json.Unmarshal([]byte(element), &message)
		if err != nil {
			log.Fatal(err)
			return nil
		}
		if val, ok := message["metadata"]; ok {
			log.Output(0, "Metadata detected")
			mapstructure.Decode(val, &m)
			/*metadataJSON, err := json.Marshal(m)
			if err != nil {
				log.Fatal(err)
				return nil
			}
			log.Println("string(metadataJSON):\n", string(metadataJSON))*/
		}
		if val, ok := message["transaction"]; ok {
			var t Transaction
			log.Output(0, "Transaction detected")
			/*tempJSON, err := json.Marshal(val)
			if err != nil {
				log.Fatal(err)
				return nil
			}
			log.Println(string(tempJSON))*/
			/*mapstructure.Decode(val, &t)*/
			tempJSON, err := json.Marshal(val)
			if err != nil {
				log.Fatal(err)
				return nil
			}
			err = json.Unmarshal([]byte(tempJSON), &t)
			if err != nil {
				log.Fatal(err)
				return nil
			}
			mo := MongoObject{
				ID: bson.NewObjectId(),
				Timestamp: t.Timestamp,
				Sampled:    t.Sampled,
				Result:    t.Result,
				Duration:  t.Duration,
				TraceID:   t.TraceID,
			}

			mo.Metadata.Service.Name = m.Service.Name
			mo.Metadata.Version = m.Service.Version
			mo.Metadata.Language = m.Service.Language
			mo.Metadata.Agent = m.Service.Agent
			mo.Metadata.Framework = m.Service.Framework
			mo.Request.URL = t.Context.Request.URL.Full
			mo.Request.Body = t.Context.Request.Body
			mo.Request.Headers = t.Context.Request.Headers
			mo.Request.Method = t.Context.Request.Method
			mo.Response.StatusCode = t.Context.Response.StatusCode
			mo.Response.Headers = t.Context.Response.Headers
			transactions = append(transactions, &mo)
		}
	}
	transactionsJSON, err := json.Marshal(transactions)
	if err != nil {
		log.Fatal(err)
		return nil
	}
	log.Println("string(transactionsJSON):\n", string(transactionsJSON))
	return transactions
}