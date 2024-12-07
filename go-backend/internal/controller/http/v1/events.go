package v1

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	cuserr "github.com/kurochkinivan/pulskrsk/internal/customErrors"
)

type eventHandler struct {
	maxMemory int64
	uploadDir string
}

func NewEventHandler(maxMemory int64) *eventHandler {
	return &eventHandler{
		maxMemory: maxMemory,
		uploadDir: "../../uploads/images/",
	}
}

func (h *eventHandler) Register(mux *http.ServeMux) {
	mux.HandleFunc(fmt.Sprintf("%s %s", http.MethodPost, baseEventPath), errMdw(logMdw(h.saveFile)))
}

func (h *eventHandler) saveFile(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	const op string = "eventHandler.saveFile"

	err := r.ParseMultipartForm(h.maxMemory)
	if err != nil {
		return cuserr.ErrReadRequestBody.WithErr(fmt.Errorf("%s: %w", op, err))
	}

	file, handler, err := r.FormFile("file")
	if err != nil {
		return cuserr.ErrRetrieveDataFromFile.WithErr(fmt.Errorf("%s: %w", op, err))
	}
	defer file.Close()

	filePath := filepath.Join(h.uploadDir, handler.Filename)
	out, err := os.Create(filePath)
	if err != nil {
		return cuserr.ErrCreateFile.WithErr(fmt.Errorf("%s: %w", op, err))
	}
	defer out.Close()

	_, err = io.Copy(out, file)
	if err != nil {
		return cuserr.SystemError(err, op, "failed to copy file")
	}

	return nil
}
