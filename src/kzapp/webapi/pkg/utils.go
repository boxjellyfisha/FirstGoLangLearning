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

func GetCurrentDir() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "..","webapi"), nil
}