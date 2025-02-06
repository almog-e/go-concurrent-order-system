package main

import (
	"fmt"
	"time"
)

func StartZoneManager(zoneName string, zoneChannel chan Order, workers int, doneZone chan bool, displayManager *DisplayManager) {
	tokenPool := make(chan struct{}, workers)

	go func() {
		for order := range zoneChannel {
			tokenPool <- struct{}{}
			go func(o Order) {
				fmt.Printf("[%s] Processing Order: %+v\n", zoneName, o)
				time.Sleep(100 * time.Millisecond)
				displayManager.AddOrder(o)
				<-tokenPool
			}(order)
		}
		doneZone <- true
	}()
}
