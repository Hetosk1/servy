package reverseproxy

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"servy/config"
	"path"
	"strings"
)

type ProxyRoute struct {
	Path string
	Address string
	Proxy *httputil.ReverseProxy
}

type ReverseProxy struct {
	port string
	routes []ProxyRoute
}

func NewReverseProxy(_port string, _mappings []config.ProxyMapping) (*ReverseProxy, error) {
	var routes []ProxyRoute

	for _, mapping := range _mappings {
		targetURL, err := url.Parse(mapping.Address)
		if err != nil {
			fmt.Printf("Invalid Proxy address: %v\n", err)
			continue
		}

		// Create a new single-host reverse proxy for this route
		proxy := httputil.NewSingleHostReverseProxy(targetURL)

		// Keep the default Director (sets scheme/host/etc.)
		defaultDirector := proxy.Director
		basePath := mapping.Proxy // e.g. "/health", "/admin", "/"

		// Customize the Director to strip the basePath before forwarding
		proxy.Director = func(req *http.Request) {
			defaultDirector(req) // let default behavior set up the backend URL

			// Strip the matching prefix (like "/health") from the incoming path
			trimmed := strings.TrimPrefix(req.URL.Path, basePath)
			if trimmed == "" {
				trimmed = "/"
			}

			req.URL.Path = path.Join(targetURL.Path, trimmed)
			req.URL.RawPath = req.URL.Path
			req.Host = targetURL.Host // make backend see its own host header
		}

		routes = append(routes, ProxyRoute{
			Path:    mapping.Proxy,
			Address: mapping.Address,
			Proxy:   proxy,
		})
	}

	return &ReverseProxy{
		port:   _port,
		routes: routes,
	}, nil
}


func (rp *ReverseProxy) ServeHTTP(w http.ResponseWriter, r *http.Request){
	for _, route := range rp.routes {
		if len(r.URL.Path) >= len(route.Path) && r.URL.Path[:len(route.Path)] == route.Path {
			fmt.Printf("Forwarding [%s] to [%s]\n", route.Path, route.Address)	
			route.Proxy.ServeHTTP(w, r)
			return 
		}
	}

	http.NotFound(w, r)
}

func (rp *ReverseProxy) Run() error{

	listenAddr := ":" + rp.port
	fmt.Println("Starting the Reverse proxy on %v", rp.port)

	return http.ListenAndServe(listenAddr, rp)

}