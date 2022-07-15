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

func getWsSessionId(host, playerId string) (string, error) {
	var (
		getSessionPath = fmt.Sprintf("http://%s/api/v1/players/%s/ws", host, playerId)
		resp           *http.Response
		err            error
	)

	//t1 := time.Now()
	for {
		if resp == nil || resp.StatusCode == http.StatusBadGateway {
			resp, err = http.Get(getSessionPath)
			if err != nil {
				return "", err
			}
			if resp.StatusCode == http.StatusOK {
				break
			}
			log.Println("the cat doesn't want to play with you, retry login!")
			time.Sleep(5 * time.Second)
		}
		//if time.Since(t1).Minutes() > 10 {
		//	return "", errors.New("the cat doesn't want to play with you")
		//}
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
		return "", errors.New("maximum number of sessions has exceeded, plz try close browser and run again")
	}
	return body["data"].(map[string]interface{})["id"].(string), nil
}
