package main

import (
	"fmt"
	"log"
	"net/http"
)

// helloHandler is a function that will handle all requests to our web server.
// writer is the tool we will use to write our response back to the user.
// request contains all the information about the user's incoming request
func helloHandler(writer http.ResponseWriter, request *http.Request) {
	fmt.Fprintf(writer, "Hello, World!")
}

func main() {
	// http.HandleFunc tells our server that any request it receives for the path "/"
	// should be handled by our "helloHandler" function
	http.HandleFunc("/", helloHandler)

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
