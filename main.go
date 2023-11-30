package main

import (
	"fmt"
	"hangman-web/functions"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"text/template"
)

type GameInfo struct {
	AttemptLeft int
	WordHidden  []string
	Hangman     string
	LettersUsed []string
}

const port = ":8080"
const redirect = 301

var attemptLeft int
var WordToGuess string  // Mot (L A V A B O)
var WordHidden []string // Mot (_ A _ _ B O)
var bool = true

var LettersUsed []string
var Name string
var level string

func Home(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "tmpl/home.html")
}

func Rules(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "tmpl/rules.html")
}

func Leaderboard(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "tmpl/leaderboard.html")
}

func Profil(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "tmpl/profil.html")
}

func Game(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "tmpl/game.html")
}

func Game1(w http.ResponseWriter, r *http.Request) {
	if len(os.Args) == 1 {
		if err := r.ParseForm(); err != nil {
			fmt.Fprintf(w, "ParseForm() err: %v", err)
			return
		}

		level = r.Form.Get("difficulty")
		Name = r.Form.Get("Pseudo")
		WordToGuess = strings.ToUpper(functions.Random(level))
		LettersUsed = []string{}
		attemptLeft = 5

		//Mot caché (_ _ _ _ _)
		WordHidden = make([]string, len(WordToGuess))
		for i := range WordHidden {
			WordHidden[i] = "_"
		}

		//Premières lettres dévoilées (_ I _ _ E)
		NbRevealedLetters := len(WordToGuess)/2 - 1
		if NbRevealedLetters < 1 {
			NbRevealedLetters = 1
		}

		indexLetters := make([]int, 0)
		for len(indexLetters) < NbRevealedLetters {
			randomIndex := rand.Intn(len(WordToGuess))
			if !functions.Contains(indexLetters, randomIndex) {
				indexLetters = append(indexLetters, randomIndex)
			}
		}
		for _, i := range indexLetters {
			WordHidden[i] = string(WordToGuess[i])
		}
	}
	http.Redirect(w, r, "/", redirect)
}

func Game3(w http.ResponseWriter, r *http.Request) {

}

func Redirect(w http.ResponseWriter, r *http.Request) {
	if bool {
		bool = false
		http.Redirect(w, r, "/game", redirect)
	} else {
		new := GameInfo{AttemptLeft: attemptLeft, WordHidden: WordHidden, Hangman: functions.PrintHangman(attemptLeft), LettersUsed: LettersUsed}
		tmpl := template.Must(template.ParseFiles("tmpl/index.html"))
		tmpl.Execute(w, new)
	}
}

func main() {

	http.HandleFunc("/", Redirect)
	http.HandleFunc("/home", Home)
	http.HandleFunc("/rules", Rules)
	http.HandleFunc("/leaderboard", Leaderboard)
	http.HandleFunc("/game", Game)
	http.HandleFunc("/profil", Profil)
	http.HandleFunc("/game1", Game1)

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	http.Handle("/img/", http.StripPrefix("/img/", http.FileServer(http.Dir("./img"))))

	fmt.Println("\n(http://localhost:8080/home) - Server started on port", port)
	http.ListenAndServe(port, nil)
}
