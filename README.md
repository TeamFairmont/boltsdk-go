# boltsdk-go
Go SDK for use when creating workers for Bolt engine

Usage:
go get github.com/TeamFairmont/boltsdk-go

Code example:
    package main

    import (
        "fmt"
        "os"
        "net/http"
    	"io/ioutil"

        "github.com/TeamFairmont/boltsdk-go/boltsdk"
    	"github.com/TeamFairmont/boltshared/mqwrapper"
    	"github.com/TeamFairmont/gabs"
    )

    func main() {
        boltsdk.EnableLogOutput(true)

        //connect to mq server
    	mq, err := mqwrapper.ConnectMQ("amqp://guest:guest@localhost:5672/")
    	if err != nil {
    		fmt.Println(err)
    		os.Exit(1)
    	}

        boltsdk.RunWorker(mq, "fetchPage", fetchPage)

        log.Println("Worker setup complete, waiting for commands... ")

    	//let our worker above process forever or until Ctrl+C
    	forever := make(chan bool)
    	<-forever
    }

    func fetchPage(payload *gabs.Container) error {
    	//request the http url
    	url := payload.Path("initial_input.url").Data().(string)
    	payload.SetP(url, "return_value.url")

    	resp, err := http.Get(url)
    	if err != nil {
    		payload.SetP("Error getting url content: "+err.Error(), "error.fetchPage")
    		return errors.New("Error getting url content: " + err.Error())
    	}

    	defer resp.Body.Close()
    	body, err := ioutil.ReadAll(resp.Body)
    	if err != nil {
    		payload.SetP("Error getting url content: "+err.Error(), "error.fetchPage")
    		return errors.New("Error getting url content: " + err.Error())
    	}
    	payload.SetP(string(body), "return_value.html")

    	return nil
    }
