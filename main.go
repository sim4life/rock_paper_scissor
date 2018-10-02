package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const NUMMOVES = 100

type Move int

const (
	Unknown Move = iota
	Rock
	Paper
	Scissors
)

func (m Move) String() string {
	return movesId[m]
}

var movesId = map[Move]string{
	Rock:     "rock",
	Paper:    "paper",
	Scissors: "scissors",
}

func (m Move) MarshalJSON() ([]byte, error) {
	return json.Marshal(movesId[m])
}

func (m *Move) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	switch strings.ToLower(s) {
	default:
		*m = Unknown
	case "rock":
		*m = Rock
	case "paper":
		*m = Paper
	case "scissors":
		*m = Scissors
	}

	return nil
}

type Decision string

const (
	Tie     Decision = "Tie"
	Player1          = "Player1"
	Player2          = "Player2"
)

type GameHistory []RoundHistory

type RoundHistory struct {
	Round  int             `json:"Round"`
	Winner Decision        `json:"Winner"`
	Inputs map[string]Move `json:"Inputs"`
}

func main() {

	gameHistory := runGameForMoves(NUMMOVES)

	currDir, err := filepath.Abs("./")
	if err != nil {
		log.Fatal(err)
	}
	filePath := filepath.Join(currDir, "rps.json")

	if err := saveGameIntoFile(gameHistory, filePath); err != nil {
		log.Println(err)
	}

	fmt.Printf("\nGame History saved to file:\n%s\n", filePath)
}

func runGameForMoves(numMoves int) GameHistory {
	gameHistory := make(GameHistory, 0)
	c1_movesHistory := make(chan GameHistory)
	c2_movesHistory := make(chan GameHistory)
	c1_move := getPlayer1Move(c1_movesHistory)
	c2_move := getPlayer2Move(c2_movesHistory)

	defer close(c1_movesHistory)
	defer close(c2_movesHistory)

	for i := 0; i < numMoves; i++ {
		c1_movesHistory <- gameHistory
		p1Move := <-c1_move
		c2_movesHistory <- gameHistory
		p2Move := <-c2_move

		winner := decideWinner(p1Move, p2Move)

		inputs := map[string]Move{
			Player1: p1Move,
			Player2: p2Move,
		}

		roundHistory := &RoundHistory{i + 1, winner, inputs}
		gameHistory = append(gameHistory, *roundHistory)
	}

	return gameHistory
}

func decideWinner(move1, move2 Move) Decision {
	if move1 == move2 {
		return Tie
	}
	if move1 == Rock && move2 == Paper {
		return Player2
	}
	if move1 == Rock && move2 == Scissors {
		return Player1
	}
	if move1 == Paper && move2 == Rock {
		return Player1
	}
	if move1 == Paper && move2 == Scissors {
		return Player2
	}
	if move1 == Scissors && move2 == Paper {
		return Player1
	}
	if move1 == Scissors && move2 == Rock {
		return Player2
	}
	return Tie
}

// Returns receive only channel of enum Move
func getPlayer1Move(c_movesHistory <-chan GameHistory) <-chan Move {
	c_move := make(chan Move)

	rand.Seed(time.Now().Unix())

	go func() {
		defer close(c_move)
		for range c_movesHistory {
			move := rand.Intn(4-1) + 1
			c_move <- Move(move)
		}
	}()

	return c_move
}

// Returns receive only channel of enum Move
func getPlayer2Move(c_movesHistory <-chan GameHistory) <-chan Move {
	c_move := make(chan Move)
	var lastRound RoundHistory
	p1Move := Rock //hard-coded first move

	go func() {
		defer close(c_move)
		for gameHistory := range c_movesHistory {
			if len(gameHistory) > 0 {
				lastRound = gameHistory[len(gameHistory)-1]
				p1Move = lastRound.Inputs[Player1]
			}
			c_move <- p1Move
		}
	}()

	return c_move
}

func getDataBytes(data GameHistory) ([]byte, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	return jsonData, nil
}

func PrettyJSON(uglyJSON []byte) (string, error) {
	var out bytes.Buffer
	err := json.Indent(&out, uglyJSON, "", "  ")
	return string(out.Bytes()), err
}

func saveGameIntoFile(gameHistoryData GameHistory, filePath string) error {
	fileHandle, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		log.Println(err)
		return err
	}
	defer fileHandle.Close()

	bytesData, err := getDataBytes(gameHistoryData)
	if err != nil {
		log.Println(err)
	}
	gameJSON, _ := PrettyJSON(bytesData)
	fmt.Printf("\n json is:\n%s\n", gameJSON)
	fileHandle.Write(bytesData)
	fileHandle.Sync()

	return nil
}
