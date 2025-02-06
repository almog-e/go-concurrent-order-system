package main

import "fmt"

func StartDispatcher(ordersChannel chan Order, pizzaZone chan Order, burgerZone chan Order, sushiZone chan Order) {
	go func() {
		for order := range ordersChannel {
			// מפיץ לפי סוג המזון
			switch order.FoodType {
			case "Pizza":
				fmt.Println("Pizza")
				pizzaZone <- order
			case "Burger":
				fmt.Println("Burger")
				burgerZone <- order
			case "Sushi":
				fmt.Println("Sushi")
				sushiZone <- order
			default:
				fmt.Printf("Unknown order type: %s\n", order.FoodType)
			}
		}
		close(pizzaZone)
		close(burgerZone)
		close(sushiZone)
	}()
}
