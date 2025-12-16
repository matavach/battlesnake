package main

import (
	"log"
	"slices"
)

type MoveChecker interface {
	process() Directions
	Weight() float32
	SetWeight() float32
}

type MoveResult struct {
	weight float32
	result Directions
}
type MoveCheck func(s *snake, gi *GameInstance) (Directions, float32)

func (m *MoveResult) Weight() float32 {
	return m.weight
}

func (m *MoveResult) SetWeight(wt float32) {
	m.weight = wt
}

func MoveFloodFill(s *snake, gi *GameInstance) (Directions, float32) {
	floodResults := make(map[string]int)
	d := NewDirections()
	var weight float32 = 0.5
	for move := range s.SafeMoves {
		visited := make(map[Coord]bool)
		nextCoord := MoveCoordinate[move](s.Head)
		flood := FloodFill(nextCoord, gi, 10, visited)
		floodResults[move] = flood
		d[move] = (float32(flood) * float32(flood))
	}
	d.Normalize()
	name, val := d.Max()
	if floodResults[name] < 5 {
		weight = 0.1
	} else if floodResults[name] < 15 {
		weight = 0.5
	} else if floodResults[name] < 40 {
		weight = 1.0
	} else {
		weight = 1.5
	}
	log.Printf("Flood fill recommends [%s] at strength [%f]\n", name, val)

	return d, weight
}

func MoveToClosestFood(s *snake, gi *GameInstance) (Directions, float32) {
	d := NewDirections()
	var weight float32 = 0.8
	modified := map[string]bool{"up": false, "left": false, "right": false, "down": false}
	type foodPath struct {
		distance int
		path     []Coord
		foodPos  Coord
	}
	foodPaths := make([]foodPath, 0)

	for _, foodPos := range gi.Food {
		paths := GetAStar(*s.Head, foodPos, gi)
		for _, p := range paths {
			distance := len(p)
			if distance > 0 {
				foodPaths = append(foodPaths, foodPath{distance, p, foodPos})
			}
		}
	}

	// early return if empty
	if len(foodPaths) == 0 {
		log.Printf("MoveToClosestFood: No paths to food\n")
		return d, weight
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
				adj := float32(1.0 / (float32(fp.distance) * float32(fp.distance)))
				d[move] += adj
				modified[move] = true
			}
		}
	}
	d.Modified(modified)
	d.Normalize()

	if len(foodPaths) > 0 && IsClosestToFood(*s.Head, foodPaths[0].foodPos, gi) {
		weight = 1.5
	}

	name, val := d.Max()
	log.Printf("MoveToClosestFood: [%s] strength [%f]\n", name, val)
	return d, weight
}
func MoveBackwards(s *snake, gi *GameInstance) (Directions, float32) {
	d := NewDirections()
	var weight float32 = 1.0

	for dir := range d {
		nextCoord := MoveCoordinate[dir](s.Head)
		if nextCoord == *s.Neck {
			d[dir] = 0
		}
	}
	logUnsafe("move into neck", d)
	return d, weight
}

func MoveWithinBoard(s *snake, gi *GameInstance) (Directions, float32) {
	d := NewDirections()
	var weight float32 = 1.0
	for dir := range d {
		nextCoord := MoveCoordinate[dir](s.Head)
		if CheckOutOfBounds(nextCoord, gi) {
			d[dir] = 0
		}
	}
	logUnsafe("board", d)
	return d, weight
}

func MoveOntoSnake(s *snake, gi *GameInstance) (Directions, float32) {
	d := NewDirections()
	var weight float32 = 1.0

	for dir := range d {
		nextCoord := MoveCoordinate[dir](s.Head)
		if CheckSnake(nextCoord, gi) {
			d[dir] = 0
		}
	}
	logUnsafe("snake: ", d)
	return d, weight
}

func MoveHeadToHeadCollision(s *snake, gi *GameInstance) (Directions, float32) {
	d := NewDirections()
	var weight float32 = 1.0
	modified := map[string]bool{"up": false, "left": false, "right": false, "down": false}

	for move := range s.SafeMoves {
		nextCoord := MoveCoordinate[move](s.Head)
		collisionWeight := float32(1.0)

		for _, otherSnake := range gi.Snakes {
			if otherSnake.ID == gi.You.ID {
				continue
			}

			for _, otherDir := range []string{"up", "down", "left", "right"} {
				otherNextCoord := MoveCoordinate[otherDir](&otherSnake.Head)

				if nextCoord == otherNextCoord {
					if s.Length > otherSnake.Length {
						collisionWeight = 2.0
						weight = 0.8
					} else if s.Length < otherSnake.Length {
						collisionWeight = 0.1
						weight = 1.0
					} else {
						collisionWeight = 0.3
						weight = 0.9
					}
				}
			}
		}

		d[move] = collisionWeight
		if collisionWeight != 1.0 {
			modified[move] = true
		}
	}

	d.Modified(modified)
	d.Normalize()

	name, val := d.Max()
	log.Printf("MoveHeadToHeadCollision: [%s] strength [%f]\n", name, val)
	return d, weight
}
