package main

import (
	"log"
	"time"
)

var oneSecond, _ = time.ParseDuration("1s")
var callsPerHourLimit = float64(4500)
var startTime = time.Now()
var totalCalls = 0

func makeGitHubCall(operation func() error) error {
	totalCalls++
	for {
		avgTime := time.Since(startTime).Seconds() / float64(totalCalls)
		if 3600/avgTime > callsPerHourLimit {
			if debug {
				log.Printf("Lowering speed of calls (rate of %d calls per hour)", int(3600/avgTime))
			}
			time.Sleep(oneSecond)
		} else {
			break
		}
	}

	// Retry a few times if error occurs
	var err error
	retries := 0
	for {
		err = operation()
		if err == nil {
			break
		}
		if debug {
			log.Printf("Retrying github call: %s", err)
		}
		retries++
		if retries > 5 {
			break
		}
	}

	return err
}
