package main

import (
	"encoding/json"
	_ "fmt"
	_"os"
	"log"
	"net/http"
	_ "strconv"
	_ "reflect"
	"io"
	"io/ioutil"
	_ "net/http/httputil"
	"compress/gzip"
	"strings"

	/*"github.com/gorilla/mux"*/
)

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
type ListRequests struct {
	Status      string `json:"status"`
	ServiceName string `json:"service_name"`
	URL         string `json:"url"`
	Methods     []struct {
		Timestamp int64  `json:"timestamp"`
		Result    string `json:"result"`
		Body      string `json:"body"`
		Headers   struct {
			ContentType string `json:"Content-Type"`
		} `json:"headers"`
		Method string `json:"method"`
	} `json:"methods"`
}

/*a.Get("/services", a.getServices)
a.Get("/service/{servicename}", a.getServiceUrls)
a.Get("/service/{servicename}/requests", a.getServiceRequests)*/

func (a *App) getEvents(w http.ResponseWriter, r *http.Request) {
	log.Output(0, "Function: getEvents [ HTTP handler function ]")
	contentEncoding := r.Header.Get("Content-Encoding")
	contentType := r.Header.Get("Content-Type")
	log.Output(0, "contentType: "+contentType)
	log.Output(0, "contentEncoding: "+contentEncoding)
	var reader io.Reader
	switch contentEncoding {
		case "gzip":
			var err error
			log.Output(0, "gzip: Reached here")
			reader, err = gzip.NewReader(r.Body)
			if err != nil {
		    	log.Output(0, "Error: "+err.Error())
		    	return
		    }
		default:
			reader = r.Body
	}
	allBody, err := copyToString(reader)
	if err != nil {
		log.Output(0, "Error: "+err.Error())
		return
	}
	transactions := getTransactions(allBody)
	log.Output(0, "Reached end of getEvents")
	//fmt.Print(reader.Read(x))
	//log.Output(0, "reflect.TypeOf(reader): "+reflect.TypeOf(reader).String())
	/*log.Output(0, "body: "+string(reader))
	log.Output(0, "reflect.TypeOf(body): "+reflect.TypeOf(reader).String())*/
}

func copyToString(r io.Reader) (res string, err error) {
	log.Output(0, "Function: copyToString")
	b, err := ioutil.ReadAll(r)
	if err != nil {
		log.Fatal(err)
		return "", err
	}
	return string(b), nil
}

func getTransactions(allBody string) []interface{} {
	log.Output(0, "Function: getTransaction")
	requestObjects := strings.Split(allBody, "\n")
	var message context 
	var transactions []interface{}
	for _, element := range requestObjects {
		if element == "" {
			/*log.Println("Empty string at key ",key)*/
			continue
		}
		err := json.Unmarshal([]byte(element), &message)
		if err != nil {
			log.Fatal(err)
			return
		}
		if val, ok := message["transaction"]; ok {
			transactions = append(transactions, val)
			/*log.Println("element:\n",element)*/
			/*log.Println("val:\n",val)
			transactionJSON, err := json.Marshal(val)
			if err != nil {
				log.Fatal(err)
				return
			}
			log.Println("string(val):\n",string(transactionJSON))*/
		}
	}
	/*for key, element := range transactions {
		log.Println("key : ",key)
		log.Println("element : ",element)
	}*/
	return transactions
}