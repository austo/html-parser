package db

import (
	"flag"
	"github.com/austo/html-parser/config"
	"testing"
)

var (
	cfgs = map[string]config.Mongo{
		"local": config.Mongo{
			Url:        "mongodb://localhost:27017/nrsv",
			Collection: "verses",
		},
	}
	env = flag.String("e", "local", "environment (e.g. dev, qa, prod, local)")
)

func TestConnectToDatabase(t *testing.T) {
	db := getDatabase(t)
	t.Log(db)
	db.Close()
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
