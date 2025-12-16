package main

type GameInstance struct {
	Bounds Coord
	Board  GameBoard
	Snakes []Battlesnake
	Food   []Coord
	You    Battlesnake
}

type GameBoard map[Coord]CoordState

func NewGameInstace(gs *GameState) GameInstance {
	return GameInstance{
		Bounds: Coord{gs.Board.Width, gs.Board.Height},
		Board:  NewGameBoard(gs),
		Snakes: gs.Board.Snakes,
		Food:   gs.Board.Food,
		You:    gs.You,
	}
}

func NewGameBoard(gs *GameState) GameBoard {
	gb := make(GameBoard, gs.Board.Height*gs.Board.Width)
	for _, f := range gs.Board.Food {
		gb[f] = CoordFood
	}
	for _, s := range gs.Board.Snakes {
		for _, b := range s.Body {
			gb[b] = CoordSnake
		}
	}
	return gb
}
