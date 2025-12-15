package main

import "slices"

var MoveCoordinate = map[string]func(c *Coord) Coord{
	"up":    func(c *Coord) Coord { return Coord{c.X, c.Y + 1} },
	"down":  func(c *Coord) Coord { return Coord{c.X, c.Y - 1} },
	"left":  func(c *Coord) Coord { return Coord{c.X - 1, c.Y} },
	"right": func(c *Coord) Coord { return Coord{c.X + 1, c.Y} },
}

// return true if coord is in bounds and not snake
func CheckValidCoord(c Coord, gb GameBoard, gs *GameState) bool {
	return !CheckOutOfBounds(c, gb, gs) && !CheckSnake(c, gb, gs)
}

// return true if coord is out of bounds
func CheckOutOfBounds(c Coord, gb GameBoard, gs *GameState) bool {
	return (c.X >= gs.Board.Width || c.X < 0 || c.Y >= gs.Board.Height || c.Y < 0)
}

// return true if coord is a snake
func CheckSnake(c Coord, b GameBoard, gs *GameState) bool {
	for coord, state := range b {
		if c == coord && state == CoordSnake {
			return true
		}
	}
	return false
}

// return true if coord is food
func CheckFood(c Coord, b GameBoard, gs *GameState) bool {
	for coord, state := range b {
		if c == coord && state == CoordFood {
			return true
		}
	}
	return false
}

func FloodFill(pos Coord, gb GameBoard, gs *GameState, depth int, visited map[Coord]bool) int {
	if depth == 0 {
		return 0
	}
	if visited[pos] {
		return 0
	}
	if !CheckValidCoord(pos, gb, gs) {
		return 0
	}
	visited[pos] = true
	count := 1
	for _, neighbor := range GetNeighbors(pos, gb, gs) {
		count += FloodFill(neighbor, gb, gs, depth-1, visited)
	}
	return count
}
func GetNeighbors(c Coord, gb GameBoard, gs *GameState) []Coord {
	neighbors := make([]Coord, 0, 4)
	for _, moveFunc := range MoveCoordinate {
		tempCoord := moveFunc(&c)
		if CheckValidCoord(tempCoord, gb, gs) {
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

func GetAStar(start Coord, end Coord, gb GameBoard, gs *GameState) []path {
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

		for _, neighbor := range GetNeighbors(currentNode.pos, gb, gs) {
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
