package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

type Response struct {
	Id        string
	IpAddress []string
}

func main() {
	// GET request

	fmt.Println("Get from Master")
	resp, err := http.Get("http://localhost:8011/fasta?id=2")
	panicOnError(err)
	defer resp.Body.Close()
	// get stautus code
	fmt.Println("Status code:", resp.StatusCode)
	b, err := io.ReadAll(resp.Body)
	slaveIP := string(b)
	var res Response

	if er := json.Unmarshal([]byte(slaveIP), &res); er != nil {
		fmt.Println(er)
	}
	fmt.Println("slaveIP: ", res)

	for _, ip := range res.IpAddress {
		if strings.Contains(slaveIP, "http://") {
			fmt.Println("Get from Slave:" + ip)

			resp, err = http.Get(ip)
			panicOnError(err)
			defer resp.Body.Close()
			fmt.Println("Status code:", resp.StatusCode)
			b, err = io.ReadAll(resp.Body)
			if res.Id == "0" {

				f, err := os.OpenFile("client/count.fasta", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
				if err != nil {
					panicOnError(err)
				}
				defer f.Close()
				if _, err = f.WriteString(string(b)); err != nil {
					panicOnError(err)
				}
			} else {
				err = os.WriteFile("client/count.fasta", b, 0644)
				panicOnError(err)
			}
		} else {
			fmt.Println("Master Returned :", slaveIP)
		}
	}

}

func panicOnError(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
