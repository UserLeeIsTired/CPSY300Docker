package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.Handle("/", http.FileServer(http.Dir("./student")))

	port := ":3000"
	fmt.Printf("Server running on http://localhost%s\n", port)
	http.ListenAndServe(port, nil)
}
