package service

import (
	"os"
	"os/exec"
	"path/filepath"
)

func ConvertHTMLStringToPDF(htmlFilePath, pdfFilePath string) error {
	// try to use the binary from the bin folder
	binaryPath, err := filepath.Abs("./bin/weasyprint_wrapper")
	if err != nil {
		return err
	}

	cmd := exec.Command(binaryPath, htmlFilePath, pdfFilePath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err = cmd.Run(); err != nil {
		// if it doesn't work, try to use the locally installed weasyprint library
		cmd = exec.Command("weasyprint", htmlFilePath, pdfFilePath)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		// Run the command
		return cmd.Run()
	}

	return nil
}
