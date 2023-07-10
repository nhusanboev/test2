package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"sync"
	"time"
)

type Config struct {
	WebServices []string `json:"web_services"`
}

type ServiceStats struct {
	ServiceURL string `json:"service_url"`
	RequestCount uint64 `json:"request_count"`
}

type LoadBalancer struct {
	services       []*url.URL
	currentIndexes []int
	stats          []*ServiceStats
	mutex          sync.Mutex
}

func NewLoadBalancer(configFile string) (*LoadBalancer, error) {
	file, err := ioutil.ReadFile(configFile)
	if err != nil {
		return nil, err
	}

	var config Config
	err = json.Unmarshal(file, &config)
	if err != nil {
		return nil, err
	}

	services := make([]*url.URL, len(config.WebServices))
	currentIndexes := make([]int, len(config.WebServices))
	stats := make([]*ServiceStats, len(config.WebServices))
	for i, addr := range config.WebServices {
		u, err := url.Parse(addr)
		if err != nil {
			return nil, err
		}
		services[i] = u
		currentIndexes[i] = 0
		stats[i] = &ServiceStats{
			ServiceURL:   addr,
			RequestCount: 0,
		}
	}

	return &LoadBalancer{
		services:       services,
		currentIndexes: currentIndexes,
		stats:          stats,
	}, nil
}

func (lb *LoadBalancer) getNextService() *url.URL {
	lb.mutex.Lock()
	defer lb.mutex.Unlock()

	// Get the current index for selecting the next service
	currentIndex := lb.currentIndexes[0]

	// Get the next service based on the current index
	service := lb.services[currentIndex]

	// Update the current index for the next request
	lb.currentIndexes[0] = (currentIndex + 1) % len(lb.services)

	return service
}


func (lb *LoadBalancer) handleRequest(w http.ResponseWriter, r *http.Request) {
	service := lb.getNextService()

	lb.mutex.Lock()
	for _, stat := range lb.stats {
		if stat.ServiceURL == service.String() {
			stat.RequestCount++
			break
		}
	}
	lb.mutex.Unlock()

	proxy := httputil.NewSingleHostReverseProxy(service)
	proxy.ServeHTTP(w, r)
}

func main() {
	configFile := "config.json"
	loadBalancer, err := NewLoadBalancer(configFile)
	if err != nil {
		log.Fatal(err)
	}

	// Запускаем горутину, которая будет периодически проверять доступность сервисов
	go loadBalancer.checkServiceAvailability()

	// Запускаем горутину, которая будет периодически показывать отчет о количестве запросов
	go loadBalancer.showRequestStats()

	http.HandleFunc("/", loadBalancer.handleRequest)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Load balancer is running on port %s\n", port)
	err = http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
	if err != nil {
		log.Fatal(err)
	}
}

func (lb *LoadBalancer) showRequestStats() {
	for {
		time.Sleep(10 * time.Second)

		lb.mutex.Lock()
		for _, stat := range lb.stats {
			if contains(lb.services, stat.ServiceURL) {
				log.Printf("Service: %s, Request count: %d\n", stat.ServiceURL, stat.RequestCount)
			}
		}
		lb.mutex.Unlock()

		fmt.Println("----------------------------")
	}
}


func (lb *LoadBalancer) checkServiceAvailability() {
	for {
		time.Sleep(5 * time.Second)

		// Create a new slice to store available services
		availableServices := make([]*url.URL, 0)

		// Check availability of each service
		for _, service := range lb.services {
			// Send a GET request to the health check endpoint of the service
			resp, err := http.Get(service.String() + "/health-check")
			if err != nil {
				log.Printf("Service %s is unavailable\n", service.String())
				continue
			}

			// If the response status code is OK, add the service to the available services
			if resp.StatusCode == http.StatusOK {
				availableServices = append(availableServices, service)
			} else {
				log.Printf("Service %s is unavailable\n", service.String())
			}

			resp.Body.Close()
		}

		// Check if any previously disabled services have been re-enabled
		for _, disabledService := range lb.stats {
			if !contains(availableServices, disabledService.ServiceURL) {
				// Parse the disabled service URL to convert it to *url.URL
				serviceURL, err := url.Parse(disabledService.ServiceURL)
				if err != nil {
					log.Printf("Failed to parse service URL: %s\n", disabledService.ServiceURL)
					continue
				}

				// Send a GET request to the health check endpoint of the disabled service
				resp, err := http.Get(serviceURL.String() + "/health-check")
				if err == nil && resp.StatusCode == http.StatusOK {
					log.Printf("Service %s is now available\n", disabledService.ServiceURL)
					availableServices = append(availableServices, serviceURL)
				}
			}
		}


		// If the available services haven't changed, continue
		if len(availableServices) == len(lb.services) {
			continue
		}

		// Update the list of available services
		lb.mutex.Lock()
		lb.services = availableServices
		lb.mutex.Unlock()
	}
}

func contains(services []*url.URL, serviceURL string) bool {
	for _, service := range services {
		if service.String() == serviceURL {
			return true
		}
	}
	return false
}
