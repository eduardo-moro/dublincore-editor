package docx

import (
	"archive/zip"
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/eduardo-moro/metadata-editor/dublincore"
)

const (
	corePropertiesPath = "docProps/core.xml"
)

// DOCX represents a DOCX document with Dublin Core metadata
type DOCX struct {
	FilePath   string
	DublinCore *dublincore.DublinCore
	FileData   []byte // Store the file content in memory
}

// ... (previous imports and constants)

// CoreProperties represents the full core.xml structure with CP namespace
type CoreProperties struct {
	XMLName      xml.Name `xml:"cp:coreProperties"`
	XMLNSCP      string   `xml:"xmlns:cp,attr"`
	XMLNSDC      string   `xml:"xmlns:dc,attr"`
	XMLNSDCTERMS string   `xml:"xmlns:dcterms,attr"`
	XMLNSXSI     string   `xml:"xmlns:xsi,attr"`

	// Dublin Core fields
	Title       []string `xml:"dc:title,omitempty"`
	Creator     []string `xml:"dc:creator,omitempty"`
	Subject     []string `xml:"dc:subject,omitempty"`
	Description []string `xml:"dc:description,omitempty"`

	// CP namespace fields
	Keywords []string `xml:"cp:keywords,omitempty"`
	Category []string `xml:"cp:category,omitempty"`
}

// ToXML converts CoreProperties to XML
func (cp *CoreProperties) ToXML() ([]byte, error) {
	cp.XMLNSCP = "http://schemas.openxmlformats.org/package/2006/metadata/core-properties"
	cp.XMLNSDC = "http://purl.org/dc/elements/1.1/"
	cp.XMLNSDCTERMS = "http://purl.org/dc/terms/"
	cp.XMLNSXSI = "http://www.w3.org/2001/XMLSchema-instance"

	header := `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>` + "\n"
	data, err := xml.MarshalIndent(cp, "", "  ")
	if err != nil {
		return nil, err
	}

	return []byte(header + string(data)), nil
}

// writeCoreProperties writes properly formatted core.xml with both DC and CP fields
func (d *DOCX) writeCoreProperties(zipWriter *zip.Writer) error {
	coreWriter, err := zipWriter.Create(corePropertiesPath)
	if err != nil {
		return fmt.Errorf("failed to create core.xml: %w", err)
	}

	// Create CoreProperties struct with both DC and CP fields
	coreProps := &CoreProperties{
		Title:       d.DublinCore.Title,
		Creator:     d.DublinCore.Creator,
		Subject:     d.DublinCore.Subject,
		Description: d.DublinCore.Description,
		Keywords:    d.DublinCore.Keywords,
		Category:    d.DublinCore.Category,
	}

	data, err := coreProps.ToXML()
	if err != nil {
		return fmt.Errorf("failed to marshal core properties: %w", err)
	}

	if _, err := coreWriter.Write(data); err != nil {
		return fmt.Errorf("failed to write core properties: %w", err)
	}

	return nil
}

// parseCoreXML parses standard DOCX core.xml with proper namespace handling
func parseCoreXML(data []byte) (*dublincore.DublinCore, error) {
	// First, try to parse with proper namespace handling
	var coreProps struct {
		XMLName     xml.Name `xml:"coreProperties"`
		Title       []string `xml:"title"`
		Creator     []string `xml:"creator"`
		Subject     []string `xml:"subject"`
		Description []string `xml:"description"`
		Keywords    []string `xml:"keywords"`
		Category    []string `xml:"category"`
	}

	if err := xml.Unmarshal(data, &coreProps); err != nil {
		return nil, fmt.Errorf("XML parsing failed: %w", err)
	}

	dc := dublincore.New()

	// Map the core properties to Dublin Core
	if len(coreProps.Title) > 0 {
		dc.Title = coreProps.Title
	}
	if len(coreProps.Creator) > 0 {
		dc.Creator = coreProps.Creator
	}
	if len(coreProps.Subject) > 0 {
		dc.Subject = coreProps.Subject
	}
	if len(coreProps.Description) > 0 {
		dc.Description = coreProps.Description
	}
	if len(coreProps.Keywords) > 0 {
		dc.Keywords = coreProps.Keywords
	}
	if len(coreProps.Category) > 0 {
		dc.Category = coreProps.Category
	}

	// If we found any data, return it
	if len(dc.Title) > 0 || len(dc.Creator) > 0 || len(dc.Keywords) > 0 || len(dc.Description) > 0 {
		return dc, nil
	}

	// If no data found with direct parsing, try alternative approach
	// Sometimes the XML uses different namespace prefixes
	return parseCoreXMLAlternative(data)
}

// parseCoreXMLAlternative tries alternative parsing approaches
func parseCoreXMLAlternative(data []byte) (*dublincore.DublinCore, error) {
	dc := dublincore.New()

	// Convert to string for manual inspection
	xmlStr := string(data)

	// Try to extract fields using string manipulation as fallback
	extractField := func(tag string) []string {
		var values []string
		startTag := "<" + tag + ">"
		endTag := "</" + tag + ">"

		start := strings.Index(xmlStr, startTag)
		for start != -1 {
			end := strings.Index(xmlStr[start:], endTag)
			if end == -1 {
				break
			}

			value := xmlStr[start+len(startTag) : start+end]
			values = append(values, strings.TrimSpace(value))

			// Find next occurrence
			start = strings.Index(xmlStr[start+end+len(endTag):], startTag)
			if start != -1 {
				start = start + end + len(endTag) // Adjust index
			}
		}
		return values
	}

	// Try different possible tag formats
	possibleTags := []string{
		"dc:title", "title", "cp:title",
		"dc:creator", "creator", "cp:creator",
		"dc:subject", "subject", "cp:subject",
		"dc:description", "description", "cp:description",
		"cp:keywords", "keywords",
		"cp:category", "category",
	}

	for _, tag := range possibleTags {
		values := extractField(tag)
		if len(values) > 0 {
			switch tag {
			case "dc:title", "title", "cp:title":
				dc.Title = values
			case "dc:creator", "creator", "cp:creator":
				dc.Creator = values
			case "dc:subject", "subject", "cp:subject":
				dc.Subject = values
			case "dc:description", "description", "cp:description":
				dc.Description = values
			case "cp:keywords", "keywords":
				dc.Keywords = values
			case "cp:category", "category":
				dc.Category = values
			}
		}
	}

	return dc, nil
}

// extractDublinCore extracts Dublin Core metadata from core.xml
func extractDublinCore(data []byte) (*dublincore.DublinCore, error) {
	// First try to parse as full core properties
	dc, err := parseCoreXML(data)
	if err == nil && (len(dc.Title) > 0 || len(dc.Creator) > 0 || len(dc.Keywords) > 0) {
		return dc, nil
	}

	// If that fails, try to parse as raw Dublin Core
	var rawDC dublincore.DublinCore
	if err := xml.Unmarshal(data, &rawDC); err != nil {
		return nil, err
	}

	return &rawDC, nil
}

// Open opens a DOCX file and reads its metadata
func Open(filePath string) (*DOCX, error) {
	// Read the entire file into memory
	fileData, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	// Create a zip reader from the file data
	reader, err := zip.NewReader(bytes.NewReader(fileData), int64(len(fileData)))
	if err != nil {
		return nil, fmt.Errorf("failed to create zip reader: %w", err)
	}

	docx := &DOCX{
		FilePath:   filePath,
		DublinCore: dublincore.New(),
		FileData:   fileData,
	}

	// Try to read existing Dublin Core metadata
	if coreFile, err := findFile(reader, corePropertiesPath); err == nil {
		coreData, err := readZipFile(coreFile)
		if err == nil {
			if dc, err := extractDublinCore(coreData); err == nil {
				docx.DublinCore = dc
			}
		}
	}

	return docx, nil
}

// Save saves the DOCX file with updated metadata
func (d *DOCX) Save(outputPath string) error {
	if outputPath == "" {
		outputPath = d.FilePath
	}

	// Create a zip reader from the original file data
	reader, err := zip.NewReader(bytes.NewReader(d.FileData), int64(len(d.FileData)))
	if err != nil {
		return fmt.Errorf("failed to create zip reader from memory: %w", err)
	}

	// Create output file
	outFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer outFile.Close()

	zipWriter := zip.NewWriter(outFile)
	defer zipWriter.Close()

	// Copy all files, replacing core.xml with updated metadata
	for _, file := range reader.File {
		if file.Name == corePropertiesPath {
			// Create new core.xml with updated metadata
			if err := d.writeCoreProperties(zipWriter); err != nil {
				return fmt.Errorf("failed to write core properties: %w", err)
			}
			continue
		}

		if err := copyZipFile(zipWriter, file); err != nil {
			return fmt.Errorf("failed to copy file %s: %w", file.Name, err)
		}
	}

	return nil
}

// Helper functions
func findFile(reader *zip.Reader, name string) (*zip.File, error) {
	for _, file := range reader.File {
		if file.Name == name {
			return file, nil
		}
	}
	return nil, fmt.Errorf("file not found: %s", name)
}

func readZipFile(file *zip.File) ([]byte, error) {
	rc, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer rc.Close()

	return io.ReadAll(rc)
}

func copyZipFile(dest *zip.Writer, src *zip.File) error {
	srcReader, err := src.Open()
	if err != nil {
		return err
	}
	defer srcReader.Close()

	destWriter, err := dest.Create(src.Name)
	if err != nil {
		return err
	}

	_, err = io.Copy(destWriter, srcReader)
	return err
}
