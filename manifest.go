package gen

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

const manifestFileName = ".genmanifest.json"

type genManifest struct {
	Version int `json:"version"`
	Mode    uint `json:"mode"`
	Tables  map[string]genManifestTable `json:"tables,omitempty"`
	Files   map[string]string          `json:"files,omitempty"`
}

type genManifestTable struct {
	ModelStructName string `json:"model_struct_name"`
	QueryStructName string `json:"query_struct_name"`
	FileName        string `json:"file_name"`
}

func loadManifest(dir string) (*genManifest, string, error) {
	path := filepath.Join(dir, manifestFileName)
	b, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return &genManifest{Version: 1, Tables: map[string]genManifestTable{}, Files: map[string]string{}}, path, nil
		}
		return nil, "", err
	}
	var m genManifest
	if err := json.Unmarshal(b, &m); err != nil {
		return nil, "", err
	}
	if m.Version == 0 {
		m.Version = 1
	}
	if m.Tables == nil {
		m.Tables = map[string]genManifestTable{}
	}
	if m.Files == nil {
		m.Files = map[string]string{}
	}
	return &m, path, nil
}

func saveManifest(path string, m *genManifest) error {
	b, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return err
	}
	b = append(b, '\n')
	return os.WriteFile(path, b, 0640)
}

func sha256Hex(b []byte) string {
	sum := sha256.Sum256(b)
	return hex.EncodeToString(sum[:])
}
