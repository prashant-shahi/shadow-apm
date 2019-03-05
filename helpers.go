package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"reflect"
	"strings"

	"gopkg.in/mgo.v2/bson"
	"github.com/mitchellh/mapstructure"
)

func respondWithError(w http.ResponseWriter, code int, msg string) {
	respondWithJson(w, code, map[string]string{"error": msg})
}

func respondWithJson(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}


func readAllString(r io.Reader) (res string, err error) {
	log.Output(0, "Function: readAllString")
	b, err := ioutil.ReadAll(r)
	if err != nil {
		log.Output(0, err.Error())
		return "", err
	}
	return string(b), nil
}

func insertMultipleTransactions(transactions []interface{}) (int, error) {
	log.Output(0, "Function: insertMultipleTransactions")
	if len(transactions) <= 0 {
		return http.StatusOK, errors.New("No Transactions found")
	}
	err := dao.BulkInsert(transactions)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	return http.StatusCreated, nil
}

func httpRequest(url string, method string, Headers interface{}, payload_string string) (response SimulateResponse, statusCode int, err error) {
	log.Output(0, "Function: httpRequest")
	var payload = []byte(payload_string)
	method = strings.ToUpper(method)
	if method != "GET" && method != "POST" {
		log.Output(0, "Error while creating http request.\tReason:"+err.Error())
		return response, http.StatusBadRequest, errors.New("Request Method not supported")
	}
    req, err := http.NewRequest(method, url, bytes.NewBuffer(payload))
    if err != nil {
		log.Output(0, "Error while creating http request.\tReason:"+err.Error())
		return response, http.StatusInternalServerError, err
	}
    req.Header.Set("Content-Type", "application/json")
    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
    	log.Println(reflect.TypeOf(err))
    	log.Println(err)
    	if strings.Contains(err.Error(), "connect: connection refused"){
    		return response, http.StatusServiceUnavailable, err
		}
    	log.Output(0, "Error sending http request.\tReason:"+err.Error())
        return response, http.StatusInternalServerError, err
    }
    defer resp.Body.Close()
    responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Output(0, "Error: "+err.Error())
		return response, http.StatusInternalServerError, err
	}
    response = SimulateResponse{
    	Status 			:	"success",
    	StatusCode 		:	resp.StatusCode,
    	Response 		:	strings.TrimSpace(string(responseBody)),
    	Headers 		:	fetchHeaders(resp.Header),
    }
    return response, http.StatusOK, nil
}

func fetchHeaders(headers http.Header) (interface{}) {
	headerArray := make(map[string]interface{}, len(headers))
	for index, header := range headers {
		/*log.Println(reflect.TypeOf(index));
		log.Println(reflect.TypeOf(header));
		log.Println(index);
		log.Println(header[0]);*/
		headerArray[index] = header[0]
	}
	return headerArray
}

func getTransactions(allBody string) []interface{} {
	log.Output(0, "Function: getTransaction")
	// Trimming leading and trailing whitespaces, tabs and new lines
	allBody = strings.TrimSpace(allBody)
	// making array of each lines
	requestObjects := strings.Split(allBody, "\n")
	var message context
	var transactions []interface{}
	var m Metadata
	for _, element := range requestObjects {
		// in case of an empty array object, we will skip
		if element == "" {
			continue
		}
		/*log.Println("element:\n",element)*/
		err := json.Unmarshal([]byte(element), &message)
		if err != nil {
			log.Output("Error: "+err.Error())
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
				log.Output("Error: "+err.Error())
				return nil
			}
			err = json.Unmarshal([]byte(tempJSON), &t)
			if err != nil {
				log.Output("Error: "+err.Error())
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
	/*transactionsJSON, err := json.Marshal(transactions)
	if err != nil {
		log.Fatal(err)
		return nil
	}
	log.Println("string(transactionsJSON):\n", string(transactionsJSON))*/
	return transactions
}

func AppendIfUnique(slice []string, i string) []string {
    for _, ele := range slice {
        if ele == i {
            return slice
        }
    }
    return append(slice, i)
}