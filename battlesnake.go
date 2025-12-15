package main

import "log"

type Config struct {
	Port       string
	MaxWorkers int8
}

var snakey = snake{}
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

	snakey, _ = NewSnake(&gs.You)
}

func end(gs *GameState) {
	turn = 0
	log.Println("GAME OVER")
}

func move(gs *GameState) BattlesnakeMoveResponse {
	log.Printf("TURN %d:\n", gs.Turn)
	response := BattlesnakeMoveResponse{}
	snakey.UpdateSnake(&gs.You)
	snakey.findFirstLevelMoves(gs)
	snakey.findSecondLevelMoves(gs)
	response.Move = snakey.NextMove
	log.Printf("MOVING: %s\n", response.Move)
	return response
}
