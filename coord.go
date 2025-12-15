package main

type CoordState int

func (cs CoordState) String() string {
	return stateName[cs]
}

const (
	CoordEmpty = iota
	CoordFood
	CoordSnake
)

var stateName = map[CoordState]string{
	CoordEmpty: "empty",
	CoordFood:  "food",
	CoordSnake: "snake",
}
