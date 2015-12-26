package core

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type GDConf struct {
	Url string `json:"url"`
}

func ReadGDConf(path string) (*GDConf, error) {
	fmt.Println("loading configuration ", path)
	fileCnt, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	conf := new(GDConf)
	err = json.Unmarshal(fileCnt, conf)
	if err != nil {
		return nil, err
	} else {
		return conf, nil
	}
}
