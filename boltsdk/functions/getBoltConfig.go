package boltsdkFunctions

import (
	"encoding/json"

	"github.com/TeamFairmont/bolt-app-sdk"
	"github.com/TeamFairmont/boltshared/config"
)

var boltConfigChan = make(chan config.Config, 1)

//GetBoltConfig takes the address to the bolt instance, a username and password, and will return the config
func GetBoltConfig(boltURL, user, pass string) (config.Config, error) {
	//call /get-config, a built in bolt call
	err := boltAppSdk.RunApp(boltURL+`/get-config`, user, pass, GetBoltConfigAppFunc)
	if err != nil {
		return config.Config{}, err
	}
	return <-boltConfigChan, nil
}

//GetBoltConfigAppFunc sends an empty payload to /get-config, and returns a bolt config via the boltConfigChan, received in GetBoltConfig
func GetBoltConfigAppFunc(p map[string]interface{}, payloadChan chan map[string]interface{}, respBodyChan chan []byte, doneChan chan bool, args []interface{}) error {
	//ready empty payload to be send
	emptyPayload := make(map[string]interface{}, 1)
	//send payload
	payloadChan <- emptyPayload
	//receive response from bolt
	respBytes := <-respBodyChan
	//unmarshal the response into a bolt config struct
	var boltConfig = config.Config{}
	err := json.Unmarshal(respBytes, &boltConfig)
	if err != nil {
		return err
	} //send the config to GetBoltConfig()
	boltConfigChan <- boltConfig

	return nil
}
