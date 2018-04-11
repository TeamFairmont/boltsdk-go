package boltsdkFunctions

import (
	"encoding/json"

	"github.com/TeamFairmont/bolt-app-sdk"
	"github.com/TeamFairmont/boltshared/config"
)

var boltConfigChan = make(chan config.Config, 1)

/*
boltURL := "http://localhost:8888/get-config"
user := "searchgroup"
pass := "open_to_public"
*/

//GetBoltConfig takesadsf  adsf
func GetBoltConfig(boltURL, user, pass string) (config.Config, error) {
	err := boltAppSdk.RunApp(boltURL+`/get-config`, user, pass, GetBoltConfigAppFunc)
	if err != nil {
		return config.Config{}, err
	}
	return <-boltConfigChan, nil
}

//GetBoltConfigAppFunc asdf
func GetBoltConfigAppFunc(p map[string]interface{}, payloadChan chan map[string]interface{}, respBodyChan chan []byte, doneChan chan bool, args []interface{}) error {
	emptyPayload := make(map[string]interface{}, 1)
	payloadChan <- emptyPayload
	respBytes := <-respBodyChan
	var boltConfig = config.Config{}
	err := json.Unmarshal(respBytes, &boltConfig)
	if err != nil {
		return err
	}
	boltConfigChan <- boltConfig

	return nil
}
