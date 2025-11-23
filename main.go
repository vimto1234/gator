package main

import (
	"fmt"

	"github.com/vimto1234/gator/internal/config"
)

func main() {
	gatorConfig, err := config.Read()
	if err != nil {
		fmt.Printf("Error reading config: %v", err)
	}

	fmt.Printf("Init config: %v\n", gatorConfig)

	gatorConfig.SetUser("david")

	newConfig, err := config.Read()
	if err != nil {
		fmt.Printf("Error reading config: %v", err)
	}
	fmt.Printf("new config: %v\n", newConfig)

}
