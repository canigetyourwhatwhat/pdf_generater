package service

import (
	"fmt"
	"log/slog"
	"os"
	"pdfGenerater/models"
	"regexp"
	"strconv"
	"strings"
)

// GenerateHTML is the main function to generate HTML from the given form
// It keeps concatenating one string to make the HTML
func GenerateHTML(form models.Form) string {
	html := "<!DOCTYPE html>\n<form>"
	// level is used to track the indentation to increase the readability of the generated HTML
	level := 0

	// validation for the form
	if len(form.Sections) == 0 || len(form.Fields) == 0 {
		slog.Error("Form should have at least one field or one section")
		os.Exit(1)
	}

	for _, section := range form.Sections {
		html += convertSectionToHTML(section, level, true)
	}

	html += "<br> <br> \n"
	for _, field := range form.Fields {
		html += convertFieldToHTML(field, level, true)
	}
	html += "\n</form>"
	return html
}

// convertSectionToHTML converts a Section element to HTML
// isOptional should be carried since the optional field is inherited from the parent Section element
func convertSectionToHTML(section models.Section, level int, isOptional bool) string {
	var indent string
	for i := 0; i < level+1; i++ {
		indent += "\n\t"
	}
	if section.Title == "" {
		slog.Error("Section should have a title", slog.String("invalid section Name", section.Name))
		os.Exit(1)
	}

	if len(section.Contents.Fields) == 0 && len(section.Contents.SubSections) == 0 {
		slog.Error("Section should have at least one field or one Section", slog.String("invalid section Name", section.Name))
		os.Exit(1)
	}

	html := fmt.Sprintf(indent+"<fieldset>"+indent+"\t<legend>%s</legend>", section.Title)

	if isOptional && section.Optional != nil {
		isOptional = *section.Optional
	}

	// add fields inside the section
	for _, field := range section.Contents.Fields {
		html += "<br> <br> \n" + convertFieldToHTML(field, level+1, isOptional)
	}

	// recursive function if there is another section inside the section.
	for _, subsection := range section.Contents.SubSections {
		html += "<br> <br>" + convertSectionToHTML(subsection, level+1, isOptional)
	}

	html += indent + "</fieldset>"
	return html
}

// convertFieldToHTML converts a Field element to HTML
func convertFieldToHTML(field models.Field, level int, isOptional bool) string {
	var indent string
	for i := 0; i < level+1; i++ {
		indent += "\t"
	}
	htmlStatement := indent + fmt.Sprintf("<label>%s</label> <br> \n"+indent, field.Caption)

	// If parent or this field is not optional, then it is required
	isRequired := "required"
	if !field.Optional || !isOptional {
		isRequired = ""
	}

	// Based on the Field type, generate the corresponding HTML element
	// If new element such as radio button is added in XML form, it should be added to this switch statement
	switch field.FieldType {
	case "TextBox":
		minLength, maxLength, lines, err := validateTextBoxType(field.Type)
		if err != nil {
			slog.Error("Format for TextBox is wrong", slog.String("Name", field.Name), slog.String("Invalid Type", field.Type))
			os.Exit(1)
		}
		htmlStatement += fmt.Sprintf("<textarea name='%s' %s minlength='%d' maxlength='%d' rows='%d'> </textarea>", field.Name, isRequired, minLength, maxLength, lines)
		return htmlStatement
	case "File":
		htmlStatement += fmt.Sprintf("<input type=file name='%s' %s/>", field.Name, isRequired)
		return htmlStatement
	case "Select":
		if field.Labels == nil || len(field.Labels) == 0 {
			slog.Error("Field should have at least one label for Select type", slog.String("invalid Field Name", field.Name))
			os.Exit(1)
		}
		availableOptions := extractEnumerationValues(field.Type)
		htmlStatement += fmt.Sprintf(" <select name='%s' %s>", field.Name, isRequired)
		// Validates whether the value in the Name attribute was included in the Type attribute
		for _, label := range field.Labels {
			found := false
			for _, option := range availableOptions {
				if label.Name == option {
					found = true
				}
			}
			if !found {
				slog.Error("Label Name is not in the available options", slog.String("invalid Name", label.Name), slog.String("Name of the corresponding Field", field.Name))
				os.Exit(1)
			}
			htmlStatement += fmt.Sprintf("\n"+indent+indent+"<option value='%s'>%s</option>", label.Name, label.Value)
		}
		htmlStatement += "\n" + indent + "</select>"
		return htmlStatement
	default:
		slog.Error("Unknown Filed type", slog.String("Field", field.FieldType), slog.String("Name", field.Name))
		os.Exit(1)
	}

	return htmlStatement + "<br>"
}

//--------------------- Helper functions ---------------------

// validateTextBoxType validates the TextBox type with schema validation
func validateTextBoxType(text string) (int, int, int, error) {
	re := regexp.MustCompile(`Text\(\[(\d+),(\d+)],Lines:(\d+)\)`)
	matches := re.FindStringSubmatch(text)
	if len(matches) != 4 {
		return 0, 0, 0, fmt.Errorf("invalid format")
	}

	minLength, err := strconv.Atoi(matches[1])
	if err != nil {
		return 0, 0, 0, err
	}

	maxLength, err := strconv.Atoi(matches[2])
	if err != nil {
		return 0, 0, 0, err
	}

	lines, err := strconv.Atoi(matches[3])
	if err != nil {
		return 0, 0, 0, err
	}

	return minLength, maxLength, lines, nil
}

// extractEnumerationValues extracts the values from the Enumeration type with schema validation
func extractEnumerationValues(enumString string) []string {
	re := regexp.MustCompile(`Enumeration\((.*?)\)`)
	matches := re.FindStringSubmatch(enumString)
	if len(matches) < 2 {
		slog.Error("HTML Select Filed type is wrong", slog.String("Type", enumString))
		os.Exit(1)
	} else if len(strings.Split(matches[1], ",")) < 2 {
		slog.Error("HTML Select Filed type should have more than 2 values to select", slog.String("Type", enumString))
		os.Exit(1)
	}
	return strings.Split(matches[1], ",")
}
