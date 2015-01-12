package config

import (
	"encoding/json"
	"io/ioutil"
)

const (
	Filename = "config.json"
)

type Mongo struct {
	Url        string
	Collection string
}

type Configs map[string]Mongo

type ByteReader interface {
	Read() ([]byte, error)
}

type fileReader string

func (r fileReader) Read() ([]byte, error) {
	return ioutil.ReadFile(string(r))
}

func ReadEnvironmentFromFile(filename, env string) (Mongo, error) {
	fr := fileReader(filename)
	return readConfig(fr, env)
}

func ReadEnvironmentsFromFile(
	filename string, deleteLocal bool) (Configs, error) {

	fr := fileReader(filename)
	cfgs, err := readConfigs(fr)
	if err != nil {
		return nil, err
	}
	if deleteLocal {
		delete(cfgs, "local")
	}
	return cfgs, nil
}

func readConfig(r ByteReader, env string) (mongo Mongo, err error) {
	cfgs, err := readConfigs(r)
	if err != nil {
		return
	}
	mongo = cfgs[env]
	return
}

func readConfigs(r ByteReader) (cfgs Configs, err error) {
	bytes, err := r.Read()
	if err != nil {
		return
	}
	err = json.Unmarshal(bytes, &cfgs)
	return
}
