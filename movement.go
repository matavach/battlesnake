package main

import (
	"log"
	"slices"
)

var FirstLevelMoves = []MoveCheck{
	MoveBackwards,
	MoveWithinBoard,
	MoveOntoSnake,
}

type MoveChecker interface {
	process() Directions
	Weight() float32
	SetWeight() float32
}

type MoveResult struct {
	move *WeightedMove
	dir  Directions
}

type WeightedMove struct {
	move   MoveCheck
	weight float32
}

func (w *WeightedMove) process(s *snake, gs *GameState) Directions {
	return w.move(s, gs)
}

func (w *WeightedMove) Weight() float32 {
	return w.weight
}

func (w *WeightedMove) SetWeight(wt float32) {
	w.weight = wt
}

type MoveCheck func(s *snake, gs *GameState) Directions

func (mc MoveCheck) process(s *snake, gs *GameState) Directions {
	return mc(s, gs)
}

func MoveFloodFill(s *snake, gs *GameState) Directions {
	d := NewDirections()
	gb := NewGameBoard(gs)

	for move := range s.SafeMoves {
		visited := make(map[Coord]bool)
		nextCoord := MoveCoordinate[move](s.Head)
		flood := FloodFill(nextCoord, gb, gs, 5, visited)
		d[move] = float32(flood)
	}
	d.Normalize()
	name, val := d.Max()
	log.Printf("Flood fill recommends [%s] at strength [%f]\n", name, val)
	return d
}

func MoveToClosestFood(s *snake, gs *GameState) Directions {
	d := NewDirections()
	gb := NewGameBoard(gs)
	modified := map[string]bool{"up": false, "left": false, "right": false, "down": false}
	type foodPath struct {
		distance int
		path     []Coord
	}
	foodPaths := make([]foodPath, 0)
	// first, find all the possible paths to food
	for _, foodPos := range gs.Board.Food {
		paths := GetAStar(*s.Head, foodPos, gb, gs)
		for _, p := range paths {
			distance := len(p)
			if distance > 0 {
				foodPaths = append(foodPaths, foodPath{distance, p})
			}
		}

	}
	// early return if empty
	if len(foodPaths) == 0 {
		return d
	}

	// I wanted to write  a SortFunc. This seems more obtuse than a loop in this case.
	slices.SortFunc(foodPaths, func(a, b foodPath) int {
		if a.distance < b.distance {
			return -1
		}
		if a.distance > b.distance {
			return 1
		}
		return 0
	})

	for move := range s.SafeMoves {
		nextCoord := MoveCoordinate[move](s.Head)
		for _, fp := range foodPaths {
			if fp.path[0] == nextCoord {
				adj := float32(1 / float32(fp.distance))
				d[move] /= adj
				modified[move] = true
			}
		}
	}
	d.Modified(modified)
	d.InverseNormalize()

	name, val := d.Max()
	log.Printf("Closest food recommends [%s] at strength [%f]\n", name, val)
	return d
}

func MoveBackwards(s *snake, gs *GameState) Directions {
	d := NewDirections()
	if s.Length < 2 {
		return d
	}

	for dir := range d {
		nextCoord := MoveCoordinate[dir](s.Head)
		if nextCoord == *s.Neck {
			d[dir] = 0
		}
	}
	logUnsafe("move into neck", d)
	return d
}

func MoveWithinBoard(s *snake, gs *GameState) Directions {
	d := NewDirections()
	for dir := range d {
		nextCoord := MoveCoordinate[dir](s.Head)
		if CheckOutOfBounds(nextCoord, nil, gs) {
			d[dir] = 0
		}
	}
	logUnsafe("board", d)
	return d
}

func MoveOntoSnake(s *snake, gs *GameState) Directions {
	d := NewDirections()
	gb := NewGameBoard(gs)

	for dir := range d {
		nextCoord := MoveCoordinate[dir](s.Head)
		if CheckSnake(nextCoord, gb, gs) {
			d[dir] = 0
		}
	}
	logUnsafe("snake: ", d)
	return d
}
