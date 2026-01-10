package statetransition

import (
	"github.com/devlongs/gean/common/types"
)

// StateTransition handles the state transition logic for the Lean Ethereum client
// This package contains:
// - State transition function
// - State changes tracking
// - Block processing logic

// Error represents errors that can occur during state transition
type Error struct {
	Message string
}

func (e Error) Error() string {
	return e.Message
}

// StateChanges represents the changes made to the state during a transition
type StateChanges struct {
	// TODO: Add fields for tracking state changes
}

// StateTransition applies a block to the current state and returns the resulting changes
func StateTransition(state *types.State, block *types.Block) (*StateChanges, error) {
	// TODO: Implement state transition logic
	return &StateChanges{}, nil
}
