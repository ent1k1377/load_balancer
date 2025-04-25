package main

import (
	"fmt"
	"github.com/ent1k1377/load_balancer/internal/config"
)

func main() {
	cfg := config.LoadConfig("config.json")

	fmt.Printf("%+v\n", cfg)
}
