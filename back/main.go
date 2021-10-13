package main

import (
	"fmt"
	"log"
	http "net/http"
)


func main() {
	hub := NewHub()
	go hub.run()

	log.Println("LETS GO :8080")

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		if r.Method == http.MethodOptions {
			return
		}
		serveWs(hub, w, r)
	})

	if err := http.ListenAndServe(fmt.Sprintf(":%d", 8080), nil); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
