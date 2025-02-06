package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

func ZoneManager(zoneName string, zoneChannel chan string, DisplayChannel chan string, tokenPool chan struct{}, ctx context.Context) {
	fmt.Printf("%s started\n", zoneName)
	var wg sync.WaitGroup

	for {
		select {
		case <-ctx.Done():
			wg.Wait()
			fmt.Printf("%s received cancel signal, exiting...\n", zoneName)
			return
		case order := <-zoneChannel:
			if order == "DONE" {
				wg.Wait()
				fmt.Printf("%s: All orders processed, exiting.\n", zoneName)
				return
			}

			// Acquire a token
			tokenPool <- struct{}{}
			wg.Add(1)
			go func(order string) {
				defer wg.Done()
				fmt.Printf("%s: Processing %s\n", zoneName, order)
				time.Sleep(1000 * time.Millisecond) // Simulate processing time

				var orderWithDetails string = zoneName + " " + order
				DisplayChannel <- orderWithDetails
				<-tokenPool // Release token
				fmt.Printf("%s: Done %s\n", zoneName, order)
			}(order)
		}
	}
}
