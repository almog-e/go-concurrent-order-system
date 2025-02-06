package main

import (
	"context"
	"fmt"
	"strings"
	"time"
)

func Dispatcher(producerChannels map[int]chan string, pizzaChannel, burgerChannel, sushiChannel chan string, ctx context.Context) {
	fmt.Println("Dispatcher started")
	activeProducers := len(producerChannels)

	for activeProducers > 0 {
		select {
		case <-ctx.Done():
			fmt.Println("Dispatcher received cancel signal, exiting...")
			return
		default:
			for id, ch := range producerChannels {
				select {
				case order := <-ch:
					if order == "DONE" {
						fmt.Printf("Dispatcher: Producer %d is done.\n", id)
						delete(producerChannels, id)
						activeProducers--
						continue
					}
					// Route order to correct zone
					time.Sleep(500 * time.Millisecond)
					if strings.Contains(order, "Pizza") {
						pizzaChannel <- order
					} else if strings.Contains(order, "Burger") {
						burgerChannel <- order
					} else if strings.Contains(order, "Sushi") {
						sushiChannel <- order
					}
				default:
					// No new messages from this producer
				}
			}
		}
	}

	// Send DONE
	pizzaChannel <- "DONE"
	burgerChannel <- "DONE"
	sushiChannel <- "DONE"
	fmt.Println("Dispatcher: All producers are done, exiting.")
}
