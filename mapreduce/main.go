package main

import (
	"flag"
	"html/template"
	"log"
	"os/exec"
	"runtime"
	"strings"

	// "encoding/json"
	"fmt"
	"io"
	"sync"

	// "log"
	"net/http"
)

type Config struct {
	IP1 string
	IP2 string
	IP3 string
}

func (c *Config) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	m := make(map[string]string)
	m["slave1"] = c.IP1
	m["slave2"] = c.IP2
	m["slave3"] = c.IP3

	var str string

	response, err := fetchResponses(m)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	totalWordCount := 0
	totalCharacterCount := 0
	for serviceName, text := range response {
		// Split the string into words
		words := strings.Fields(text)

		// Count the number of words in the current string
		wordCount := len(words)
		characters := len(text)
		totalCharacterCount += characters
		totalWordCount += wordCount

		data := fmt.Sprintf("\n%s Characters: %v Words: %v\n", serviceName, characters, wordCount)

		str += data
	}
	data := fmt.Sprintf("\nTotal Characters: %v \nTotal Words: %v\n", totalCharacterCount, totalWordCount)
	str += data

	display(str)

	//fmt.Fprintf(w, html)
}

func main() {
	var port string
	var IP1 string
	var IP2 string
	var IP3 string
	var mapreduceIP string

	flag.StringVar(&IP1, "ip1", "https://8d93-105-35-226-5.eu.ngrok.io/fasta", "slave 1 IP address")
	flag.StringVar(&IP2, "ip2", "https://f7c1-105-40-174-223.eu.ngrok.io/fasta", "slave 2 IP address")
	flag.StringVar(&IP3, "ip3", "http://localhost:8099/fasta", "slave 3 IP address")
	flag.StringVar(&mapreduceIP, "mapr", "http://localhost:8010/fasta", "map reduce IP address")
	flag.StringVar(&port, "port", "8010", "port")
	flag.Parse()

	config := Config{
		IP1: IP1,
		IP2: IP2,
		IP3: IP3,
	}

	addr := fmt.Sprintf(":%s", port)

	http.Handle("/fasta", &config)
	// http.HandleFunc("/", index)
	fmt.Println("starting map reduce")
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		log.Fatal("Server error:", err)
	}

	fmt.Println("Server is running on", addr)
}
func fetchResponses(m map[string]string) (map[string]string, error) {
	responses := make(map[string]string)
	var wg sync.WaitGroup
	var mutex sync.Mutex

	for key, value := range m {
		wg.Add(1)

		go func(serviceName string, url string) {
			defer wg.Done()

			resp, err := http.Get(url)
			if err != nil {
				fmt.Printf("Error fetching %s: %s\n", url, err.Error())
				return
			}
			defer resp.Body.Close()

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				fmt.Printf("Error reading response body from %s: %s\n", url, err.Error())
				return
			}

			// Store the response in the array
			mutex.Lock()
			responses[serviceName] = string(body)
			mutex.Unlock()
		}(key, value)
	}

	wg.Wait()

	return responses, nil
}

// openBrowser opens the specified URL in the default browser
func openBrowser(url string) error {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "darwin":
		cmd = "open"
		args = []string{url}
	case "linux":
		cmd = "xdg-open"
		args = []string{url}
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start", url}
	default:
		return fmt.Errorf("unsupported platform")
	}

	return exec.Command(cmd, args...).Start()
}

func display(text string) {
	// Define the handler function
	handler := func(w http.ResponseWriter, r *http.Request) {

		// Define the HTML template with the response text
		html := `
			<!DOCTYPE html>
			<html>
			<head>
				<title>Response</title>
			</head>
			<body>
				<h1>Response:</h1>
				<pre>
					{{.}}
				</pre>
			</body>
			</html>
		`

		// Create a template from the HTML string
		tmpl := template.Must(template.New("response").Parse(html))

		// Execute the template with the response text
		err := tmpl.Execute(w, text)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			log.Println("Error executing template:", err)
			return
		}
	}

	// Register the handler function to handle requests
	http.HandleFunc("/", handler)

	// Start the server in a separate goroutine
	go func() {
		log.Println("Server listening on http://localhost:8081")
		err := http.ListenAndServe(":8081", nil)
		if err != nil {
			log.Fatal("Server error:", err)
		}
	}()

	// Launch the browser with the server's URL
	openBrowser("http://localhost:8081")

	// Wait for a key press to exit
	fmt.Println("Press any key to exit...")
	var input string
	fmt.Scanln(&input)
}
