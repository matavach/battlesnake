package main

import (
	"sync"
)

type GameManager struct {
	GameState GameState
	mu        *sync.RWMutex
}

var GameInstance = GameManager{}

func NewGameManager(gs GameState) *GameManager {
	return &GameManager{
		GameState: gs,
		mu:        &sync.RWMutex{},
	}
}

func (gm *GameManager) GetBoard() *Board {
	gm.mu.RLock()
	defer gm.mu.RUnlock()
	return &gm.GameState.Board
}

func (gm *GameManager) SetGameState(state GameState) {
	gm.mu.Lock()
	defer gm.mu.Unlock()
	gm.GameState = state
}

func (gm *GameManager) ReadGameState() GameState {
	gm.mu.RLock()
	defer gm.mu.Unlock()
	return gm.GameState
}
