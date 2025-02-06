package main

import (
	"fmt"
	"math/rand"
	"time"
)

// Order struct represents a food order
type Order struct {
	RestaurantID int
	OrderID      int
	FoodType     string
}

// StartProducers spawns goroutines for each restaurant
func StartProducers(restaurants []Restaurant, ordersChannel chan Order, doneChannel chan bool) {
	rand.Seed(time.Now().UnixNano())
	foodNames := []string{"Pizza", "Burger", "Sushi"}

	for _, restaurant := range restaurants {
		go func(r Restaurant) {
			for i := 1; i <= r.Orders; i++ {
				randomIndex := rand.Intn(len(foodNames))
				randomFood := foodNames[randomIndex]

				ordersChannel <- Order{RestaurantID: r.ID, OrderID: i, FoodType: randomFood}

				fmt.Printf("Restaurant %d sent order %d\n", r.ID, i)
			}
			doneChannel <- true

		}(restaurant)
	}
}
