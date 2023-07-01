package main

import (
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/fbuys/myurls/internal/myurls"
)

func main() {
	for _, url := range myurls.GetAllUrls() {
		http.HandleFunc(url.Id, redirect(url.Address))
	}

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

func redirect(address string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		http.Redirect(w, req, address, 302)
	}
}
