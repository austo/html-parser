package db

import (
	"fmt"
	"github.com/austo/html-parser/config"
	"github.com/austo/html-parser/nrsv"
	"gopkg.in/mgo.v2"
	// "gopkg.in/mgo.v2/bson"
)

type Db struct {
	session *mgo.Session
	cfg     config.Mongo
}

type VerseKey struct {
	BookIndex    uint8
	ChapterIndex uint8
	VerseIndex   uint16
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

func (db *Db) InsertVerseRecord(vr nrsv.VerseRecord) error {
	return db.getCollection().Insert(vr)
}

func (db *Db) InsertVerseRecords(vrs []interface{}) error {
	return db.getCollection().Insert(vrs...)
}

func (db *Db) DeleteVerseRecord(vr nrsv.VerseRecord) (err error) {
	vk := getVerseKey(vr)
	info, err := db.getCollection().RemoveAll(vk)
	if err != nil {
		return
	}
	removed := info.Removed
	s := ""
	if removed > 1 {
		s = "s"
	}
	fmt.Printf("%d document%s removed\n", removed, s)
	return
}

func getVerseKey(vr nrsv.VerseRecord) (vk VerseKey) {
	vk.BookIndex = vr.BookIndex
	vk.ChapterIndex = vr.ChapterIndex
	vk.VerseIndex = vr.VerseIndex
	return
}

func (db *Db) getCollection() *mgo.Collection {
	return db.session.DB("").C(db.cfg.Collection)
}
