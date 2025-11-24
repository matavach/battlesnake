package main

type MoveChecker interface {
	process() Directions
}

type MoveCheck func() Directions

func (mc MoveCheck) process() Directions {
	return mc()
}

func (s *snake) OneDimensionalCheck(c *Coord, d *Directions) {
	if s.Head.X-1 == c.X {
		d.left = false
	}
	if s.Head.X+1 == c.X {
		d.right = false
	}
	if s.Head.Y-1 == c.Y {
		d.down = false
	}
	if s.Head.Y+1 == c.Y {
		d.up = false
	}
}

func (s *snake) MoveBackwards(gs *GameState) Directions {
	d := Directions{}
	// using switch because each case is mutually exclusive - ie, can't have neck be below and left of head
	s.OneDimensionalCheck(&s.Neck, &d)
	return d
}

func (s *snake) MoveWithinBoard(gs *GameState) Directions {
	d := Directions{}
	if s.Head.X == 0 {
		d.left = false
	}
	if s.Head.X == gs.Board.Width-1 {
		d.right = false
	}
	if s.Head.Y == 0 {
		d.down = false
	}
	if s.Head.Y == gs.Board.Height-1 {
		d.up = false
	}
	return d
}

func (s *snake) MoveOntoSnake(gs *GameState) Directions {
	d := Directions{}
	adjacent := map[Coord]bool{}
	for _, enemy := range gs.Board.Snakes {
		// skip myself
		if enemy.ID == s.Data.ID {
			continue
		}
		for _, cell := range enemy.Body {
			adjacent[cell] = true
		}
	}

	if adjacent[Coord{s.Head.X - 1, s.Head.Y}] {
		d.left = false
	}
	if adjacent[Coord{s.Head.X + 1, s.Head.Y}] {
		d.right = false
	}
	if adjacent[Coord{s.Head.X, s.Head.Y - 1}] {
		d.down = false
	}
	if adjacent[Coord{s.Head.X, s.Head.Y + 1}] {
		d.up = false
	}
	return d
}
