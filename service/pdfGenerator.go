package service

import (
	"os"
	"os/exec"
	"path/filepath"
)

func ConvertHTMLStringToPDF(htmlFilePath, outputFilePath string) error {
	// Get absolute path to the binary
	binaryPath, err := filepath.Abs("./bin/weasyprint_wrapper")
	if err != nil {
		return err
	}

	cmd := exec.Command(binaryPath, htmlFilePath, outputFilePath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
