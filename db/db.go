package db

import (
	"fmt"
	"github.com/austo/html-parser/config"
	"gopkg.in/mgo.v2"
	// "gopkg.in/mgo.v2/bson"
)

type Db struct {
	session *mgo.Session
	cfg     config.Mongo
}

func (db *Db) Connect() (err error) {
	session, err := mgo.Dial(db.cfg.Url)
	if err != nil {
		return
	}
	db.session = session
	return
}

func (db *Db) Close() {
	if db.session == nil {
		return
	}
	db.session.Close()
}

func (db *Db) isConnected() error {
	if db.session == nil {
		return fmt.Errorf("not connected")
	}
	return nil
}

func NewDB(cfg config.Mongo) *Db {
	db := new(Db)
	db.cfg = cfg
	return db
}

func Connect(cfg config.Mongo) (db *Db, err error) {
	db = NewDB(cfg)
	err = db.Connect()
	return
}
