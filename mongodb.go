package main

import (
	"log"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type ApmDAO struct {
	Server   string
	Database string
}

var db *mgo.Database

const (
	TRANSACTION = "transactions"
	METADATA    = "metadata"
)

type MongoObject struct {
	TraceID   string `json:"trace_id"`
	Timestamp int64  `json:"timestamp"`
	Result    string `json:"result"`
	Metadata  struct {
		Service struct {
			Name string `json:"name"`
		} `json:"service"`
		Version  interface{} `json:"version"`
		Language struct {
			Version string `json:"version"`
			Name    string `json:"name"`
		} `json:"language"`
		Agent struct {
			Version string `json:"version"`
			Name    string `json:"name"`
		} `json:"agent"`
		Framework struct {
			Version string `json:"version"`
			Name    string `json:"name"`
		} `json:"framework"`
	} `json:"metadata"`
	Request struct {
		URL     string `json:"url"`
		Body    string `json:"body"`
		Headers struct {
			ContentLength string `json:"content-length"`
			ContentType   string `json:"content-type"`
			Host          string `json:"host"`
			Accept        string `json:"accept"`
			UserAgent     string `json:"user-agent"`
		} `json:"headers"`
		Method string `json:"method"`
	} `json:"request"`
	Response struct {
		StatusCode int `json:"status_code"`
		Headers    struct {
			ContentLength string `json:"Content-Length"`
			ContentType   string `json:"Content-Type"`
		} `json:"headers"`
	} `json:"response"`
	Duration float64 `json:"duration"`
}

func (m *ApmDAO) Connect() {
	session, err := mgo.Dial(m.Server)
	if err != nil {
		log.Fatal(err)
	}
	db = session.DB(m.Database)
}

func (m *ApmDAO) FindAll() ([]MongoObject, error) {
	var movies []MongoObject
	err := db.C(COLLECTION).Find(bson.M{}).All(&movies)
	return movies, err
}

/*func (m *ApmDAO) FindById(id string) (Movie, error) {
	var movie Movie
	err := db.C(COLLECTION).FindId(bson.ObjectIdHex(id)).One(&movie)
	return movie, err
}

func (m *ApmDAO) Insert(movie Movie) error {
	err := db.C(COLLECTION).Insert(&movie)
	return err
}

func (m *ApmDAO) Delete(movie Movie) error {
	err := db.C(COLLECTION).Remove(&movie)
	return err
}

func (m *ApmDAO) Update(movie Movie) error {
	err := db.C(COLLECTION).UpdateId(movie.ID, &movie)
	return err
}*/
