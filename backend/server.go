package main

import (
	"fmt"
	"log"
	"net/http"
)

func foo(w http.ResponseWriter, req *http.Request) {
	log.Println("/foo invoked")
	fmt.Fprintf(w, "Hello from Foo!!\n")
}

func bar(w http.ResponseWriter, req *http.Request) {
	log.Println("/bar invoked")
	fmt.Fprintf(w, "Hello from Bar!!\n")
}

func main() {
	prodServer := http.NewServeMux()
	prodServer.HandleFunc("/foo", foo)

	sandServer := http.NewServeMux()
	sandServer.HandleFunc("/bar", bar)

	go func() { log.Fatal(http.ListenAndServe(":8001", prodServer)) }()
	go func() { log.Fatal(http.ListenAndServe(":8002", sandServer)) }()
	log.Println(">>> Prod and Sand servers started in :8001 & :8002")
	select {}
}
