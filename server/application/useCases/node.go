package usecases

import (
	"io"
	"os"

	"github.com/google/uuid"
)

type NodeIDUseCase struct {
}

func NewNodeIDUseCase() *NodeIDUseCase {
	return &NodeIDUseCase{}
}

func (r *NodeIDUseCase) ReadFromFile() ([]byte, error) {
	file, err := os.Open(".ID")
	if err != nil {
		return nil, err
	}

	defer func() {
		if cerr := file.Close(); err == nil {
			err = cerr
		}
	}()

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	return data, nil // We should check for errors here to make the function safer and more robust in real-world applications
}

func (r *NodeIDUseCase) ReadOrCreateFile() ([]byte, error) {
	data, err := r.ReadFromFile()
	if err != nil {
		// If file does not exist, create it and write the default ID to this new file
		err = os.WriteFile(".ID", []byte(uuid.New().String()), 0644)
		if err != nil {
			return nil, err
		}
		data, err = r.ReadFromFile() // Attempting read after writing the default ID to a new file
	}
	return data, err
}
