package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"strconv"
	"net/http"
	"errors"
	"reflect"

	"gopkg.in/mgo.v2/bson"
	_"github.com/mitchellh/mapstructure"
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
		log.Fatal(err)
		return "", err
	}
	return string(b), nil
}

func insertMultipleTransactions(transactions []interface{}) (int, error) {
	log.Output(0, "Function: insertMultipleTransactions")
	if len(transactions) <= 0 {
		log.Output(0, "No Transactions found")
		return http.StatusCreated, errors.New("No Transactions found")
	}
	/*if err := json.NewDecoder(r.Body).Decode(&mo); err != nil {
		return http.StatusBadRequest, error.New("Invalid request payload")
	}*/
	/*vals := interface{}*/
	for key, val := range transactions {
		var mo MongoObject
		tempJSON, err := json.Marshal(val)
		if err != nil {
			log.Fatal(err)
			return http.StatusInternalServerError, err
		}
		log.Println(string(tempJSON))
		err = json.Unmarshal([]byte(tempJSON), &mo)
		if err != nil {
			log.Fatal(err)
			return http.StatusInternalServerError, nil
		}
		mo.ID = bson.NewObjectId()
		log.Println(reflect.TypeOf(val))
		log.Println(val)
		/*log.Println(mo.ID)
		log.Println("Before decode")
		log.Println(mo)
		mapstructure.Decode(val, &mo)
		log.Println("After decode")
		log.Println(mo)*/
		if err := dao.Insert(mo); err != nil {
			return http.StatusInternalServerError, err
		}
		log.Output(0, "Key "+strconv.Itoa(key)+"\tTransaction Insertion Success")
	}
	return http.StatusCreated, nil
}