package main

import (
	"fmt"
	"log/slog"
	"os"
	"pdfGenerater/service"
)

func main() {
	// configure slog for logging
	lvl := new(slog.LevelVar)
	lvl.Set(slog.LevelDebug)
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{
		Level: lvl,
	})))

	// -------------  Parse file into struct  -------------
	form, err := service.ParseInputFile("input.xml")
	if err != nil {
		fmt.Println("Failed to parse XML file: ", err)
		return
	}

	// -------------  Generate HTML  -------------
	htmlFilePath := "output.html"
	htmlOutputStr := service.GenerateHTML(*form)
	if err = os.WriteFile(htmlFilePath, []byte(htmlOutputStr), 0644); err != nil {
		fmt.Println("Error writing HTML file:", err)
	} else {
		fmt.Println("HTML file generated successfully at path:", htmlFilePath)
	}

	// -------------  Generate PDF from HTML  -------------
	outputFilePath := "output.pdf"
	if err = service.ConvertHTMLStringToPDF(htmlFilePath, outputFilePath); err != nil {
		fmt.Println("Failed converting PDF from HTML:", err)
	} else {
		fmt.Println("PDF generated successfully at path:", outputFilePath)
	}
}
