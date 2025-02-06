# go concurrent Restaurant Zone Management System

## Overview
This program simulates a **restaurant order management system** using **Go**, where orders from multiple restaurants are dynamically processed and dispatched to specific zones (Pizza, Burger, Sushi).

---

## **Features**
- Concurrent order processing using Go’s **goroutines** and **channels**.
- Dynamic order dispatching to zones based on order type.
- Configurable system via a simple `config.txt` file.
- Graceful shutdown and cleanup of all tasks.

---

## **Project Structure**
```
restaurant-zone-system/
├── dispatcher.go        # Dispatcher that routes orders to appropriate zones
├── go.mod               # Module file for managing dependencies
├── main.go              # Main entry point of the program
├── producers.go         # Generates and sends restaurant orders
├── server.go            # HTTP server for monitoring orders
├── zoneManagers.go      # Manages processing of orders in each zone
└── config.txt           # Configuration file for restaurants and zones
```

---

## **How to Run**
1. **Clone the project:**
   ```bash
   git clone https://github.com/your-username/restaurant-zone-system.git
   cd restaurant-zone-system
   ```

2. **Ensure Go is installed:**
   ```bash
   go version
   ```

3. **Run the program:**
   ```bash
   go run .
   ```

4. **View the orders:**
   Visit the following URL to see the list of completed orders:
   ```
   http://localhost:8080/orders
   ```

---

## **Configuration (`config.txt`)**
Define the configuration of the system in `config.txt`. Example:
```plaintext
RESTAURANT 1
Orders: 10
Queue Size: 5

RESTAURANT 2
Orders: 8
Queue Size: 4

PizzaZone Workers: 3
BurgerZone Workers: 2
SushiZone Workers: 1
Zone Queue Size: 10
HTTP Server Port: 8080
```
- **RESTAURANT:** Defines the ID and settings for each restaurant.
- **Orders:** Number of orders the restaurant will generate.
- **Queue Size:** Buffer size for each restaurant’s order queue.
- **Zone Workers:** Number of workers handling orders in each zone.
- **Zone Queue Size:** Buffer size for the zone queues.
- **HTTP Server Port:** Port number for the HTTP server.

---

## **How It Works**
1. **Producers:** Each restaurant generates orders for food items (Pizza, Burger, Sushi) and sends them through a channel.
2. **Dispatcher:** The dispatcher routes orders to the correct zone based on the order type.
3. **Zone Managers:** Each zone manager processes orders concurrently using worker tokens to control concurrency.
4. **HTTP Server:** The server monitors completed orders and provides real-time feedback.

---

## **System Flow**
1. Producers generate orders and send them to the dispatcher.
2. The dispatcher routes orders to the correct zone channel (Pizza, Burger, Sushi).
3. Zone managers process orders concurrently and update the list of completed orders.
4. The HTTP server allows viewing completed orders through the browser.

---

## **Stopping the Program**
To stop the program gracefully:
1. Press `Enter` in the terminal.
2. This triggers a cancellation signal, and all running goroutines will clean up and exit.

---


