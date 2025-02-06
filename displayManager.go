package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
)

// DisplayManager manages completed orders and serves them via HTTP
type DisplayManager struct {
	mu      sync.Mutex
	orders  []Order
	httpSrv *http.Server
}

// AddOrder adds an order to the DisplayManager
func (dm *DisplayManager) AddOrder(order Order) {
	dm.mu.Lock()
	defer dm.mu.Unlock()
	dm.orders = append(dm.orders, order)
	fmt.Printf("Order added to DisplayManager: %+v\n", order)
}

// GetOrders handles HTTP requests to retrieve the list of completed orders
func (dm *DisplayManager) GetOrders(w http.ResponseWriter, r *http.Request) {
	dm.mu.Lock()
	defer dm.mu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(dm.orders); err != nil {
		http.Error(w, "Failed to encode orders", http.StatusInternalServerError)
		fmt.Printf("Error encoding orders: %v\n", err)
	}
}

// StartServer starts the HTTP server for the DisplayManager
func (dm *DisplayManager) StartServer(port int) {
	mux := http.NewServeMux()
	mux.HandleFunc("/orders", dm.GetOrders)

	dm.httpSrv = &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: mux,
	}

	go func() {
		fmt.Printf("Starting HTTP server on port %d\n", port)
		if err := dm.httpSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("HTTP server error: %v\n", err)
		}
	}()
}

// StopServer stops the HTTP server gracefully
func (dm *DisplayManager) StopServer() {
	if dm.httpSrv != nil {
		fmt.Println("Stopping HTTP server...")
		if err := dm.httpSrv.Close(); err != nil {
			fmt.Printf("Error stopping server: %v\n", err)
		}
	}
}
