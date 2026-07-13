package sync

import (
	"encoding/json"
	"errors"
	"io/fs"
	"os"
	"path/filepath"
)

// LoadState loads local state from the specified path.
func LoadState(path string) (*LocalState, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return &LocalState{
				RegisteredRecords: make(map[string]string),
			}, nil
		}
		return nil, err
	}

	var state LocalState
	if err := json.Unmarshal(data, &state); err != nil {
		return nil, err
	}

	if state.RegisteredRecords == nil {
		state.RegisteredRecords = make(map[string]string)
	}

	return &state, nil
}

// SaveState saves local state to the specified path.
func SaveState(path string, state *LocalState) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}
