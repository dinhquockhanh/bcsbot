package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

var (
	errMaxConnection = errors.New("maximum number of sessions has exceeded, plz try close browser and run again")
)

func getWsSessionId(host, playerId, token string) (string, error) {
	var (
		getSessionPath = fmt.Sprintf("http://%s/api/v1/players/%s/ws?token=%s", host, playerId, token)
		resp           *http.Response
		err            error
	)

	for {
		if resp == nil || resp.StatusCode == http.StatusBadGateway {
			resp, err = http.Get(getSessionPath)
			if err != nil {
				return "", err
			}
			if resp.StatusCode == http.StatusOK {
				break
			}
			log.Printf("the cat doesn't want to play with %s, retry login!", playerId)
			time.Sleep(5 * time.Second)
		}
	}

	resBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	body := map[string]interface{}{}
	if err := json.Unmarshal(resBody, &body); err != nil {
		return "", err
	}
	if body["code"] != "0000" {
		return "", errMaxConnection
	}
	return body["data"].(map[string]interface{})["id"].(string), nil
}
