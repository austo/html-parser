package db

import (
	"flag"
	"github.com/austo/html-parser/config"
	"github.com/austo/html-parser/nrsv"
	"testing"
)

var (
	cfgs = map[string]config.Mongo{
		"local": config.Mongo{
			Url:        "mongodb://localhost:27017/nrsv",
			Collection: "verses",
		},
	}
	vr = nrsv.VerseRecord{42, "Luke", 23, 51,
		"had not agreed to their plan and action. He came from the Jewish town of Arimathea, and he was waiting expectantly for the kingdom of God."}
	env = flag.String("e", "local", "environment (e.g. dev, qa, prod, local)")
)

func TestConnectToDatabase(t *testing.T) {
	db := getDatabase(t)
	t.Log(db)
	db.Close()
}

func TestVerseInsert(t *testing.T) {
	db := getDatabase(t)
	err := db.InsertVerseRecord(vr)
	if err != nil {
		t.Fatal(err)
	}
	err = db.DeleteVerseRecord(vr)
	if err != nil {
		t.Fatal(err)
	}
}

func getDatabase(t *testing.T) *Db {
	flag.Parse()
	if _, ok := cfgs[*env]; !ok {
		t.Fatalf("environment \"%s\" is not supported for mongodb", *env)
	}
	db, err := Connect(cfgs[*env])
	if err != nil {
		t.Fatal("db is nil")
	}
	return db
}
