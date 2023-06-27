package main

import (
	"errors"
	"fmt"
	"net/http"
	"os"
)

func main() {
	http.HandleFunc("/", redirect)

	startServer()
}

func startServer() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Println("starting server on port", port)
	err := http.ListenAndServe(":"+port, nil)

	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("server closed\n")
	} else if err != nil {
		fmt.Printf("error starting server: %s\n", err)
		os.Exit(1)
	}
}

func redirect(w http.ResponseWriter, req *http.Request) {
	http.Redirect(w, req, "https://fbuys.dev/", 302)
}
