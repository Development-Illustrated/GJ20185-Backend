package main

import (
    "bufio"
    "fmt"
	"os"
    // "net"
    "net/http"
    "net/url"
)

func main() {
    readFromTerminal();
}

func readFromTerminal() {
	buf := bufio.NewReader(os.Stdin)
    fmt.Print("> ")
    sentence, err := buf.ReadBytes('\n')
    if err != nil {
        fmt.Println(err)
	// } else if 1 == 1{
	// 	fmt.Println("Hi Umair");
	} else {
        fmt.Println(string(sentence))
        sendMessageViaTCP("Up", "http://127.0.0.1", "8080")
    }
}

func sendMessageViaTCP(message string, backendAddr string, port string) {
    
    // http.PostForm("http://127.0.0.1:8080/",
    // url.Values{"key": {"Value"}, "id": {"123"}})

    http.PostForm(backendAddr+":"+port+"/",
    url.Values{"key": {"Value"}, "id": {"123"}})
    

    http.PostForm(backendAddr+":"+port+"/TodoShow",
    url.Values{"key": {"Value"})

    // fmt.Println(resp)
    // fmt.Println(err)
}