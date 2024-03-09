package jsonfile

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/oktavarium/gophkeeper/internal/client/crypto"
)

// JSONFile holds a Go value of type Data and persists it to a JSON file.
// Data is accessed and modified using the Read and Write methods.
// Create a JSONFile using the New or Load functions.
type JSONFile[Data any] struct {
	path   string
	crypto *crypto.Crypto
	mu     sync.RWMutex
	bytes  []byte
	data   *Data
}

// New creates a new empty JSONFile at the given path.
func New[Data any](path string, c *crypto.Crypto) (*JSONFile[Data], error) {
	p := &JSONFile[Data]{path: path, crypto: c, bytes: []byte("{}"), data: new(Data)}
	if err := p.Write(func(*Data) error { return nil }); err != nil {
		return nil, fmt.Errorf("jsonfile.New: %w", err)
	}
	return p, nil
}

// Load loads an existing JSONFileData from the given path.
//
// If the file does not exist, Load returns an error that can be
// checked with os.IsNotExist.
//
// Load and New are separate to avoid creating a new file when
// starting a service, which could lead to data loss. To both load an
// existing file or create it (which you may want to do in a development
// environment), combine Load with New, like this:
//
//	db, err := jsonfile.Load[Data](path)
//	if os.IsNotExist(err) {
//		db, err = jsonfile.New[Data](path)
//	}
func Load[Data any](path string, c *crypto.Crypto) (*JSONFile[Data], error) {
	var err error
	var encryptedBytes []byte

	p := &JSONFile[Data]{path: path, crypto: c, data: new(Data)}

	encryptedBytes, err = os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("jsonfile.Load: %w", err)
	}

	p.bytes, err = p.crypto.DecryptData(encryptedBytes)
	if err != nil {
		return nil, fmt.Errorf("error on data decrypting: %w", err)
	}

	err = json.Unmarshal(p.bytes, p.data)
	if err != nil {
		return nil, fmt.Errorf("jsonfile.Load: %w", err)
	}
	return p, nil
}

// Read calls fn with the current copy of the data.
func (p *JSONFile[Data]) Read(fn func(data *Data)) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	fn(p.data)
}

// Write calls fn with a copy of the data, then writes the changes to the file.
// If fn returns an error, Write does not change the file and returns the error.
func (p *JSONFile[Data]) Write(fn func(*Data) error) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	data := new(Data) // operate on copy to allow concurrent reads and rollback
	if err := json.Unmarshal(p.bytes, data); err != nil {
		return fmt.Errorf("JSONFile.Write: %w", err)
	}
	if err := fn(data); err != nil {
		return err
	}
	b, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("JSONFile.Write: %w", err)
	}
	if bytes.Equal(b, p.bytes) {
		return nil // no change
	}

	f, err := os.CreateTemp(filepath.Dir(p.path), filepath.Base(p.path)+".tmp")
	if err != nil {
		return fmt.Errorf("JSONFile.Write: temp: %w", err)
	}

	var encryptedBytes []byte
	encryptedBytes, err = p.crypto.EncryptData(b)
	if err != nil {
		return fmt.Errorf("error on encrypting data: %w", err)
	}

	_, err = f.Write(encryptedBytes)
	if err1 := f.Close(); err1 != nil && err == nil {
		err = err1
	}
	if err != nil {
		return fmt.Errorf("JSONFile.Write: %w", err)
	}

	err = os.Rename(f.Name(), p.path)
	if err != nil {
		return fmt.Errorf("JSONFile.Write: rename: %w", err)
	}

	data = new(Data) // avoid any aliased memory
	err = json.Unmarshal(b, data)
	if err != nil {
		return fmt.Errorf("JSONFile.Write: %w", err)
	}

	p.data = data
	p.bytes = b
	return nil
}
