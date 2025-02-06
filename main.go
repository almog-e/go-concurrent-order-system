package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

type Restaurant struct {
	ID        int `json:"ID"`
	Orders    int `json:"Orders"`
	QueueSize int `json:"Queue_Size"`
}

type Config struct {
	Restaurants       []Restaurant `json:"RESTAURANTS"`
	PizzaZoneWorkers  int          `json:"PizzaZone_Workers"`
	BurgerZoneWorkers int          `json:"BurgerZone_Workers"`
	SushiZoneWorkers  int          `json:"SushiZone_Workers"`
	ZoneQueueSize     int          `json:"Zone_Queue_Size"`
	HTTPServerPort    int          `json:"HTTP_Server_Port"`
}

func LoadConfig(filePath string) (*Config, error) {
	// Open the JSON file
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open config file: %v", err)
	}
	defer file.Close()

	// Decode the JSON into the Config struct
	var config Config
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		return nil, fmt.Errorf("failed to parse config file: %v", err)
	}

	return &config, nil
}

func main() {
	// Create a channel to communicate between goroutines
	configFilePath := "config.json"
	config, err := LoadConfig(configFilePath)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Print the configuration
	fmt.Printf("HTTP Server Port: %d\n", config.HTTPServerPort)
	fmt.Println("Restaurants:")
	for _, restaurant := range config.Restaurants {
		fmt.Printf("Restaurant ID: %d, Orders: %d, Queue Size: %d\n",
			restaurant.ID, restaurant.Orders, restaurant.QueueSize)
	}
	doneChannel := make(chan bool, len(config.Restaurants))
	ordersChannel := make(chan Order, 10) // Buffered channel for orders
	pizzaZone := make(chan Order, config.ZoneQueueSize)
	burgerZone := make(chan Order, config.ZoneQueueSize)
	sushiZone := make(chan Order, config.ZoneQueueSize)

	StartProducers(config.Restaurants, ordersChannel, doneChannel)
	StartDispatcher(ordersChannel, pizzaZone, burgerZone, sushiZone)
	doneZonesChannel := make(chan bool, 3)

	displayManager := &DisplayManager{}

	go displayManager.StartServer(config.HTTPServerPort)

	StartZoneManager("PizzaZone", pizzaZone, config.PizzaZoneWorkers, doneZonesChannel, displayManager)
	StartZoneManager("BurgerZone", burgerZone, config.BurgerZoneWorkers, doneZonesChannel, displayManager)
	StartZoneManager("SushiZone", sushiZone, config.SushiZoneWorkers, doneZonesChannel, displayManager)

	// Wait for all zones to complete
	go func() {
		for i := 0; i < 3; i++ { // Wait for 3 true values
			<-doneZonesChannel
		}
		fmt.Println("All zones completed!")
		close(doneZonesChannel)
	}()

	// Wait for all producers to complete
	for i := 0; i < len(config.Restaurants); i++ {
		<-doneChannel
	}
	close(ordersChannel)
	sigs := make(chan os.Signal, 1)
	fmt.Println("Server is running. Press Ctrl+C to stop.")

	signal.Notify(sigs, os.Interrupt, syscall.SIGTERM)

	<-sigs // מחכים לסיגנל
	fmt.Println("Shutting down the server...")
	fmt.Println("All goroutines completed!")

}
