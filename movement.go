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

func MoveFloodFillWide(s *snake, gi *GameInstance) (Directions, float32) {
	floodResults := make(map[string]int)
	d := NewDirections()
	smallestFlood := 0
	var weight float32 = 0.6
	for move := range s.SafeMoves {
		visited := make(map[Coord]bool)
		nextCoord := MoveCoordinate[move](s.Head)
		flood := FloodFill(nextCoord, gi, 20, visited)
		floodResults[move] = flood
		if smallestFlood == 0 || flood < smallestFlood {
			smallestFlood = flood
		}
		d[move] = (float32(flood) * float32(flood))
	}
	d.Normalize()
	name, _ := d.Max()
	if smallestFlood < 5 {
		weight = 10
	} else if smallestFlood < 15 {
		weight = 5
	} else if smallestFlood < 40 {
		weight = 0.5
	} else {
		weight = 0.2
	}
	if len(gi.Snakes) == 1 {
		weight *= 3
	}
	log.Printf("Flood fill recommends [%s], weight[%f]\n", name, weight)

	return d, weight
}

func MoveFloodFillNarrow(s *snake, gi *GameInstance) (Directions, float32) {
	floodResults := make(map[string]int)
	d := NewDirections()
	smallestFlood := 0
	var weight float32 = 0.6
	for move := range s.SafeMoves {
		visited := make(map[Coord]bool)
		nextCoord := MoveCoordinate[move](s.Head)
		flood := FloodFill(nextCoord, gi, 5, visited)
		floodResults[move] = flood
		if smallestFlood == 0 || flood < smallestFlood {
			smallestFlood = flood
		}
		d[move] = (float32(flood) * float32(flood))
	}
	d.Normalize()
	name, _ := d.Max()
	if smallestFlood < 5 {
		weight = 10
	} else if smallestFlood < 10 {
		weight = 3
	} else if smallestFlood < 15 {
		weight = 1
	} else {
		weight = 0.5
	}
	if len(gi.Snakes) == 1 {
		weight *= 1.5
	}
	log.Printf("Flood fill recommends [%s], weight[%f]\n", name, weight)

	return d, weight
}

func MoveToClosestFood(s *snake, gi *GameInstance) (Directions, float32) {
	d := NewDirections()
	var weight float32 = 1.0
	modified := map[string]bool{"up": false, "left": false, "right": false, "down": false}
	type foodPath struct {
		distance int
		path     []Coord
		foodPos  Coord
	}
	foodPaths := make([]foodPath, 0)

foodLoop:
	for _, foodPos := range gi.Food {
		for _, otherSnakes := range gi.Snakes {
			if ManhattanDistance(otherSnakes.Head, foodPos) < 3 && ManhattanDistance(*s.Head, foodPos) > 3 {
				continue foodLoop
			}
		}
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

	if s.Length > LongestEnemy(gi).Length {
		weight *= 0.5
	}

	if len(foodPaths) > 0 && IsClosestToFood(*s.Head, foodPaths[0].foodPos, gi) {
		weight *= 3
	}

	if s.Data.Health > 60 {
		weight *= 0.8
	} else if s.Data.Health > 30 {
		weight *= 1
	} else {
		weight *= 1.5
	}

	name, _ := d.Max()
	log.Printf("MoveToClosestFood: [%s] weight [%f]\n", name, weight)
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

func MoveHeadtoHeadSafety(s *snake, gi *GameInstance) (Directions, float32) {
	d := NewDirections()
	var weight float32 = 1.0
dirLoop:
	for dir := range d {
		nextCoord := MoveCoordinate[dir](s.Head)
		for _, otherSnake := range gi.Snakes {
			for _, snakeDirection := range GetNeighbors(otherSnake.Head, gi) {
				if nextCoord == snakeDirection && s.Length <= otherSnake.Length {
					d[dir] = 0
					continue dirLoop
				}
			}
		}
	}
	logUnsafe("collision: ", d)
	return d, weight
}

func MoveHeadToHeadCollision(s *snake, gi *GameInstance) (Directions, float32) {
	d := NewDirections()
	var weight float32 = 0.0
	modified := map[string]bool{"up": false, "left": false, "right": false, "down": false}

	for move := range s.SafeMoves {
		nextCoord := MoveCoordinate[move](s.Head)
		var nextCoordWeight float32 = 0.2
		for _, otherSnake := range gi.Snakes {
			if nextCoordWeight == 0 {
				break
			}
			for _, snakeDirection := range GetNeighbors(otherSnake.Head, gi) {
				if nextCoord == snakeDirection {
					if s.Length <= otherSnake.Length {
						nextCoordWeight = 0
						weight = 1.0
						break
					} else {
						nextCoordWeight = 5

					}
				}
			}
		}
		d[move] = nextCoordWeight
		if nextCoordWeight != float32(0.2) {
			modified[move] = true
		}
	}

	d.Normalize()

	name, _ := d.Max()
	log.Printf("MoveHeadToHeadCollision: [%s], weight [%f]\n", name, weight)
	return d, weight
}

func MoveTowardEdges(s *snake, gi *GameInstance) (Directions, float32) {
	d := NewDirections()
	var weight float32 = 5
	modified := map[string]bool{"up": false, "left": false, "right": false, "down": false}

	edgeThreshold := 0 // How close to edge counts as "at edge"

	for move := range s.SafeMoves {
		nextCoord := MoveCoordinate[move](s.Head)

		// Check if we're at or moving to an edge
		isAtEdge := (nextCoord.X <= edgeThreshold || nextCoord.X >= gi.Bounds.X-1-edgeThreshold ||
			nextCoord.Y <= edgeThreshold || nextCoord.Y >= gi.Bounds.Y-1-edgeThreshold)

		if isAtEdge {
			d[move] = 1.5 // Strong preference for edge positions
			modified[move] = true
		} else {
			// Distance-based scoring for non-edge positions
			// Closer to edge = higher score
			distToNearestEdge := nextCoord.X
			if gi.Bounds.X-1-nextCoord.X < distToNearestEdge {
				distToNearestEdge = gi.Bounds.X - 1 - nextCoord.X
			}
			if nextCoord.Y < distToNearestEdge {
				distToNearestEdge = nextCoord.Y
			}
			if gi.Bounds.Y-1-nextCoord.Y < distToNearestEdge {
				distToNearestEdge = gi.Bounds.Y - 1 - nextCoord.Y
			}

			// Inverse scoring: closer to edge = higher score
			edgeScore := 1.0 / (1.0 + float32(distToNearestEdge)*0.5)
			d[move] = edgeScore
		}
	}

	d.Modified(modified)
	d.Normalize()

	name, val := d.Max()
	log.Printf("MoveTowardEdges: [%s] strength [%f]\n", name, val)
	return d, weight
}

func MoveCondense(s *snake, gi *GameInstance) (Directions, float32) {
	d := NewDirections()
	var weight float32 = 0.6
	modified := map[string]bool{"up": false, "left": false, "right": false, "down": false}

	// Get the direction from head to neck (direction we came from)
	neckDir := "up"
	if s.Head.X > s.Neck.X {
		neckDir = "left"
	} else if s.Head.X < s.Neck.X {
		neckDir = "right"
	} else if s.Head.Y > s.Neck.Y {
		neckDir = "down"
	}

	// Get perpendicular directions (parallel to body)
	perpendicularDirs := []string{}
	if neckDir == "up" || neckDir == "down" {
		perpendicularDirs = []string{"left", "right"}
	} else {
		perpendicularDirs = []string{"up", "down"}
	}

	// Score moves that move parallel to body (condensing)
	for move := range s.SafeMoves {
		condenseScore := float32(0.8) // Default lower score

		for _, perpDir := range perpendicularDirs {
			if move == perpDir {
				condenseScore = 1.2 // Prefer parallel moves
				modified[move] = true
				break
			}
		}

		d[move] = condenseScore
	}

	d.Modified(modified)
	d.Normalize()

	name, val := d.Max()
	log.Printf("MoveCondense: [%s] strength [%f]\n", name, val)
	return d, weight
}
