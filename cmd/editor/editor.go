package editor

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/eduardo-moro/metadata-editor/docx"
	"github.com/eduardo-moro/metadata-editor/dublincore"
	"github.com/eduardo-moro/metadata-editor/ui"
	"github.com/urfave/cli/v2"
)

func Main() {
	app := &cli.App{
		Name:  "dublin-core-editor",
		Usage: "Edit Dublin Core metadata in DOCX files with a nice TUI",
		Commands: []*cli.Command{
			{
				Name:    "edit",
				Aliases: []string{"e"},
				Usage:   "Edit metadata with TUI interface",
				Action: func(c *cli.Context) error {
					if c.NArg() == 0 {
						return fmt.Errorf("please provide a DOCX file path")
					}
					filePath := c.Args().First()
					outputPath := c.String("output")
					return editWithTUI(filePath, outputPath)
				},
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "output",
						Aliases: []string{"o"},
						Usage:   "Output file (default: overwrite original)",
					},
				},
			},
			{
				Name:    "debug",
				Aliases: []string{"d"},
				Usage:   "Debug the internal structure of a DOCX file",
				Action:  debugDOCX,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "file",
						Aliases:  []string{"f"},
						Usage:    "DOCX file to debug",
						Required: true,
					},
				},
			},
			{
				Name:    "view",
				Aliases: []string{"v"},
				Usage:   "View current metadata",
				Action:  viewMetadata,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "file",
						Aliases:  []string{"f"},
						Usage:    "DOCX file to view",
						Required: true,
					},
				},
			},
		},
		// Default action if no command is specified
		Action: func(c *cli.Context) error {
			if c.NArg() == 0 {
				return fmt.Errorf("please provide a DOCX file path and command\nUse --help for usage information")
			}
			// Default to edit command if file is provided without command
			filePath := c.Args().First()
			return editWithTUI(filePath, "")
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

// Add the viewMetadata function
func viewMetadata(c *cli.Context) error {
	filePath := c.String("file")

	if err := validateFileExists(filePath); err != nil {
		return err
	}

	doc, err := docx.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open DOCX file: %w", err)
	}

	fmt.Printf("📂 File: %s\n", filePath)
	fmt.Println("Current metadata:")
	printCurrentMetadata(doc.DublinCore)

	return nil
}

func editWithTUI(filePath, outputPath string) error {
	// Open the DOCX file
	doc, err := docx.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open DOCX file: %w", err)
	}

	fmt.Printf("📂 Opening: %s\n", filePath)
	fmt.Println("Current metadata:")
	printCurrentMetadata(doc.DublinCore)
	fmt.Println("\nLoading TUI editor...")
	fmt.Println("Note: Type your metadata and press Enter to submit.")
	fmt.Println()

	// Store original metadata for comparison
	originalDC := &dublincore.DublinCore{}
	originalDC.Title = append([]string{}, doc.DublinCore.Title...)
	originalDC.Creator = append([]string{}, doc.DublinCore.Creator...)
	originalDC.Keywords = append([]string{}, doc.DublinCore.Keywords...)
	originalDC.Description = append([]string{}, doc.DublinCore.Description...)
	originalDC.Category = append([]string{}, doc.DublinCore.Category...)

	// Run the BubbleTea TUI
	updatedDC, cancelled, err := ui.RunEditor(doc.DublinCore)
	if err != nil {
		return fmt.Errorf("TUI editor failed: %w", err)
	}

	if cancelled {
		fmt.Println("❌ Edit cancelled. No changes made.")
		return nil
	}

	// Simple change detection - compare string representations
	changesMade := hasChanges(originalDC, updatedDC)
	if !changesMade {
		fmt.Println("✅ No changes made. File remains unchanged.")
		return nil
	}

	// Update the document with new metadata
	doc.DublinCore = updatedDC

	// Handle output path
	if outputPath == "" {
		backupPath := filePath + ".backup"
		if err := createBackup(filePath, backupPath); err != nil {
			return fmt.Errorf("backup failed: %w", err)
		}
		fmt.Printf("✅ Created backup: %s\n", backupPath)
		outputPath = filePath
	}

	// Save changes
	if err := doc.Save(outputPath); err != nil {
		return fmt.Errorf("failed to save DOCX file: %w", err)
	}

	fmt.Printf("\n✅ Metadata updated successfully in %s\n", outputPath)
	fmt.Println("\nUpdated metadata:")
	printMetadata(doc.DublinCore)

	return nil
}

// Simple change detection function
func hasChanges(original, updated *dublincore.DublinCore) bool {
	if strings.Join(original.Title, "|") != strings.Join(updated.Title, "|") {
		return true
	}
	if strings.Join(original.Creator, "|") != strings.Join(updated.Creator, "|") {
		return true
	}
	if strings.Join(original.Keywords, "|") != strings.Join(updated.Keywords, "|") {
		return true
	}
	if strings.Join(original.Description, "|") != strings.Join(updated.Description, "|") {
		return true
	}
	return false
}

func printCurrentMetadata(dc *dublincore.DublinCore) {
	fmt.Printf("📝 Title:       %s\n", getValueOrNone(dc.Title))
	fmt.Printf("👤 Creator(s):  %s\n", getValueOrNone(dc.Creator))
	fmt.Printf("🔑 Keywords:    %s\n", getValueOrNone(dc.Keywords))
	fmt.Printf("📋 Description: %s\n", getValueOrNone(dc.Description))
	fmt.Printf("📂 Category:    %s\n", getValueOrNone(dc.Category))
}

func printMetadata(dc *dublincore.DublinCore) {
	fmt.Printf("📝 Title:       %s\n", strings.Join(dc.Title, ", "))
	fmt.Printf("👤 Creator(s):  %s\n", strings.Join(dc.Creator, ", "))
	fmt.Printf("🔑 Keywords:    %s\n", strings.Join(dc.Keywords, ", "))
	fmt.Printf("📋 Description: %s\n", strings.Join(dc.Description, ", "))
	fmt.Printf("📂 Category:    %s\n", strings.Join(dc.Category, ", "))
}

func getValueOrNone(values []string) string {
	if len(values) == 0 || (len(values) == 1 && values[0] == "") {
		return "(none)"
	}
	return strings.Join(values, ", ")
}

func createBackup(src, dst string) error {
	input, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	return os.WriteFile(dst, input, 0644)
}

func hasRealChanges(original, updated *dublincore.DublinCore) bool {
	// Simple comparison - if any field is different, changes were made
	if strings.Join(original.Title, ",") != strings.Join(updated.Title, ",") {
		return true
	}
	if strings.Join(original.Creator, ",") != strings.Join(updated.Creator, ",") {
		return true
	}
	if strings.Join(original.Keywords, ",") != strings.Join(updated.Keywords, ",") {
		return true
	}
	if strings.Join(original.Description, ",") != strings.Join(updated.Description, ",") {
		return true
	}
	return false
}

func debugDOCX(c *cli.Context) error {
	filePath := c.String("file")

	if err := validateFileExists(filePath); err != nil {
		return err
	}

	fmt.Printf("🔍 Debugging: %s\n", filePath)

	// Read the file directly to see what's in core.xml
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return fmt.Errorf("failed to get file info: %w", err)
	}

	reader, err := zip.NewReader(file, fileInfo.Size())
	if err != nil {
		return fmt.Errorf("failed to create zip reader: %w", err)
	}

	// Look for core.xml
	coreFile, err := findZipFile(reader, "docProps/core.xml")
	if err != nil {
		return fmt.Errorf("core.xml not found: %w", err)
	}

	coreData, err := readZipFile(coreFile)
	if err != nil {
		return fmt.Errorf("failed to read core.xml: %w", err)
	}

	fmt.Println("=== Raw core.xml content ===")
	fmt.Println(string(coreData))
	fmt.Println("===========================")

	// Try to parse it
	doc, err := docx.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open with docx parser: %w", err)
	}

	fmt.Println("=== Parsed metadata ===")
	printMetadata(doc.DublinCore)

	return nil
}

func findZipFile(reader *zip.Reader, name string) (*zip.File, error) {
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

func validateFileExists(filePath string) error {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return fmt.Errorf("file does not exist: %s", filePath)
	}
	return nil
}
