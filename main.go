package main

import (
	"fmt"
	"log"
	"os"

	"load-balancer/config" // YAML Configuration package
	"load-balancer/loadbalancer" // Load Balancer core package
)

func main() {
	// Loaded YAML Configuration 
	cfg, err := config.LoadConfig("/etc/servy/servy.yaml")
	if err != nil {
		log.Fatalf("Error loading configuration: %v", err)
	}

	config.CheckAndLogConfig(cfg)
	fmt.Println("Configuration loaded successfully.")

	// Check if loadbalancer switch is on
	if cfg.LoadBalancer.Service != "on" {
		fmt.Println("Load balancer service is not set to 'on'. Exiting.")
		os.Exit(0)
	}

	// Load Servers into the load balancer
	servers := []loadbalancer.Server{}
	for _, address := range cfg.LoadBalancer.Servers {
		server, err := loadbalancer.NewSimpleServer(address)
		if err != nil {
			log.Fatalf("Error creating simple server for address %s: %v", address, err)
		}
		servers = append(servers, server)
	}

	// Spinning up the Load Balancer
	lb := loadbalancer.NewLoadBalancer(cfg.LoadBalancer.Port, servers)
	if err := lb.Run(); err != nil {
		log.Fatalf("Load Balancer failed to run: %v", err)
	}
}
