package main

import (
	"log"
	"time"
)

func logUnsafe(moveType string, d Directions) {
	unsafe := false
	for dir, strength := range d {
		if strength <= 0 {
			log.Printf("Found unsafe move [%s] for '%s' check\n", dir, moveType)
			unsafe = true
		}
	}
	if !unsafe {
		log.Printf("No unsafe moves for '%s' found\n", moveType)
	}
}

// LogExecutionTime logs how long a function takes
func LogExecutionTime(name string, start time.Time) {
	duration := time.Since(start)
	log.Printf("[%s] completed in %v\n", name, duration)
}

// LogMovement logs move evaluation with timing
func LogMovement(moveType string, duration time.Duration) {
	log.Printf("[MOVEMENT] %s evaluated in %v\n", moveType, duration)
}
