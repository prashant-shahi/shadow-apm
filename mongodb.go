package main

import (
	"errors"
	"log"
	"strconv"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var db *mgo.Database

func (m *ApmDAO) Connect() (error) {
	session, err := mgo.Dial(m.Server)
	if err != nil {
		return err
	}
	/*cred := mgo.Credential{
		Username: DBUSER,
		Password: DBPWD,
	}
	err = session.Login(&cred)
	if err != nil {
		return err
	}*/
	db = session.DB(m.Database)
	return nil
}

func (m *ApmDAO) FindAll(query bson.M) ([]MongoObject, error) {
	var mo []MongoObject
	err := db.C(COLLECTION).Find(query).All(&mo)
	return mo, err
}

func (m *ApmDAO) FindById(id string) (MongoObject, error) {
	var mo MongoObject
	err := db.C(COLLECTION).FindId(bson.ObjectIdHex(id)).One(&mo)
	return mo, err
}

func (m *ApmDAO) FindByQuery(query bson.M) ([]MongoObject, error) {
	var mo []MongoObject
	/*tempJSON, err := json.Marshal(queries[0][""])
	if err != nil {
		return nil, err
	}
	err = db.C(COLLECTION).Find(bson.M{string(tempJSON)}).All(&mo)*/
	err := db.C(COLLECTION).Find(query).All(&mo)
	return mo, err
}

func (m *ApmDAO) FindDistinct(field string, query bson.M) ([]interface{}, error) {
	var fullarray []interface{}
	err := db.C(COLLECTION).Find(query).Distinct(field, &fullarray)
	return fullarray, err
}

func (m *ApmDAO) Insert(mo MongoObject) error {
	err := db.C(COLLECTION).Insert(&mo)
	return err
}

// Add some data
func (m *ApmDAO) BulkInsert(mos []interface{}) error {
	log.Output(0, "Function: BulkInsert [ MongoDB handler function ]")
	if len(mos) <= 0 {
		return errors.New("No MongoObjects found")
	}
	bulk := db.C(COLLECTION).Bulk()
	/*var mongoobjects []interface{}
	for _, element  := range mo{
	    mongoobjects = append(mongoobjects, element)
	}*/
	bulk.Insert(mos...)
	bulkresult, err := bulk.Run()
	if err != nil {
		return err
	}
	log.Output(0, "Bulk Result:\n Matched:"+strconv.Itoa(bulkresult.Matched)+"\t Modified:"+strconv.Itoa(bulkresult.Modified))
	return nil
}

func (m *ApmDAO) BulkDelete(mo []MongoObject) error {
	return nil
}

func (m *ApmDAO) Delete(mo MongoObject) error {
	err := db.C(COLLECTION).Remove(&mo)
	return err
}

func (m *ApmDAO) Update(mo MongoObject) error {
	err := db.C(COLLECTION).UpdateId(mo.ID, &mo)
	return err
}