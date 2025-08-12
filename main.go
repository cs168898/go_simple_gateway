package main

import (
	"io"
	"log"
	"net/http"
	"net/url"
	"sync"
)

// create a map with string keys and slice of bytes.
var cache = make(map[string][]byte)

// a mutex to protect cache from concurrent access.
var mutex = &sync.RWMutex{}

// proxyHandler is a function that will receive the original request and forward it to our target service.
// writer is the tool we will use to write our response back to the user.
// request contains all the information about the user's incoming request
func proxyHandler(writer http.ResponseWriter, request *http.Request) {
	// define our target url to forward our requests to.

	// create the cache key using the url string
	cacheKey := request.URL.String()

	// lock the reader to safely read the cache
	mutex.RLock()
	// now we check the cache to see if the key already exists
	cachedResponse, found := cache[cacheKey]
	mutex.RUnlock()

	if found {
		log.Println("Cache HIT for key:", cacheKey)

		// if the cached response is already in the map , return it, no need proxy.
		writer.Write(cachedResponse)
		return
	}

	// if not in cache, proxy the request

	log.Println("Cache MISS for key:", cacheKey, ". Forwarding request...")
	// the underscore means that there is a possible error variable value there , but we dont need it, so throw it away.
	targetUrl := "https://httpbin.org"
	target, _ := url.Parse(targetUrl)
	target.Path = request.URL.Path
	target.RawQuery = request.URL.RawQuery

	log.Println("forwarding request to:", target.String())

	// create a new request to the target Url.
	proxyReq, err := http.NewRequest(request.Method, target.String(), request.Body)
	if err != nil {
		// error variable contains something , means that an error exists...
		// if we failed to create the request, send an error to the user.
		http.Error(writer, "Error creating proxy request", http.StatusInternalServerError)
		log.Println("Error creating proxy request:", err)
	}

	// we copy the headers from the original request to our new proxy request
	proxyReq.Header = request.Header

	//create a client to send the request
	client := &http.Client{}

	resp, err := client.Do(proxyReq)
	if err != nil {
		// error variable contains something , means that an error exists...
		// if we failed to create the request, send an error to the user.
		http.Error(writer, "Error forwarding request", http.StatusInternalServerError)
		log.Println("Error forwarding request:", err)
		return
	}

	// close the response body when done with it
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(writer, "Error reading response body", http.StatusInternalServerError)
		return
	}

	mutex.Lock()

	// store into cache map
	cache[cacheKey] = body

	mutex.Unlock()

	// copy the status code from the target's response to our own response
	writer.WriteHeader(resp.StatusCode)

	// copy the headers from the target's response to our response
	for key, values := range resp.Header {
		for _, value := range values {
			writer.Header().Add(key, value)
		}
	}

	// write the body
	writer.Write(body)

}

func main() {
	// http.HandleFunc tells our server that any request it receives for the path "/"
	// should be handled by our "helloHandler" function
	http.HandleFunc("/", proxyHandler)

	// http.ListenAndServe starts the web server.
	// it listens on port 8080 for any incoming network connections
	// and if it fails to sdtart, it will cause a fatal error.
	log.Println("Server starting on port 8080...")

	// the semicolon separates the initializer from the condition
	// the listen and serve function will return a non nil error if it fails to start a server
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Could not start server: %s\n", err)
	}

}
