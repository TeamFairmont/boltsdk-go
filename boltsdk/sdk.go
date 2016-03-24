package boltsdk

import (
	"log"

	"github.com/TeamFairmont/boltshared/mqwrapper"
	"github.com/TeamFairmont/boltshared/validation"
	"github.com/TeamFairmont/gabs"
)

var enableLogOutput bool

// EnableLogOutput enables (true) or disables (false) the SDK's output
// to stdout for logging / debugging purposes. Defaults to false
func EnableLogOutput(output bool) {
	enableLogOutput = output
}

// WorkerFunc is the user-function type used to process work
// Needs to take a gabs Payload container
type WorkerFunc func(*gabs.Container) error

// RunWorker sets up an AMQP channel, then spins up a worker goroutine
// using the options and work function
func RunWorker(mq *mqwrapper.Connection, commandName string, wf WorkerFunc) error {
	logOut("Starting worker for", commandName)

	//set base QoS parms
	ch, _ := mq.Connection.Channel()

	//connect to the proper mq consumer for this command
	q, res, err := mqwrapper.CreateConsumeNamedQueue(commandName, ch)
	_ = q
	if err != nil {
		return err
	}

	// spin up the goroutine to process work
	go func() {
		for d := range res {
			logOut(commandName, "in")

			//grab the message body and parse to json obj
			payload, err := gabs.ParseJSON(d.Body)
			if err != nil {
				logOut("err:", commandName, err)
			} else {
				//run work func
				err := wf(payload)
				if err != nil {
					logOut("err:", commandName, err)
					pushError(commandName, err)
				}

				//validate payload structure
				err = validate.CheckPayloadStructure(payload)
				if err != nil {
					logOut("err:", commandName, err)
					pushError(commandName, err)
				}

				//push our response to the temp mq replyTo path
				err = mqwrapper.PublishCommand(ch, d.CorrelationId, d.ReplyTo, payload, "")
				if err != nil {
					logOut("err:", commandName, err)
					pushError(commandName, err)
				}

			}

			d.Ack(false) //tell mq we've handled the message
			logOut(commandName, "out")
		}
	}()

	logOut("Started worker for", commandName)
	return nil
}

func logOut(v ...interface{}) {
	if enableLogOutput {
		log.Println(v...)
	}
}

func pushError(commandName string, err error) {
	//TODO
}
