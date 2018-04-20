package config

import (
	"encoding/json"
	"io/ioutil"

	log "github.com/inconshreveable/log15"
)

var (
	jsonObj map[string]interface{}
)

func LoadConfigFile(filename string) error {
	var err error = nil
	jsonObj, err = readFile(filename)

	if err != nil {
		log.Error("Read Json Error")
	}

	return err
}

func readFile(filename string) (j map[string]interface{}, e error) {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Error("ReadFile: ", err.Error())
		return nil, err
	}

	if err := json.Unmarshal(bytes, &j); err != nil {
		log.Error("Unmarshal: ", err.Error())
		return nil, err
	}

	return j, nil
}

func GetStringValue(key string, defVal string) string {
	if v, ok := jsonObj[key]; ok {
		return v.(string)
	}
	return defVal
}

func GetValue(key string, defVal interface{}) interface{} {
	if v, ok := jsonObj[key]; ok {
		//fmt.Println(reflect.TypeOf(v))
		return v
	}
	return defVal
}
