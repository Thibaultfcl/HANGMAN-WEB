package main

import (
	"fmt"
	"net/http"
)

const port = ":8080"

func Home(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "tmpl/home.html")
}

func main() {

	http.HandleFunc("/home", Home)

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	http.Handle("/img/", http.StripPrefix("/img/", http.FileServer(http.Dir("./img"))))

	fmt.Println("\n(http://localhost:8080/home) - Server started on port", port)
	http.ListenAndServe(port, nil)
}
