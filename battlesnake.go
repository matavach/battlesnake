package main

import "log"

type Config struct {
	Port       string
	MaxWorkers int8
}

var turn int

func main() {
	turn = 0
	c := Config{
		Port: "8000",
	}
	Run(&c)
}

func info() BattlesnakeInfoResponse {
	log.Printf("Info Requested\n")
	return BattlesnakeInfoResponse{
		APIVersion: "1",
		Author:     "ematijevich",
		Color:      "#808080",
		Head:       "default",
		Tail:       "default",
	}
}

func start(gs *GameState) {
	log.Println("GAME START")
}

func end(gs *GameState) {
	turn = 0
	log.Println("GAME OVER")
}

func move(gs *GameState) BattlesnakeMoveResponse {
	log.Printf("TURN %d:\n", gs.Turn)
	response := BattlesnakeMoveResponse{}
	gi := NewGameInstace(gs)
	snakey, _ := NewSnake(&gi.You)
	snakey.findFirstLevelMoves(&gi)
	snakey.findSecondLevelMoves(&gi)
	response.Move = snakey.NextMove
	log.Printf("MOVING: %s\n", response.Move)
	return response
}
