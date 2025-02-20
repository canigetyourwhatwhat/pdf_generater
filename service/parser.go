package service

import (
	"encoding/xml"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"pdfGenerater/models"
)

func ParseInputFile(filePath_ string) (*models.Form, error) {
	var form *models.Form
	var err error
	switch filepath.Ext(filePath_) {
	case ".xml":
		form, err = parseXML(filePath_)
	case ".json":
		form, err = parseJSON(filePath_)
	default:
		return nil, errors.New("unknown file type")
	}

	return form, err
}

// parseXML parses XML into a Form struct
func parseXML(filePath string) (*models.Form, error) {
	file, err := os.ReadFile(filePath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			slog.Info(fmt.Sprintf("file was not found in %s", filePath))
			os.Exit(1)
		} else {
			return nil, err
		}
	}

	var form models.Form
	err = xml.Unmarshal(file, &form)
	if err != nil {
		return nil, err
	}

	return &form, nil
}

// parseJSON parses JSON into a Form struct
func parseJSON(_ string) (*models.Form, error) {
	slog.Warn("JSON parsing is not implemented yet")
	return nil, nil
}
