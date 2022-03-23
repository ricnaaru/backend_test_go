package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func GetMap(r *http.Request) map[string]interface{} {
	var payload interface{}
	buffer, err := ioutil.ReadAll(r.Body)

	if err != nil {
		fmt.Sprintf("error1 %d", err)
		log.Panic(err)
		return nil
	}

	err2 := r.Body.Close()
	if err2 != nil {
		fmt.Sprintf("error2 %d", err)
		log.Panic(err2)
		return nil
	}

	err3 := json.Unmarshal(buffer, &payload)
	if err3 != nil {
		fmt.Sprintf("error3 %d", err)
		log.Panic(err3)
		return nil
	}

	m := payload.(map[string]interface{})

	return m
}

func CastString(o interface{}) string {
	if o == nil {
		return ""
	}

	return o.(string)
}
