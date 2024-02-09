# Hangman Web

Hangman Web is a simple web-based Hangman game implemented in Go. It allows users to play Hangman, view the leaderboard, and manage user profiles. The game features different difficulty levels and keeps track of user scores.

## Installation

1. Clone the repository: `git clone https://ytrack.learn.ynov.com/git/fthibaul/HANGMAN-WEB.git`

2. Navigate to the project directory: `cd hangman-web`

3. Install dependencies: `go get -u github.com/mattn/go-sqlite3`

4. Run the application: `go run main.go`

   The server will start at http://localhost:8080/home.

## Features

- **Home Page:** Provides an introduction to the game.
- **Rules Page:** Displays the rules of the Hangman game.
- **Game Page:** Allows users to play Hangman with different difficulty levels.
- **Leaderboard:** Shows the top scores of all users and the user's personal top scores.
- **Profile Page:** Allows the user to login or to create an account to keep track of his best scores.
- **Admin Panel:** Accessible only by admins, this panel allows admins to view and manage user accounts.

## Usage

1. Access the home page at [http://localhost:8080/home](http://localhost:8080/home).
2. Navigate through the different pages:
   - **Rules:** Learn how to play Hangman.
   - **Game:** Play the Hangman game with various difficulty levels.
   - **Leaderboard:** View the top scores.
   - **Profile:** Check your user profile.
   - **Admin Panel:** Manage user accounts if you have admin privileges.
3. Enjoy playing Hangman and compete for high scores!

## Dependencies

- [github.com/mattn/go-sqlite3](https://github.com/mattn/go-sqlite3): SQLite database driver for Go.

- [winlibs.com/#download-release](https://winlibs.com/#download-release): Intall GCC

## Acknowledgments

This Hangman Web project was created by a group of students from Ynov.