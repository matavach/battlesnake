package main

type GameBoard map[Coord]CoordState

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
