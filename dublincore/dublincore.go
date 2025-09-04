package dublincore

import (
	"encoding/xml"
	"time"
)

// DublinCore represents the Dublin Core metadata elements
type DublinCore struct {
	XMLName     xml.Name `xml:"http://purl.org/dc/elements/1.1/ dc"`
	Title       []string `xml:"title,omitempty"`
	Creator     []string `xml:"creator,omitempty"`
	Subject     []string `xml:"subject,omitempty"`
	Description []string `xml:"description,omitempty"`
	Publisher   []string `xml:"publisher,omitempty"`
	Contributor []string `xml:"contributor,omitempty"`
	Date        []string `xml:"date,omitempty"`
	Type        []string `xml:"type,omitempty"`
	Format      []string `xml:"format,omitempty"`
	Identifier  []string `xml:"identifier,omitempty"`
	Source      []string `xml:"source,omitempty"`
	Language    []string `xml:"language,omitempty"`
	Relation    []string `xml:"relation,omitempty"`
	Coverage    []string `xml:"coverage,omitempty"`
	Rights      []string `xml:"rights,omitempty"`

	// Custom fields for CP namespace
	Keywords []string `xml:"http://purl.org/dc/terms/ keyword,omitempty"`
	Category []string `xml:"http://purl.org/dc/terms/ type,omitempty"` // Using type for category
}

// New creates a new DublinCore instance with default values
func New() *DublinCore {
	return &DublinCore{
		Date:     []string{time.Now().Format(time.RFC3339)},
		Format:   []string{"application/vnd.openxmlformats-officedocument.wordprocessingml.document"},
		Category: []string{"curriculo"}, // Default category
	}
}

// SetTitle sets the title
func (dc *DublinCore) SetTitle(title string) {
	dc.Title = []string{title}
}

// AddCreator adds a creator
func (dc *DublinCore) AddCreator(creator string) {
	dc.Creator = append(dc.Creator, creator)
}

// SetDescription sets the description
func (dc *DublinCore) SetDescription(description string) {
	dc.Description = []string{description}
}

// AddKeyword adds a keyword
func (dc *DublinCore) AddKeyword(keyword string) {
	dc.Keywords = append(dc.Keywords, keyword)
}

// SetCategory sets the category (always to "curriculo")
func (dc *DublinCore) SetCategory() {
	dc.Category = []string{"curriculo"}
}

// ToXML converts Dublin Core metadata to XML
func (dc *DublinCore) ToXML() ([]byte, error) {
	return xml.MarshalIndent(dc, "", "  ")
}

// FromXML parses Dublin Core metadata from XML
func FromXML(data []byte) (*DublinCore, error) {
	var dc DublinCore
	err := xml.Unmarshal(data, &dc)
	if err != nil {
		return nil, err
	}
	return &dc, nil
}
