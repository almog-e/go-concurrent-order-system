package main

import (
	"context"
	"fmt"
	"math/rand"
	"time"
)

func StartProducer(restaurant Restaurant, ch chan string, ctx context.Context) {
	fmt.Printf("Restaurant %d started\n", restaurant.ID)
	foodNames := []string{"Pizza", "Burger", "Sushi"}
	ordersLeft := restaurant.Orders
	m := map[string]int{
		"Pizza":  0,
		"Burger": 0,
		"Sushi":  0,
	}

	for ordersLeft > 0 {
		select {
		case <-ctx.Done():
			fmt.Printf("Restaurant %d received cancel signal, exiting...\n", restaurant.ID)
			return
		default:
			// Simulate producing orders
			foodIdx := rand.Intn(len(foodNames)) // Pick a random food type
			order := fmt.Sprintf("Restaurant %d: %s %d", restaurant.ID, foodNames[foodIdx], m[foodNames[foodIdx]])
			fmt.Printf("%s\n", order)
			m[foodNames[foodIdx]]++
			time.Sleep(500 * time.Millisecond)
			ch <- order
			ordersLeft--
		}
	}

	// Send DONE signal as a string
	ch <- "DONE"
	fmt.Printf("Restaurant %d completed all orders.\n", restaurant.ID)
}
