package main

type Directions map[string]float32

func NewDirections() Directions {
	return Directions{
		"up":    1,
		"down":  1,
		"left":  1,
		"right": 1,
	}
}

func (d Directions) Max() (name string, value float32) {
	var maxVal float32
	var dirName string
	for dir, v := range d {
		if v > maxVal {
			maxVal = v
			dirName = dir
		}
	}
	return dirName, maxVal
}

func (d Directions) Min() (name string, value float32) {
	var minVal float32 = float32(^uint32(0) >> 1)
	var dirName string
	for dir, v := range d {
		if v < minVal {
			minVal = v
			dirName = dir
		}
	}
	return dirName, minVal
}
func (d Directions) Normalize() {
	var _, maxVal = d.Max()
	for dir, value := range d {
		d[dir] = value / maxVal
	}
}

func (d Directions) InverseNormalize() {
	var highest float32
	var mult float32
	for dir, value := range d {
		if value == 0 {
			d[dir] = 0
			continue
		}
		new := 1.0 / value
		if new > highest {
			highest = new
			mult = 1 / highest
		}
		d[dir] = mult * new
	}
}

// maps are cool
func (d Directions) RemoveUnsafe() {
	for dir, val := range d {
		if val == 0 {
			delete(d, dir)
		}
	}
}

func (d Directions) Modified(m map[string]bool) {
	for move := range d {
		if !m[move] {
			delete(d, move)
		}
	}
}

func (d Directions) ApplyWeight(w float32) {
	for dir, value := range d {
		d[dir] = float32(value / w)
	}
}
