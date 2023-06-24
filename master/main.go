package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"math"
	"net"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

type SlaveDevice struct {
	id     string
	ipAddr string
	data   string
}

type Response struct {
	Id        string
	IpAddress []string
}
type Config struct {
	Slave1    SlaveDevice
	Slave2    SlaveDevice
	Slave3    SlaveDevice
	MapReduce SlaveDevice
}

func (h *Config) makeChunks(fileName string) {
	// read file in chunks
	f, err := os.Open(fileName)
	panicOnErrorM(err)
	defer f.Close()
	sizeOfChuck, err := f.Seek(0, 2)
	sizeOfChuck = int64(math.Ceil(float64(sizeOfChuck) / 3.0))

	panicOnErrorM(err)

	f.Seek(0, 0)
	reader := bufio.NewReader(f)
	for i := 1; i < 4; i++ {

		b := make([]byte, sizeOfChuck)
		numberOfBytesRead, err := reader.Read(b)
		if err != nil {
			if !errors.Is(err, io.EOF) {
				fmt.Println("Error reading file:", err)
			}
			break
		}
		// fmt.Println(string(b[0:numberOfBytesRead]))
		if i == 1 {
			postChunk(h.Slave1.ipAddr, string(b[0:numberOfBytesRead]))
		} else if i == 2 {
			postChunk(h.Slave2.ipAddr, string(b[0:numberOfBytesRead]))
		} else {
			postChunk(h.Slave3.ipAddr, string(b[0:numberOfBytesRead]))
		}

	}
}
func (h *Config) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	fmt.Println("Handling GET req")
	// http://localhost:8090/fasta?id=0
	query := req.URL.Query()
	id := query.Get("id")
	fmt.Println("ID:    ", id)
	fmt.Println(query)

	var result string
	var err error
	var idInt int
	var resp Response
	if id != "" {
		idInt, err = strconv.Atoi(id)
		switch idInt {
		case 0:
			resp.Id = "0"
			resp.IpAddress = append(resp.IpAddress, []string{h.Slave1.ipAddr, h.Slave2.ipAddr, h.Slave3.ipAddr, h.MapReduce.ipAddr}...)
		case 1:
			resp.Id = "1"
			resp.IpAddress = append(resp.IpAddress, []string{h.Slave1.ipAddr, h.MapReduce.ipAddr}...)
		case 2:
			resp.Id = "2"
			resp.IpAddress = append(resp.IpAddress, []string{h.Slave2.ipAddr, h.MapReduce.ipAddr}...)
		case 3:
			resp.Id = "3"
			resp.IpAddress = append(resp.IpAddress, []string{h.Slave3.ipAddr, h.MapReduce.ipAddr}...)
		default:
			result = "Write id value from 0 to 3"
		}
		fmt.Println("result: ", result)
	} else {
		result = "Write id value from 0 to 3 ,example:(http://localhost:8090/fasta?id=1)"
		fmt.Println("ID:    ", id)
	}
	// if we had any error return status 500 and error
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	b, er := json.Marshal(resp)
	if er != nil {
		fmt.Println(er)
	}
	// set header return data
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, string(b))
}

func main() {
	var port string
	var IP1 string
	var IP2 string
	var IP3 string
	var mapreduceIP string

	flag.StringVar(&IP1, "ip1", "http://localhost:8097/fasta", "slave 1 IP address")
	flag.StringVar(&IP2, "ip2", "http://localhost:8098/fasta", "slave 2 IP address")
	flag.StringVar(&IP3, "ip3", "http://localhost:8099/fasta", "slave 3 IP address")
	flag.StringVar(&mapreduceIP, "mapr", "http://localhost:8010/fasta", "map reduce IP address")
	flag.StringVar(&port, "port", "8011", "port")
	flag.Parse()

	fmt.Println(IP1, IP2, IP3)

	var slave1 = SlaveDevice{id: "1", ipAddr: IP1}
	var slave2 = SlaveDevice{id: "2", ipAddr: IP2}
	var slave3 = SlaveDevice{id: "3", ipAddr: IP3}
	var mapreduce = SlaveDevice{id: "4", ipAddr: mapreduceIP}

	conf := Config{
		Slave1:    slave1,
		Slave2:    slave2,
		Slave3:    slave3,
		MapReduce: mapreduce,
	}

	conf.makeChunks("master/Big_Data.txt")
	conf.masterAsServer(port)
}

func (h *Config) masterAsServer(port string) {
	addr := fmt.Sprintf(":%s", port)
	http.Handle("/fasta", h)
	fmt.Println("starting master server")
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		log.Fatal("Server error:", err)
	}

	fmt.Println("Server is running on", addr)

}
func panicOnErrorM(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// checkClientAvailability validates if a client is available.
func checkClientAvailability(ipAddress string, timeout time.Duration) error {
	parsedURL, err := url.Parse(ipAddress)
	if err != nil {
		return err
	}

	host, port, err := net.SplitHostPort(parsedURL.Host)
	if err != nil {
		return err
	}
	address := fmt.Sprintf("%s:%s", host, port)
	conn, err := net.DialTimeout("tcp", address, timeout)
	if err != nil {
		return err
	}

	defer conn.Close()
	return nil
}

func postChunk(address, data string) {
	b := strings.NewReader(data)
	resp, err := http.Post(address, "text/plain", b)
	panicOnErrorM(err)
	defer resp.Body.Close()
	// get stautus code
	fmt.Println("Status code:", resp.StatusCode)
	//timeout := time.Second * 2
	//
	//err := checkClientAvailability(address, timeout)
	//if err != nil {
	// fmt.Printf("Client at %s:%s is not available\n", address, err)
	// log.Fatal("Slave not available.")
	//} else {
	// fmt.Printf("Client at %s is available\n", address)
	// resp, err := http.Post(address, "text/plain", b)
	// panicOnErrorM(err)
	// defer resp.Body.Close()
	// // get stautus code
	// fmt.Println("Status code:", resp.StatusCode)
	//}
}
