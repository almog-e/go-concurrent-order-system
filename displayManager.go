package main

import (
	"fmt"
)

func StartDisplayManager(DisplayChannel chan string) {
	fmt.Println("Display Manager started")
	for display := range DisplayChannel {
		fmt.Println("\033[34m" + "Manager Display: " + display + "\033[0m")
	}
	fmt.Println("Display Manager stopped")
}
