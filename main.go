package main

import (
	"fmt"
	"net/http"
)

const port = ":8080"

func Home(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "tmpl/home.html")
}

func Rules(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "tmpl/rules.html")
}

func Leaderboard(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "tmpl/leaderboard.html")
}

func Game(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "tmpl/game.html")
}

func Profil(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "tmpl/profil.html")
}

func main() {

	http.HandleFunc("/home", Home)
	http.HandleFunc("/rules", Rules)
	http.HandleFunc("/leaderboard", Leaderboard)
	http.HandleFunc("/game", Game)
	http.HandleFunc("/profil", Profil)

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	http.Handle("/img/", http.StripPrefix("/img/", http.FileServer(http.Dir("./img"))))

	fmt.Println("\n(http://localhost:8080/home) - Server started on port", port)
	http.ListenAndServe(port, nil)
}
