package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

// Configuration struct
type Config struct {
	Restaurants   []Restaurant
	ZoneWorkers   map[string]int
	ZoneQueueSize int
}

// Restaurant struct
type Restaurant struct {
	ID        int
	Orders    int
	QueueSize int
}

func main() {
	time.Sleep(6 * time.Second)
	configFile := "config.txt"
	config, err := parseConfig(configFile)
	if err != nil {
		fmt.Println("Error reading config file:", err)
		return
	}

	// Create context with cancellation
	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup

	// Create channels for each restaurant
	channels := make(map[int]chan string) // Use chan string
	for _, restaurant := range config.Restaurants {
		channels[restaurant.ID] = make(chan string, restaurant.QueueSize)
	}

	// Create zone-specific channels and token pools
	pizzaChannel := make(chan string, config.ZoneQueueSize)
	burgerChannel := make(chan string, config.ZoneQueueSize)
	sushiChannel := make(chan string, config.ZoneQueueSize)

	pizzaTokens := make(chan struct{}, config.ZoneWorkers["PizzaZone"])
	burgerTokens := make(chan struct{}, config.ZoneWorkers["BurgerZone"])
	sushiTokens := make(chan struct{}, config.ZoneWorkers["SushiZone"])
	DisplayChannel := make(chan string)
	// Start the HTTP server for Display Manager
	go StartDisplayManager(DisplayChannel)

	// Start Zone Managers
	wg.Add(3)
	go func() {
		defer wg.Done()
		ZoneManager("PizzaZone", pizzaChannel, DisplayChannel, pizzaTokens, ctx)
	}()
	go func() {
		defer wg.Done()
		ZoneManager("BurgerZone", burgerChannel, DisplayChannel, burgerTokens, ctx)
	}()
	go func() {
		defer wg.Done()
		ZoneManager("SushiZone", sushiChannel, DisplayChannel, sushiTokens, ctx)
	}()

	// Start Dispatcher
	wg.Add(1)
	go func() {
		defer wg.Done()
		Dispatcher(channels, pizzaChannel, burgerChannel, sushiChannel, ctx)
	}()

	// Start producers for each restaurant
	for _, restaurant := range config.Restaurants {
		wg.Add(1)
		go func(restaurant Restaurant) {
			defer wg.Done()
			StartProducer(restaurant, channels[restaurant.ID], ctx)
		}(restaurant)
	}

	// Wait for all goroutines to finish
	go func() {
		wg.Wait()
		close(DisplayChannel)
		fmt.Println("Main: All work completed.")
	}()

	time.Sleep(8 * time.Second)

	// Perform graceful shutdown
	fmt.Println("Shutting down...")
	fmt.Println("3")
	time.Sleep(1 * time.Second)
	fmt.Println("2")
	time.Sleep(1 * time.Second)
	fmt.Println("1")
	time.Sleep(1 * time.Second)
	fmt.Println("Bye!")

	cancel() // Cancel the context for all running goroutines
}

func parseConfig(filename string) (Config, error) {
	file, err := os.Open(filename)
	if err != nil {
		return Config{}, err
	}
	defer file.Close()

	var config Config
	config.ZoneWorkers = make(map[string]int)

	scanner := bufio.NewScanner(file)
	var currentRestaurant *Restaurant

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue // Skip empty lines
		}

		if strings.HasPrefix(line, "RESTAURANT") {
			if currentRestaurant != nil {
				config.Restaurants = append(config.Restaurants, *currentRestaurant)
			}
			idStr := strings.TrimPrefix(line, "RESTAURANT ")
			id, err := strconv.Atoi(strings.TrimSpace(idStr))
			if err != nil {
				return Config{}, fmt.Errorf("invalid restaurant ID format: %s", line)
			}

			currentRestaurant = &Restaurant{ID: id}
		} else if strings.HasPrefix(line, "Orders:") {
			if currentRestaurant != nil {
				currentRestaurant.Orders = parseValue(line)
			}
		} else if strings.HasPrefix(line, "Queue Size:") {
			if currentRestaurant != nil {
				currentRestaurant.QueueSize = parseValue(line)
			}
		} else if strings.Contains(line, "Zone Workers") {
			parts := strings.Split(line, " ")
			if len(parts) >= 3 {
				zone := parts[0]
				workers := parseValue(line)
				config.ZoneWorkers[zone] = workers
			}
		} else if strings.HasPrefix(line, "Zone Queue Size:") {
			config.ZoneQueueSize = parseValue(line)
		}
	}

	// Add the last restaurant
	if currentRestaurant != nil {
		config.Restaurants = append(config.Restaurants, *currentRestaurant)
	}

	if err := scanner.Err(); err != nil {
		return Config{}, err
	}

	return config, nil
}

func parseValue(line string) int {
	parts := strings.Split(line, ":")
	if len(parts) < 2 {
		return 0
	}
	value, _ := strconv.Atoi(strings.TrimSpace(parts[1]))
	return value
}
