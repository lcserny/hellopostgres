package hellopostgres

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"path"
)

func Run(cfgPath string) error {
	dbCfg, err := initDbConfig(cfgPath)
	if err != nil {
		return err
	}

	httpCfg, err := initHttpConfig(cfgPath)
	if err != nil {
		return err
	}

	greetingResource := NewGreetingResource(*dbCfg, *httpCfg)
	defer greetingResource.Close()
	return greetingResource.Init()
}

func initDbConfig(cfgPath string) (*dbConfig, error) {
	dbJsonBytes, err := ioutil.ReadFile(path.Join(cfgPath, "db.json"))
	if err != nil {
		return nil, errors.New("please setup a 'db.json' file in the specified configs folder")
	}

	dbCfg := dbConfig{}
	if err := json.Unmarshal(dbJsonBytes, &dbCfg); err != nil {
		return nil, errors.New("invalid db.json file specified")
	}

	return &dbCfg, nil
}

func initHttpConfig(cfgPath string) (*httpConfig, error) {
	httpJsonBytes, err := ioutil.ReadFile(path.Join(cfgPath, "http.json"))
	if err != nil {
		return nil, errors.New("please setup a 'http.json' file in the specified configs folder")
	}

	httpCfg := httpConfig{}
	if err := json.Unmarshal(httpJsonBytes, &httpCfg); err != nil {
		return nil, errors.New("invalid http.json file specified")
	}

	return &httpCfg, nil
}
