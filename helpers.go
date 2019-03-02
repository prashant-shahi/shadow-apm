package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"errors"
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
		return http.StatusOK, errors.New("No Transactions found")
	}
	err := dao.BulkInsert(transactions)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	return http.StatusCreated, nil
}