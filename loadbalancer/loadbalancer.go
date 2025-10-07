package loadbalancer

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
)

// Server defines the interface for a backend server.
type Server interface {
	GetAddress() string
	IsAlive() bool
	Serve(rw http.ResponseWriter, r *http.Request)
}

// simpleServer is a concrete implementation of the Server interface.
type simpleServer struct {
	address string
	proxy   *httputil.ReverseProxy
}

// NewSimpleServer creates a new simpleServer instance.
func NewSimpleServer(address string) (Server, error) {
	serverUrl, err := url.Parse(address)
	if err != nil {
		return nil, fmt.Errorf("failed to parse server address %s: %w", address, err)
	}

	return &simpleServer{
		address: address,
		proxy:   httputil.NewSingleHostReverseProxy(serverUrl),
	}, nil
}

// GetAddress returns the server's URL address.
func (s *simpleServer) GetAddress() string {
	return s.address
}

// IsAlive always returns true for a simple server.
// In a real application, this would involve a health check.
func (s *simpleServer) IsAlive() bool {
	return true
}

// Serve handles the request by proxying it to the backend server.
func (s *simpleServer) Serve(rw http.ResponseWriter, r *http.Request) {
	s.proxy.ServeHTTP(rw, r)
}

// LoadBalancer manages a pool of servers and distributes requests.
type LoadBalancer struct {
	port            string
	roundRobinCount int
	servers         []Server
}

// NewLoadBalancer creates a new LoadBalancer instance.
func NewLoadBalancer(port string, servers []Server) *LoadBalancer {
	return &LoadBalancer{
		port:            port,
		roundRobinCount: 0,
		servers:         servers,
	}
}

// GetNextAvailableServer uses a round-robin strategy to select the next available server.
func (lb *LoadBalancer) GetNextAvailableServer() Server {
	server := lb.servers[lb.roundRobinCount%len(lb.servers)]

	// Simple check, but in a robust system, this loop needs care to avoid infinite looping
	// if all servers are down.
	for !server.IsAlive() {
		lb.roundRobinCount++
		server = lb.servers[lb.roundRobinCount%len(lb.servers)]
	}

	// Increment for the next request.
	lb.roundRobinCount++

	return server
}

// ServeProxy is the HTTP handler that forwards requests to the chosen backend server.
func (lb *LoadBalancer) ServeProxy(rw http.ResponseWriter, r *http.Request) {
	targetServer := lb.GetNextAvailableServer()
	fmt.Printf("Forwarding request to: %s from %s\n", targetServer.GetAddress(), r.RemoteAddr)
	targetServer.Serve(rw, r)
}

// Run starts the HTTP server on the configured port.
func (lb *LoadBalancer) Run() error {
	listenAddr := ":" + lb.port
	fmt.Printf("Starting load balancer on %s\n", listenAddr)

	http.HandleFunc("/", lb.ServeProxy)
	return http.ListenAndServe(listenAddr, nil)
}