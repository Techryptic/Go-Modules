package main

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

func main() {

	creds := [3]string{ //Need to put exactly how many creds are in the list below
		"admin:123456",
		"user:password",
		"root:toor"}

	// Command line params
	path := flag.String("l", "", "File containing IP Addresses (one per line)")
	singleIP := flag.String("s", "", "Single IP Address")
	nJobs := flag.Int("n", 1, "Number of concurrent requests to make")
	flag.Parse()

	// Checking command line params
	if *path == "" && *singleIP == "" {
		fmt.Println("Options must be supplied.")
		flag.Usage()
		os.Exit(1)
	} else {

		// Setup for timing the function
		startTime := time.Now()

		// Open the file, parse the lines
		addresses, err := readFile(*path)
		if err != nil {
			fmt.Printf("Failed to read file: %v\n", err)
			os.Exit(2)
		}

		// Create a channel to pass the ip addresses into
		//ipChan := make(chan string)
		// TODO: If you need to do specific results, you can create a second channel to pass them back with
		// resChan := make(chan resultStruct) // Where `type resultStruct struct { /* */ }`

		//go // Create a channel to pass the ip addresses into struct requestCombo
		type requestCombo struct {
			ipAddress string
			credss    string
		}

		inChan := make(chan requestCombo)

		// Start up worker goroutines
		wg := &sync.WaitGroup{}
		wg.Add(*nJobs)
		for i := 0; i < *nJobs; i++ {
			go func() {
				for combo := range inChan {
					//fmt.Println(combo.ipAddress)
					//fmt.Println(combo.credss)
					address := combo.ipAddress
					credz := combo.credss
					err := FirstRequest(address, credz)
					if err != nil {
						//fmt.Printf("Failed to FirstRequest (%s): %v\n", address, err)
					}
				}
				wg.Done()
			}()
		}
		// Start feeding the worker goroutines
		for _, cred := range creds {
			for _, ip := range addresses {
				inChan <- requestCombo{
					ipAddress: ip,
					credss:    cred,
				}
			}
		}
		close(inChan)
		wg.Wait()

		// Report to user.
		totalTime := time.Since(startTime)
		fmt.Printf("Finished processing %d IPs with %d goroutines in %v.\n", len(addresses), *nJobs, totalTime)

	}
}

func readFile(path string) ([]string, error) {
	// Note: this function will read the full contents of the file into memory.
	// If it is an issue, use a buffered channel instead of returning a string array.

	var array []string

	// Open the file
	f, err := os.Open(path)
	if err != nil {
		return array, err
	}

	// Read it line by line
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := string(scanner.Text())
		array = append(array, line)
	}

	return array, nil
}

func getPageContent(port int, url, address, Protocol, HTTP_Method, req_BodyText, h1, v1, h2, v2 string) (*http.Response, error) {
	//This bit will take first priority of the protocol supplied in the argument before the one set in the request
	if !strings.Contains(address, "://") {
		address = (Protocol + "://" + address)
	}
	fullurl := strings.Join([]string{address, ":", strconv.Itoa(port)}, url, "") //Combinding Protocol:ip:port
	var query = []byte(req_BodyText)                                    //Request Body Text
	req, err := http.NewRequest(strings.ToUpper(HTTP_Method), fullurl, bytes.NewBuffer(query))
	req.Close = true
	//req.Header.Set("Cookies", "text/plain")   //Static Request Header
	req.Header.Set("User-Agent", "Firefox") //Static Request Header
	req.Header.Set(h1, v1)
	req.Header.Set(h2, v2)

	tr := &http.Transport{
		Proxy:              http.ProxyFromEnvironment,
		MaxIdleConns:       10,
		IdleConnTimeout:    30 * time.Second,
		DisableCompression: false,
		TLSClientConfig: &tls.Config{
			//CipherSuites:             []uint16{tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA, tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA},
			PreferServerCipherSuites: true,
			InsecureSkipVerify:       true,
			MinVersion:               tls.VersionTLS10,
			MaxVersion:               tls.VersionTLS12,
		},
	}
	client := http.Client{
		Transport: tr,
		Timeout:   time.Duration(5 * time.Second)}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func FirstRequest(address, credz string) error {
    	credz := strings.Split(credz, ":")
   	_ = credz[0] //user
	_ = credz[1] //pass

	Protocol := "http" //http or https
	port := 80
	url := "/"
	HTTP_Method := "get"
	req_BodyText := "Body Text"
	h1, v1 := "Accept-Encoding", "gzip, deflate" //Request Header & Value
	h2, v2 := "Content-Type", "text/plain"       //Request Header & Value

	resp, err := getPageContent(port, url, address, Protocol, HTTP_Method, req_BodyText, h1, v1, h2, v2) //Sending Contents to getPageContent func
	if err != nil {
		//log.Fatal(err)
		//fmt.Printf("Server not responding %s\n", err.Error()) //Catch-all errors
	} else {
		//fmt.Printf("(%s) returned status: %s\n", address, resp.Status) //Prints Successful Connections
		defer resp.Body.Close()
		//respBody, err := ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			//log.Fatal(err)
		}
		//resp_Body := string(respBody[:])   // Response Body, to print: //fmt.Println(resp_Body)
		resp_StatusCode := resp.StatusCode // 200 / 500 ect
		/*
			resp_Server := resp.Header.Get("Server")
			resp_ContentType := resp.Header.Get("Content-Type")
			resp_ContentLength := resp.Header.Get("Content-Length")
			resp_SetCookie := resp.Header.Get("Set-Cookie")
			resp_Location := resp.Header.Get("Location")
			resp_XPoweredBy := resp.Header.Get("X-Powered-By")
		*/

		/* Example of Carrying over a Cookie to the SecondRequest function
		resp_XDISRequestID := resp.Header.Get("X-DIS-Request-ID") 	//Setting cookie value to the parameter resp_XDISRequestID
		fmt.Println(resp_XDISRequestID) 							//Printing cookie for debug purposes
		Check := strings.Contains(resp_Body, "cookie") 				//See if the String 'cookie' is inside the request body.

		if Check == true && resp_StatusCode == 200 {				//Runs the check..
			fmt.Println("Default Creds Found!")       				//Printing Results
			SecondRequest(resp_XDISRequestID, address) 				//Starting Second Request, will carry over the response cookie + address
		}
		*/

		//Strings to check inside Response Body
		//Check1 := strings.Contains(resp_Body, "administrator") //String Contains
		//Check2 := !strings.Contains(resp_Body, "logged out")   //String does not Contain

		//Checking for match
		if resp_StatusCode == 200 {
			fmt.Println("Default Creds Found!" + credz + address)
			//resp_XDISRequestID := "tech"
			//SecondRequest(resp_XDISRequestID, address)
		}

		/* For debuging :)
		//os.Setenv("http_proxy", "http://127.0.0.1:8080")
		//os.Setenv("https_proxy", "https://127.0.0.1:8080")
		//os.Setenv("HTTP_PROXY", "http://127.0.0.1:8080")
		//os.Setenv("HTTPS_PROXY", "https://127.0.0.1:8080")

		//Print all headers
		for k, v := range resp.Header {
			fmt.Print(k)
			fmt.Print(" : ")
			fmt.Println(v)
		}
		*/
	}
	return nil
}

func SecondRequest(Cookie, address string) {

	Protocol := "http" //http or https
	port := 80
	HTTP_Method := "get"
	req_BodyText := "Body Text"
	h1, v1 := "Accept-Encoding", "gzip, deflate" //Request Header & Value
	h2, v2 := "Content-Type", "text/plain"       //Request Header & Value

	resp, err := getPageContent(port, address, Protocol, HTTP_Method, req_BodyText, h1, v1, h2, v2) //Sending Contents to getPageContent func
	if err != nil {
		//log.Fatal(err)
		fmt.Printf("Server not responding %s\n", err.Error()) //Catch-all errors
	} else {
		fmt.Printf("(%s) returned status: %s\n", address, resp.Status) //Prints Successful Connections
		defer resp.Body.Close()
		respBody, err := ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			log.Fatal(err)
		}
		resp_Body := string(respBody[:])   // Response Body, to print: //fmt.Println(resp_Body)
		resp_StatusCode := resp.StatusCode // 200 / 500 ect

		Check := strings.Contains(resp_Body, "cookie") //String Contains

		if Check == true && resp_StatusCode == 200 {
			fmt.Println("Default Creds Found!") //Printing Results
		}
	}
}
