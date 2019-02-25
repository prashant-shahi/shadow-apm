package main

import (
	_ "compress/gzip"
	"encoding/json"
	_ "fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	_ "net/http/httputil"
	_ "os"
	_ "reflect"
	_ "strconv"
	"strings"

	"github.com/mitchellh/mapstructure"
	/*"github.com/gorilla/mux"*/)

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
		gzipReader, err := gzipReader(r.Body)
		if err != nil {
			log.Output(0, "Error: "+err.Error())
			return
		}
		reader = gzipReader
	default:
		reader = r.Body
	}
	allBody, err := copyToString(reader)
	if err != nil {
		log.Output(0, "Error: "+err.Error())
		return
	}
	transactions := getTransactions(allBody)
	log.Println(transactions)
	log.Output(0, "Reached end of getEvents")
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
	var t Transaction
	var m Metadata
	for _, element := range requestObjects {
		if element == "" {
			/*log.Println("Empty string at key ",key)*/
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
			/*log.Println(m)*/
			metadataJSON, err := json.Marshal(m)
			if err != nil {
				log.Fatal(err)
				return nil
			}
			log.Println("string(metadataJSON):\n", string(metadataJSON))
		}
		if val, ok := message["transaction"]; ok {
			log.Output(0, "Transaction detected")
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
				Timestamp: t.Timestamp,
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
			/*log.Println(mo)*/
			transactions = append(transactions, mo)
			/*log.Println("element:\n",element)*/
			/*log.Println("val:\n",val)*/
			/*transactionJSON, err := json.Marshal(val)
			if err != nil {
				log.Fatal(err)
				return
			}
			log.Println("string(val):\n",string(transactionJSON))*/
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
