package main

import (
	"database/sql"
	"fmt"
	"hangman-web/functions"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"text/template"

	_ "github.com/mattn/go-sqlite3"
)

// Struct of the game infos
type GameInfo struct {
	AttemptLeft   int
	HiddenWord    []string
	Hangman       string
	LettersUsed   []string
	WordToGuess   string
	Score         int
	ConnectedUser string
}

//struct of the user infos
type UserInfo struct {
	Id       int
	Username string
	Password string
	IsAdmin  bool
	IsBanned bool
}

//struct for the score
type Score struct {
	Username string
	Score    int
}

//constant for the redirect
const port = ":8080"
const redirect = 301

//variables for the game
var attemptLeft int
var WordToGuess string  // Mot (L A V A B O)
var HiddenWord []string // Mot (_ A _ _ B O)
var LettersUsed []string
var level string

//variable to stock the connected user
var connectedUser string

//home page
func Home(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "tmpl/home.html")
}

//rules page
func Rules(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "tmpl/rules.html")
}

//function to handle the leaderboard
func Leaderboard(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	//using the function to get the 10th first score of all users
	topScores, err := GetTopScores(db, 10)
	if err != nil {
		http.Error(w, "Error retrieving top scores", http.StatusInternalServerError)
		return
	}

	//get the 10th first score of the connected user
	var userTopScores []Score
	if connectedUser != "" {
		userTopScores, err = GetUserTopScores(db, connectedUser, 10)
		if err != nil {
			http.Error(w, "Error retrieving user top scores", http.StatusInternalServerError)
			return
		}
	}

	//creat a template to print the leaderboard
	tmpl := template.Must(template.ParseFiles("tmpl/leaderboard.html"))
	data := struct {
		TopScores     []Score
		UserTopScores []Score
		ConnectedUser string
	}{
		TopScores:     topScores,
		UserTopScores: userTopScores,
		ConnectedUser: connectedUser,
	}
	tmpl.Execute(w, data)
}

//function to handle the profil
func Profil(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	//get if the user connected is admin or is banned
	var isAdmin, isBanned bool
	if connectedUser != "" {
		isAdmin, isBanned = IsAdminOrIsBanned(db)
	}
	//creat a template with the infos we get
	new := UserInfo{Username: connectedUser, IsAdmin: isAdmin, IsBanned: isBanned}
	tmpl := template.Must(template.ParseFiles("tmpl/profil.html"))
	tmpl.Execute(w, new)
}

//function to handle the start of the game
func Game(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "tmpl/game.html")
}

//function to initialize the game infos
func Game1(w http.ResponseWriter, r *http.Request) {
	//we check if there is only one argument
	if len(os.Args) == 1 {
		//we check if there is an error when ParseForm analyze the data of the http request
		if err := r.ParseForm(); err != nil {
			fmt.Fprintf(w, "ParseForm() err: %v", err)
			return
		}

		//we get the value of difficulty chose by the user
		level = r.Form.Get("difficulty")
		//and we initialize the word to guess depending of the difficulty
		WordToGuess = strings.ToUpper(functions.Random(level))
		LettersUsed = []string{}
		attemptLeft = 10

		//Mot caché (_ _ _ _ _)
		HiddenWord = make([]string, len(WordToGuess))
		for i := range HiddenWord {
			HiddenWord[i] = "_"
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
			HiddenWord[i] = string(WordToGuess[i])
		}
	}
	//we redirect to the game loop
	http.Redirect(w, r, "/", redirect)
}

//function that checks if the entered letter has been used
func Game2(w http.ResponseWriter, r *http.Request) {
	//we check if there is an error when ParseForm analyze the data of the http request
	if err := r.ParseForm(); err != nil {
		fmt.Fprintf(w, "ParseForm() err: %v", err)
		return
	}

	//here we check if the letter has been used
	letter := strings.ToUpper(r.Form.Get("try"))
	alreadyIn := false
	for i := range LettersUsed {
		if LettersUsed[i] == letter {
			alreadyIn = true
		}
	}

	if !alreadyIn {
		HiddenWord, attemptLeft, LettersUsed = functions.Compare(HiddenWord, attemptLeft, WordToGuess, letter, LettersUsed[:])
	}

	if attemptLeft <= 0 && r.Method == "POST" {
		http.Redirect(w, r, "/lose", redirect)
	}

	HiddenWordStr := strings.Join(HiddenWord, "")
	if HiddenWordStr == WordToGuess && r.Method == "POST" {
		http.Redirect(w, r, "/win", redirect)
	}

	//we redirect to the game loop
	http.Redirect(w, r, "/", redirect)
}

//function that handle the win page and calculate the score
func Win(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	//we calculate the score depending of the attempt left and the difficulty chose
	Score := functions.Score(attemptLeft, level)
	if connectedUser != "" {
		//add the score to the db
		_, err := db.Exec("INSERT INTO scores (username, score) VALUES (?, ?)", connectedUser, Score)
		if err != nil {
			panic(err.Error())
		}
	}
	//creat a template with the score and the word to guess
	new := GameInfo{Score: Score, WordToGuess: WordToGuess}
	tmpl := template.Must(template.ParseFiles("tmpl/win.html"))
	tmpl.Execute(w, new)
}

//function that handle the lose page
func Lose(w http.ResponseWriter, r *http.Request) {
	//creat a template with the word to guess
	new := GameInfo{WordToGuess: WordToGuess}
	tmpl := template.Must(template.ParseFiles("tmpl/lose.html"))
	tmpl.Execute(w, new)
}

//function that handle the login
func Login(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	//we check if the method is a POST
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	//we get the username and the password
	username := r.FormValue("login-username")
	password := r.FormValue("login-password")

	//request SQL to get the user
	row := db.QueryRow("SELECT * FROM users WHERE username=? AND password=?", username, password)
	var storedUsername, storedPassword string
	var isAdmin, isBanned bool
	var id int

	//scan and get the data
	err := row.Scan(&id, &storedUsername, &storedPassword, &isAdmin, &isBanned)

	//we compare the data
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Fprintln(w, "User not found")
		} else {
			http.Error(w, fmt.Sprintf("Error: %v", err), http.StatusInternalServerError)
		}
		return
	}
	//change the user connected
	connectedUser = username
	//redirect to profil
	http.Redirect(w, r, "/profil", redirect)
}

//function that handle the sign up
func Signup(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	//we check if the method is a POST
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	//we get the username and the password
	username := r.FormValue("username")
	password := r.FormValue("password")

	//request SQL to check if the user already exist
	row := db.QueryRow("SELECT * FROM users WHERE username=?", username)
	var storedUsername, storedPassword string
	var isAdmin, isBanned bool
	var id int
	//we scan to get the data
	err := row.Scan(&id, &storedUsername, &storedPassword, &isAdmin, &isBanned)

	//we check if the username is already used
	if err == nil {
		fmt.Fprintln(w, "This username already exist, please select another one")
		return
	} else if err != sql.ErrNoRows {
		http.Error(w, fmt.Sprintf("Error: %v", err), http.StatusInternalServerError)
		return
	}

	//we creat a new user in the db
	_, err = db.Exec("INSERT INTO users (username, password, isAdmin, isBanned) VALUES (?, ?, FALSE, FALSE)", username, password)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error: %v", err), http.StatusInternalServerError)
		return
	}
	//change the user connected
	connectedUser = username
	//redirect to profil
	http.Redirect(w, r, "/profil", redirect)
}

//function to handle the logout
func Logout(w http.ResponseWriter, r *http.Request) {
	//we reset the variable
	connectedUser = ""
	//we redirect to profil
	http.Redirect(w, r, "/profil", redirect)
}

//function that handle the game loop
func Redirect(w http.ResponseWriter, r *http.Request) {
	//creat a template with all the game infos needed
	new := GameInfo{AttemptLeft: attemptLeft, HiddenWord: HiddenWord, Hangman: functions.PrintHangman(attemptLeft), LettersUsed: LettersUsed}
	tmpl := template.Must(template.ParseFiles("tmpl/index.html"))
	tmpl.Execute(w, new)
}

//function that handle the Admin panel
func Admin(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	//we check if the user is admin
	isAdmin, _ := IsAdminOrIsBanned(db)

	if !isAdmin {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	//we get all user from the db
	users, err := GetUsersFromDatabase(db)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	//creat a struct to return
	data := struct {
		Users []UserInfo
	}{
		Users: users,
	}

	//creat a template with the struct of data
	tmpl := template.Must(template.ParseFiles("tmpl/admin.html"))
	tmpl.Execute(w, data)
}

//function that creat a table User
func createTableUser(db *sql.DB) {
	//creating the user table if not already created
	_, err := db.Exec(`
        CREATE TABLE IF NOT EXISTS users (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            username TEXT NOT NULL,
            password TEXT NOT NULL,
			isAdmin BOOL NOT NULL,
			isBanned BOOL NOT NULL
        )
    `)
	if err != nil {
		panic(err.Error())
	}
}

//function that creat a table Score
func createTableScores(db *sql.DB) {
	//creating the score table if not already created
	_, err := db.Exec(`
        CREATE TABLE IF NOT EXISTS scores (
            username TEXT NOT NULL,
            score INT NOT NULL
        )
    `)
	if err != nil {
		panic(err.Error())
	}
}

//function that get the top score for the leaderboard
func GetTopScores(db *sql.DB, limit int) ([]Score, error) {
	//get the user and the score by order to get the higher scores
	rows, err := db.Query("SELECT username, score FROM scores ORDER BY score DESC LIMIT ?", limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var topScores []Score
	//scan each one by one
	for rows.Next() {
		var username string
		var score int
		//scan the username and the score
		if err := rows.Scan(&username, &score); err != nil {
			return nil, err
		}
		//add it to the table
		topScores = append(topScores, Score{Username: username, Score: score})
	}

	return topScores, nil
}

//function that get the top score of the user for the leaderboard
func GetUserTopScores(db *sql.DB, username string, limit int) ([]Score, error) {
	//get the user and the score by order to get the higher scores of one specific user
	rows, err := db.Query("SELECT username, score FROM scores WHERE username = ? ORDER BY score DESC LIMIT ?", username, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var userTopScores []Score
	//scan each one by one
	for rows.Next() {
		var score int
		//scan the username and the score
		if err := rows.Scan(&username, &score); err != nil {
			return nil, err
		}
		//add it to the table
		userTopScores = append(userTopScores, Score{Username: username, Score: score})
	}

	return userTopScores, nil
}

//function that check if the user is Admin or Banned
func IsAdminOrIsBanned(db *sql.DB) (bool, bool) {
	//request SQL
	row := db.QueryRow("SELECT IsAdmin, IsBanned FROM users WHERE username=?", connectedUser)
	var isAdmin, isBanned bool
	//scan the data
	err := row.Scan(&isAdmin, &isBanned)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
	// return the data
	return isAdmin, isBanned
}

//function that get the infos of the user
func GetUsersFromDatabase(db *sql.DB) ([]UserInfo, error) {
	//request SQL
	rows, err := db.Query("SELECT Id, Username, Password, IsAdmin, IsBanned FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []UserInfo

	//scan each one by one
	for rows.Next() {
		var user UserInfo
		//scan the data
		err := rows.Scan(&user.Id, &user.Username, &user.Password, &user.IsAdmin, &user.IsBanned)
		if err != nil {
			return nil, err
		}
		//add the user
		users = append(users, user)
	}

	return users, nil
}

//function that remove the admin
func RemoveAdmin(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	userIdStr := r.PostForm.Get("userId")

	userId, err := strconv.Atoi(userIdStr)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	//use that function update the admin status
	err = UpdateUserAdminStatus(db, userId, false)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	//redirect to the admin panel
	http.Redirect(w, r, "/admin", redirect)
}

//function that grant the admin
func GrantAdmin(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	userIdStr := r.PostForm.Get("userId")

	userId, err := strconv.Atoi(userIdStr)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	//use that function update the admin status
	err = UpdateUserAdminStatus(db, userId, true)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	//redirect to the admin panel
	http.Redirect(w, r, "/admin", redirect)
}

//function that update the admin status
func UpdateUserAdminStatus(db *sql.DB, userId int, makeAdmin bool) error {
	//request SQL
	query := "UPDATE users SET IsAdmin = ? WHERE Id = ?"

	//exect the request with the parameters
	_, err := db.Exec(query, makeAdmin, userId)
	return err
}

//function that ban the user
func Ban(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	userIdStr := r.PostForm.Get("userId")

	userId, err := strconv.Atoi(userIdStr)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	////use that function update the ban status
	err = UpdateUserBanStatus(db, userId, true)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	//redirect to the admin panel
	http.Redirect(w, r, "/admin", redirect)
}

//function that unban the user
func Unban(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	userIdStr := r.PostForm.Get("userId")

	userId, err := strconv.Atoi(userIdStr)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	////use that function update the ban status
	err = UpdateUserBanStatus(db, userId, false)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	//redirect to the admin panel
	http.Redirect(w, r, "/admin", redirect)
}

func UpdateUserBanStatus(db *sql.DB, userId int, makeAdmin bool) error {
	query := "UPDATE users SET IsBanned = ? WHERE Id = ?"

	_, err := db.Exec(query, makeAdmin, userId)
	return err
}

func main() {
	//open the database with sqlite3
	db, err := sql.Open("sqlite3", "database.db")
	if err != nil {
		panic(err.Error())
	}
	//creat the 2 tables
	createTableUser(db)
	createTableScores(db)
	defer db.Close()

	//handle the request
	http.HandleFunc("/", Redirect)
	http.HandleFunc("/home", Home)
	http.HandleFunc("/rules", Rules)
	http.HandleFunc("/game", Game)
	http.HandleFunc("/game1", Game1)
	http.HandleFunc("/game2", Game2)
	http.HandleFunc("/lose", Lose)
	http.HandleFunc("/logout", Logout)
	http.HandleFunc("/profil", func(w http.ResponseWriter, r *http.Request) { Profil(w, r, db) })
	http.HandleFunc("/leaderboard", func(w http.ResponseWriter, r *http.Request) { Leaderboard(w, r, db) })
	http.HandleFunc("/win", func(w http.ResponseWriter, r *http.Request) { Win(w, r, db) })
	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) { Login(w, r, db) })
	http.HandleFunc("/signup", func(w http.ResponseWriter, r *http.Request) { Signup(w, r, db) })
	http.HandleFunc("/admin", func(w http.ResponseWriter, r *http.Request) { Admin(w, r, db) })
	http.HandleFunc("/removeAdmin", func(w http.ResponseWriter, r *http.Request) { RemoveAdmin(w, r, db) })
	http.HandleFunc("/grantAdmin", func(w http.ResponseWriter, r *http.Request) { GrantAdmin(w, r, db) })
	http.HandleFunc("/ban", func(w http.ResponseWriter, r *http.Request) { Ban(w, r, db) })
	http.HandleFunc("/unban", func(w http.ResponseWriter, r *http.Request) { Unban(w, r, db) })

	//load the CSS and the images
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	http.Handle("/img/", http.StripPrefix("/img/", http.FileServer(http.Dir("./img"))))

	//start the local host
	fmt.Println("\n(http://localhost:8080/home) - Server started on port", port)
	http.ListenAndServe(port, nil)
}