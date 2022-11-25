package handler

import (
	"io"
	"net/http"
	"os"
)

const PathFront = "../storage/front/"

func WriteFile(w http.ResponseWriter, filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = io.Copy(w, file)
	return err
}

func WriteError(w http.ResponseWriter) {
	io.WriteString(w, "Error 500")
}
