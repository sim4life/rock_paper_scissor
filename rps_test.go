package main

import (
	"log"
	"io/ioutil"
	"encoding/json"
	"testing"
)

func TestJson(t *testing.T) {
	bytesData, err := ioutil.ReadFile("rps.json")
    if err != nil {
		log.Fatal(err)
    }
	
	var gh GameHistory
	if err := json.Unmarshal([]byte(bytesData), &gh); err != nil {
		t.Errorf("Can not parse Game History JSON.")
		log.Fatal(err)
	}

	ghLength := len(gh)
    if ghLength != 100 {
       t.Errorf("Game History length was incorrect, got: %d, want: %d.", ghLength, 100)
	}
}