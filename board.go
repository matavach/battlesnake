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
		Snakes: GetEnemySnakes(gs),
		Food:   gs.Board.Food,
		You:    gs.You,
	}
}

func GetEnemySnakes(gs *GameState) []Battlesnake {
	snakes := make([]Battlesnake, 0, len(gs.Board.Snakes)-1)
	for _, s := range gs.Board.Snakes {
		if s.ID == gs.You.ID {
			continue
		}
		snakes = append(snakes, s)
	}
	return snakes
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
