package main

type snake struct {
	Data      Battlesnake
	SafeMoves Directions
	Head      Coord
	Neck      Coord
}

var SafeMoves = []MoveCheck{
	s.MoveBackwards(),
}

func NewSnake(bs Battlesnake) (*snake, error) {
	s := &snake{
		Data: bs,
		SafeMoves: Directions{
			up:    true,
			left:  true,
			right: true,
			down:  true,
		},
		Head: bs.Head,
		Neck: bs.Body[1],
	}
	return s, nil
}

func (s *snake) FindSafeMoves(gs *GameState) {
	checkers := []func() Directions{}
}
