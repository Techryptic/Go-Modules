package main

import (
    "fmt"
    "log"
    "io/ioutil"
    "net/http"
    "strconv"
    "strings"
    "os"
    "bufio"
)

func main() {
    //go hello()
    toLookFor :="jquery"

    file := os.Args[1]
    f, _ := os.Open(file)
    scanner := bufio.NewScanner(f)

    for scanner.Scan() {
        url := scanner.Text()
        fmt.Println(url)

    response, err := http.Get(url)
    if err != nil {
        fmt.Println("Argh! Broken")
        log.Fatal(err)
    }

    defer response.Body.Close()

    responseData, err := ioutil.ReadAll(response.Body)
    if err != nil {
        fmt.Println("Argh! Broken2")
        log.Fatal(err)
    }
    responseString := string(responseData)

    fmt.Println("HTTP Response Status: " + strconv.Itoa(response.StatusCode))

    if response.StatusCode >= 200 && response.StatusCode <= 299 {
        fmt.Println("HTTP Status OK!")
        if strings.Contains(responseString, toLookFor) {
            fmt.Printf("Found "+toLookFor+" in "+url+" \n")
        }
    } else {
        fmt.Println("Argh! Broken")
    }
    }
}

//func hello() {  
//    fmt.Println("Hi, I'm a goroutine!")
//}

