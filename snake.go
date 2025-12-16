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
	FirstLevel  []MoveCheck
	SecondLevel []MoveCheck
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
		FirstLevel: []MoveCheck{
			MoveBackwards,
			MoveOntoSnake,
			MoveWithinBoard,
		},
		SecondLevel: []MoveCheck{
			MoveFloodFill,
			MoveToClosestFood,
			MoveHeadToHeadCollision,
		},
	}
	return s, nil
}

func (s *snake) findSecondLevelMoves(gi *GameInstance) {
	defer LogExecutionTime("Second Level Moves:", time.Now())

	dirCh := make(chan MoveResult, len(s.SecondLevel))

	go func() {
		for _, move := range s.SecondLevel {
			dir, weight := move(s, gi)
			dirCh <- MoveResult{weight, dir}
		}
		close(dirCh)
	}()

	results := make([]MoveResult, 0, len(s.SecondLevel))
	for result := range dirCh {
		results = append(results, result)
	}

	d := NewDirections()
	for dirName := range d {
		var sum float32
		var weightSum float32
		for _, result := range results {
			if result.result[dirName] == 0 {
				d[dirName] = 0
				break
			}
			sum += result.result[dirName] * result.weight
			weightSum += result.weight
		}
		if d[dirName] != 0 {
			d[dirName] = sum / weightSum
		}
	}

	dirName, str := d.Max()
	s.NextMove = dirName
	log.Printf("Next move: '%s' with strength: '%f'", s.NextMove, str)
}

func (s *snake) findFirstLevelMoves(gi *GameInstance) {
	defer LogExecutionTime("Second Level Moves:", time.Now())

	dirCh := make(chan MoveResult, len(s.FirstLevel))

	go func() {
		for _, move := range s.FirstLevel {
			dir, weight := move(s, gi)
			dirCh <- MoveResult{weight, dir}
		}
		close(dirCh)
	}()

	results := make([]MoveResult, 0, len(s.FirstLevel))
	for result := range dirCh {
		results = append(results, result)
	}

	d := NewDirections()
	for dirName := range d {
		for _, result := range results {
			if result.result[dirName] == 0 {
				delete(d, dirName)
				break
			}
		}
	}
	s.SafeMoves = d
}
