package main 

import (
	"log" // prints timestamped messages
	"net/http" // is the webserver
	"os" // reads environment variables 
	"time" // provides durations 
)

func main(){
	// Read from environment and fall back to port :8080
	addr := os.Getenv("HTTP_ADDR")
	if addr == "" {
		addr = ":8080"
	}
	// :8080 means port 8080 on every network interface
	// localhost:8080 will only reach inside the container 

	// Creates a router: looks at requests & decides which func handles it
	mux := http.NewServeMux()
	// Registers a route. Rejects other methods for me 
	mux.HandleFunc("GET /healthz", handleHealthz)
	
	// http.ListenAndServe(addr, mux) because the shortcut has no timeouts

	// configures the server explicitly instead of using the shortcut 
	srv := &http.Server{
		Addr: addr, 
		Handler: mux, 
		// defend against "Slowloris" attack, where a client sends headers 
		// one byte at a time to hold connections open until you run out 
		ReadHeaderTimeout: 5*time.Second, 
		// Cap how long reading a request response can take Read & Write
		ReadTimeout: 10*time.Second, 
		WriteTimeout: 10*time.Second, 
		// Closes kept-alive connections that sit unused
		IdleTimeout: 60*time.Second, 
	}
	log.Printf("api listening on  %s", addr)
	// blocks forever while serving and only returns if something goes wrong.
	// The most common error is bind: address already in use
	// meaning another copy of the server is still running on 8080
	if err := srv.ListenAndServe(); err != nil {
		// above is Go's stanard error-handling pattern: functions return errors as values and check them explicitly 
		log.Fatal(err) // print and exits with status 1, tells Docker process failed 
	}
}
// The handler
// Every Go handler uses w (write the response) and r (incoming request)
// r is unused and go allows unused function parameters 
func handleHealthz(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plan; charset=utf-8")
	// optional below, 200 is by default
	w.WriteHeader(http.StatusOK)
	// returns a byte count error 
	// if the client alrdy disconnected _, _= deliberately discards both
	_, _ = w.Write([]byte("ok\n"))
}

