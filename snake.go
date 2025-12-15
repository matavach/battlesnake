package main

import (
	"log"
	"time"
)

type snake struct {
	Data        *Battlesnake
	SafeMoves   Directions
	Head        *Coord
	Neck        *Coord
	NextMove    string
	Length      int
	SecondLevel []*WeightedMove
}

// var weights = map[MoveChecker]float32{
// 	MoveFloodFill: 0.75,
// }

func NewSnake(bs *Battlesnake) (snake, error) {
	s := snake{
		Data:      bs,
		SafeMoves: NewDirections(),
		Head:      &bs.Head,
		Neck:      &bs.Body[1],
		Length:    len(bs.Body),
		SecondLevel: []*WeightedMove{
			{MoveFloodFill, 0.2},
			{MoveToClosestFood, 1.0},
		},
	}
	return s, nil
}

func (s *snake) UpdateSnake(bs *Battlesnake) error {
	s.Data = bs
	s.SafeMoves = NewDirections()
	s.Head = &bs.Head
	s.Neck = &bs.Body[1]
	s.Length = len(bs.Body)
	return nil
}

func (s *snake) processFirstMoveResults(ch chan Directions) {
	results := make([]Directions, 0, len(FirstLevelMoves))
	for c := range ch {
		results = append(results, c)
	}
	d := NewDirections()

	for dir := range d {
		var sum float32
		for _, result := range results {
			if result[dir] == 0 {
				d[dir] = 0
				break
			}
			sum += result[dir]
		}
		if d[dir] != 0 {
			d[dir] = sum / float32(len(results))
		}
	}
	s.SafeMoves = d
}

func (s *snake) processSecondMoveResults(ch chan MoveResult) Directions {
	results := make([]MoveResult, 0, len(s.SecondLevel))
	for c := range ch {
		results = append(results, c)
	}
	d := NewDirections()

	for dir := range d {
		var sum float32
		for _, result := range results {
			if result.dir[dir] == 0 {
				d[dir] = 0
				break
			}
			sum += result.dir[dir] * result.move.weight

		}
		if d[dir] != 0 {
			d[dir] = sum / float32(len(results))
		}
	}
	return d
}

func (s *snake) findSecondLevelMoves(gs *GameState) {
	defer LogExecutionTime("Second Level Moves:", time.Now())

	ch := make(chan MoveResult, len(s.SecondLevel))

	go func() {
		for _, move := range s.SecondLevel {
			// now := time.Now()
			res := move.process(s, gs)
			// LogMovement(fmt.Sprintf("%d", i), time.Since(now))
			ch <- MoveResult{move, res}
		}
		close(ch)
	}()
	res := s.processSecondMoveResults(ch)
	dir, str := res.Max()
	s.NextMove = dir
	log.Printf("Next move: '%s' with strength: '%f'", s.NextMove, str)
}

func (s *snake) findFirstLevelMoves(gs *GameState) {
	defer LogExecutionTime("findFirstLevelMoves", time.Now())

	ch := make(chan Directions, len(FirstLevelMoves))

	go func() {
		for _, move := range FirstLevelMoves {
			// now := time.Now()
			res := move.process(s, gs)
			// LogMovement("some move", time.Since(now))
			ch <- res
		}
		close(ch)
	}()

	s.processFirstMoveResults(ch)
	s.SafeMoves.RemoveUnsafe()
}
