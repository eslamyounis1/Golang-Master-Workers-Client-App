package main

import (
	// "encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	// "strconv"
	// "strings"
)

type Config struct {
	ServiceName string
}

func main() {
	args := os.Args
	if len(args) < 3 {
		fmt.Println("Usage: go run main.go [port] [servicename]")
		return
	}
	port := args[1]
	serviceName := args[2]
	addr := fmt.Sprintf(":%s", port)

	handle := Config{ServiceName: serviceName}
	// http.HandleFunc("/", index)
	name := fmt.Sprintf("starting %s", serviceName)
	fmt.Println(name)
	http.Handle("/fasta", &handle)
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		log.Fatal("Server error:", err)
	}

	fmt.Println("Server is running on", addr)
}
func (h *Config) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		h.handleGet(w, req)
	case http.MethodPost:
		h.handlePost(w, req)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}

}
func (h *Config) handleGet(w http.ResponseWriter, req *http.Request) {
	fmt.Println("Handling GET req")
	defer req.Body.Close()

	fileName := fmt.Sprintf("%s.fasta", h.ServiceName)

	// Read the file
	content, err := os.ReadFile(fileName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Println("content ", string(content))

	// Set the response content type
	w.Header().Set("Content-Type", "text/plain")

	// Write the file content as the response body
	fmt.Fprintf(w, "%s", content)
}
func (h *Config) handlePost(w http.ResponseWriter, req *http.Request) {
	fmt.Println("Handling POST req")
	defer req.Body.Close()

	// read req body
	b, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Println("Request: ", string(b))
	fmt.Println("service: ", h.ServiceName)
	fileName := fmt.Sprintf("%s.fasta", h.ServiceName)
	err = os.WriteFile(fileName, b, 0644)
	if err != nil {
		log.Fatal(err)
	}

	w.WriteHeader(http.StatusCreated)
}
