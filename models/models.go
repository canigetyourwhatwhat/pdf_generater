package models

import "encoding/xml"

type Form struct {
	XMLName  xml.Name  `xml:"Form"`
	Sections []Section `xml:"Section"`
	Fields   []Field   `xml:"Field"`
}

type Section struct {
	Name     string   `xml:"Name,attr"`
	Optional *bool    `xml:"Optional,attr"`
	Title    string   `xml:"Title"`
	Contents Contents `xml:"Contents"`
}

// Contents is a struct that contains fields and subsections
// If new element such as Comment Box is added, it should be added to this struct
type Contents struct {
	Fields      []Field   `xml:"Field"`
	SubSections []Section `xml:"Section"`
}

type Field struct {
	Name      string  `xml:"Name,attr"`
	Type      string  `xml:"Type,attr"`
	Optional  bool    `xml:"Optional,attr"`
	FieldType string  `xml:"FieldType,attr"`
	Caption   string  `xml:"Caption"`
	Labels    []Label `xml:"Labels>Label"`
}

type Label struct {
	Name  string `xml:"Name,attr"`
	Value string `xml:",chardata"`
}
