package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/chrisfregly/tictactoe"
)

var game *tictactoe.TicTacToe

type tictactoeboard [3][3]*tictactoe.Player

var mu sync.Mutex

type GameState struct {
	Turn    tictactoe.Player  `json:"turn"`
	Winner  *tictactoe.Player `json:"winner,omitempty"`
	Board   tictactoeboard    `json:"board"`
	GameEnd bool              `json:"game_end"`
}

func handleGetGame(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()

	if game == nil {
		gameVal := tictactoe.NewTicTacToe()
		game = &gameVal
	}

	state := GameState{
		Turn:    game.GetTurn(),
		Winner:  game.GetWinner(),
		Board:   game.GetBoard(),
		GameEnd: game.IsGameOver(),
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(state); err != nil {
		http.Error(w, "Invalid JSON format", http.StatusInternalServerError)
		return
	}
}

func handlePostMove(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()

	if game == nil {
		http.Error(w, "Game already ended or not started. Please start a new game.", http.StatusBadRequest)
		return
	}

	var move struct {
		Player tictactoe.Player `json:"player"`
		Row    int              `json:"row"`
		Column int              `json:"column"`
	}

	err := json.NewDecoder(r.Body).Decode(&move)
	if err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	if game.IsGameOver() {
		http.Error(w, "Game already ended. Please reset the game.", http.StatusBadRequest)
		return
	}

	if move.Player != game.GetTurn() {
		http.Error(w, "Not the player's turn", http.StatusBadRequest)
		return
	}

	if err := game.Move(move.Player, move.Row, move.Column); err != nil {
		http.Error(w, fmt.Sprintf("Invalid move: %v", err), http.StatusBadRequest)
		return
	}

	state := GameState{
		Turn:    game.GetTurn(),
		Winner:  game.GetWinner(),
		Board:   game.GetBoard(),
		GameEnd: game.IsGameOver(),
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(state); err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}
}

func handleDeleteAndGet(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		handleGetGame(w, r)
	case http.MethodDelete:
		handleDeleteGame(w, r)
	}
}

func handleDeleteGame(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()

	game = nil

	w.WriteHeader(http.StatusOK)
}

func main() {

	http.HandleFunc("/game", handleDeleteAndGet)
	http.HandleFunc("/game/move", handlePostMove)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
