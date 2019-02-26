package main

import (
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

func (m *ApmDAO) FindAll() ([]MongoObject, error) {
	var mo []MongoObject
	err := db.C(COLLECTION).Find(bson.M{}).All(&mo)
	return mo, err
}

func (m *ApmDAO) FindById(id string) (MongoObject, error) {
	var mo MongoObject
	err := db.C(COLLECTION).FindId(bson.ObjectIdHex(id)).One(&mo)
	return mo, err
}

func (m *ApmDAO) FindByQuery(query string) ([]MongoObject, error) {
	var mo []MongoObject
	/*tempJSON, err := json.Marshal(queries[0][""])
	if err != nil {
		return nil, err
	}
	err = db.C(COLLECTION).Find(bson.M{string(tempJSON)}).All(&mo)*/
	err := db.C(COLLECTION).Find(bson.M{}).All(&mo)
	return mo, err
}

func (m *ApmDAO) FindDistinct(query string) ([]MongoObject, error) {
	var mo []MongoObject
	/*tempJSON, err := json.Marshal(queries[0][""])
	if err != nil {
		return nil, err
	}
	err = db.C(COLLECTION).Find(bson.M{string(tempJSON)}).All(&mo)*/
	err := db.C(COLLECTION).Find(nil).Distinct(query, &mo)
	return mo, err
}

func (m *ApmDAO) Insert(mo MongoObject) error {
	err := db.C(COLLECTION).Insert(&mo)
	return err
}

func (m *ApmDAO) Delete(mo MongoObject) error {
	err := db.C(COLLECTION).Remove(&mo)
	return err
}

func (m *ApmDAO) Update(mo MongoObject) error {
	err := db.C(COLLECTION).UpdateId(mo.ID, &mo)
	return err
}