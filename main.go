package main

import (
	"fmt"
	"log"
	"sync"
	"servy/config"
	"servy/loadbalancer"
	"servy/reverseproxy"
)

func main() {
	// Load YAML configuration
	cfg, err := config.LoadConfig("/etc/servy/servy.yaml")
	if err != nil {
		log.Fatalf("Error loading configuration: %v", err)
	}

	config.CheckAndLogConfig(cfg)
	fmt.Println("✅ Configuration loaded successfully.")

	// WaitGroup to keep main alive while goroutines run
	var wg sync.WaitGroup

	// ========== Load Balancer ==========
	if cfg.LoadBalancer.Service == "on" {
		wg.Add(1)
		go func() {
			defer wg.Done()
			fmt.Println("🚀 Starting Load Balancer service...")

			var servers []loadbalancer.Server
			for _, address := range cfg.LoadBalancer.Servers {
				server, err := loadbalancer.NewSimpleServer(address)
				if err != nil {
					log.Fatalf("❌ Error creating Load Balancer server (%s): %v", address, err)
				}
				servers = append(servers, server)
			}

			lb := loadbalancer.NewLoadBalancer(cfg.LoadBalancer.Port, servers)
			if err := lb.Run(); err != nil {
				log.Fatalf("💥 Load Balancer failed: %v", err)
			}
		}()
	} else {
		fmt.Println("⛔ Load Balancer service: OFF")
	}

	// ========== Reverse Proxy (NOT IMPLEMENTED) ==========
	

	if cfg.ReverseProxy.Service == "on" {
		wg.Add(1)
		go func() {
			defer wg.Done()
			fmt.Println("Starting Reverse Proxy services...")

			rp, err := reverseproxy.NewReverseProxy(cfg.ReverseProxy.Port, cfg.ReverseProxy.Proxies)
			if err != nil {
				log.Fatalf("Error starting reverse proxy services: %v", err)
			}

			err = rp.Run() // no colon here
			if err != nil {
				log.Fatalf("Error running reverse proxy services: %v", err)
				return
			}

		}()
	} else {
		fmt.Println("⛔ Reverse Proxy service: OFF")
	}

	// Wait for all running services (currently only the load balancer) to finish.
	wg.Wait()
}
