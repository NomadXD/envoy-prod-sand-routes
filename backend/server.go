package main

import (
	"fmt"
	"log"
	"net/http"
)

func foo(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "Hello from Foo!!\n")
}

func bar(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "Hello from Bar!!\n")
}

func main() {
	prodServer := http.NewServeMux()
	prodServer.HandleFunc("/foo", foo)

	sandServer := http.NewServeMux()
	sandServer.HandleFunc("/bar", bar)

	go func() { log.Fatal(http.ListenAndServe(":8001", prodServer)) }()
	go func() { log.Fatal(http.ListenAndServe(":8002", sandServer)) }()
	log.Println(">>> Prod and Sand servers started")
	select {}
}
