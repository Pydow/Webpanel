package main

import (
	"fmt"
	"log"
	"net/http"
	. "webpanel/components/provider"
	. "webpanel/components/routes"
)

func main() {
	StartRedis()
	StartMongoDB()
	http.HandleFunc("/", Home)

	mux := http.NewServeMux()
	fs := http.FileServer(http.Dir("styles"))

	mux.Handle("/styles/", http.StripPrefix("/styles/", fs))

	mux.HandleFunc("/login", LoginPage)
	mux.HandleFunc("/handle/login", LoginHandler)

	mux.HandleFunc("/", Home)

	fmt.Println("Webserver gestartet.")
	log.Fatal(http.ListenAndServe(":8000", mux))
}


