package main

import "slices"

var MoveCoordinate = map[string]func(c *Coord) Coord{
	"up":    func(c *Coord) Coord { return Coord{c.X, c.Y + 1} },
	"down":  func(c *Coord) Coord { return Coord{c.X, c.Y - 1} },
	"left":  func(c *Coord) Coord { return Coord{c.X - 1, c.Y} },
	"right": func(c *Coord) Coord { return Coord{c.X + 1, c.Y} },
}

// return true if coord is in bounds and not snake
func CheckValidCoord(c Coord, gi *GameInstance) bool {
	return !CheckOutOfBounds(c, gi) && !CheckSnake(c, gi)
}

// return true if coord is out of bounds
func CheckOutOfBounds(c Coord, gi *GameInstance) bool {
	return (c.X >= gi.Bounds.X || c.X < 0 || c.Y >= gi.Bounds.Y || c.Y < 0)
}

// return true if coord is a snake
func CheckSnake(c Coord, gi *GameInstance) bool {
	for coord, state := range gi.Board {
		if c == coord && state == CoordSnake {
			return true
		}
	}
	return false
}

// return true if coord is food
func CheckFood(c Coord, gi *GameInstance) bool {
	for coord, state := range gi.Board {
		if c == coord && state == CoordFood {
			return true
		}
	}
	return false
}

func FloodFill(pos Coord, gi *GameInstance, depth int, visited map[Coord]bool) int {
	if depth == 0 {
		return 0
	}
	if visited[pos] {
		return 0
	}
	if !CheckValidCoord(pos, gi) {
		return 0
	}
	visited[pos] = true
	count := 1
	for _, neighbor := range GetNeighbors(pos, gi) {
		count += FloodFill(neighbor, gi, depth-1, visited)
	}
	return count
}
func GetNeighbors(c Coord, gi *GameInstance) []Coord {
	neighbors := make([]Coord, 0, 4)
	for _, moveFunc := range MoveCoordinate {
		tempCoord := moveFunc(&c)
		if CheckValidCoord(tempCoord, gi) {
			neighbors = append(neighbors, tempCoord)
		}
	}
	return neighbors
}

type NodeAStar struct {
	parent    *NodeAStar
	pos       Coord
	distance  int
	heuristic int
	cost      int
}

type path []Coord

func GetAStar(start Coord, end Coord, gi *GameInstance) []path {
	open := make(map[Coord]*NodeAStar)
	closed := make(map[Coord]bool)
	startNode := &NodeAStar{
		pos: start,
	}
	minCost := -1
	paths := make([]path, 0)

	open[start] = startNode
	for len(open) > 0 {
		var currentNode *NodeAStar
		var currentPos Coord

		for pos, node := range open {
			if currentNode == nil || node.cost < currentNode.cost {
				currentNode = node
				currentPos = pos
			}
		}

		if minCost >= 0 && currentNode.cost > minCost {
			break
		}

		delete(open, currentPos)

		if currentNode.pos == end {
			tempPath := make(path, 0, 10)
			tempNode := currentNode
			for tempNode != nil {
				tempPath = append(tempPath, tempNode.pos)
				tempNode = tempNode.parent
			}
			slices.Reverse(tempPath)
			tempPath = tempPath[1:] // path[0] is the starting position
			paths = append(paths, tempPath)
			minCost = currentNode.cost
			continue
		}

		for _, neighbor := range GetNeighbors(currentNode.pos, gi) {
			if closed[neighbor] {
				continue
			}
			childNode := &NodeAStar{
				parent:    currentNode,
				pos:       neighbor,
				distance:  currentNode.distance + 1,
				heuristic: (neighbor.X-end.X)*(neighbor.X-end.X) + (neighbor.Y-end.Y)*(neighbor.Y-end.Y),
			}
			childNode.cost = childNode.distance + childNode.heuristic
			if existingNode, exists := open[neighbor]; exists {
				if childNode.cost < existingNode.cost {
					open[neighbor] = childNode
				}
			} else {
				open[neighbor] = childNode
			}
		}

		closed[currentPos] = true

	}
	return paths
}

func IsClosestToFood(coord Coord, foodNode Coord, gi *GameInstance) bool {

	myDist := ManhattanDistance(coord, foodNode)

	for _, otherSnake := range gi.Snakes {
		if otherSnake.ID == gi.You.ID {
			continue
		}

		otherDist := ManhattanDistance(otherSnake.Head, foodNode)

		if otherDist < myDist {
			return false
		}
	}

	return true
}

// ManhattanDistance calculates the Manhattan distance between two coordinates
func ManhattanDistance(a, b Coord) int {
	dx := a.X - b.X
	if dx < 0 {
		dx = -dx
	}
	dy := a.Y - b.Y
	if dy < 0 {
		dy = -dy
	}
	return dx + dy
}

func LongestEnemy(gi *GameInstance) Battlesnake {
	var longest int = 0
	enemy := Battlesnake{}

	for _, s := range gi.Snakes {
		if s.Length > longest {
			longest = s.Length
			enemy = s
		}
	}
	return enemy
}
