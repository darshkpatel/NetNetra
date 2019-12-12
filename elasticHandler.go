package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"sync"
)

// converts to json and send data to elastic
//Refer https://github.com/dvas0004/GolangSimpleDnsSniffer
// ToDo: Use WaitGroups
func sendToElastic(data deserializedPacket, wg *sync.WaitGroup) {
	defer wg.Done()

	var jsonMsg, jsonErr = json.Marshal(data)
	if jsonErr != nil {
		log.Panic(jsonErr)
	}

	// Use flags/configfile to get server parameters
	esIndex := "packetData"
	esDocType := "syslog"
	esServer := "localhost"
	request, reqErr := http.NewRequest("POST", "http://"+esServer+":9200/"+esIndex+"/"+esDocType,
		bytes.NewBuffer(jsonMsg))
	if reqErr != nil {
		panic(reqErr)
	}

	client := &http.Client{}
	resp, elErr := client.Do(request)

	if elErr != nil {
		panic(elErr)
	}

	defer resp.Body.Close()

}
