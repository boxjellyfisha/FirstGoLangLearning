package pkg

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
)

func JsonResponse(w http.ResponseWriter, response any) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func GetResourceDir() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	var dbPath = filepath.Join(dir, "res")

	if !dirExists(dbPath) {
		os.Mkdir(dbPath, 0755)
	}

	return dbPath, nil
}

func dirExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.IsDir()
}
