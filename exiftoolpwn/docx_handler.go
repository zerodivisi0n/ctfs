package main

import (
	"archive/zip"
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
)

func docxHandler(w http.ResponseWriter, r *http.Request) error {
	docxFile, err := os.CreateTemp("", "exiftoolpwn*.docx")
	if err != nil {
		return fmt.Errorf("create docx file: %w", err)
	}
	defer os.Remove(docxFile.Name()) // clean up

	if _, err = io.Copy(docxFile, r.Body); err != nil {
		return fmt.Errorf("read input data: %w", err)
	}
	if err := docxFile.Close(); err != nil {
		return fmt.Errorf("close input file: %w", err)
	}

	title, err := readDocxTitle(docxFile.Name())
	if err != nil {
		return fmt.Errorf("read docx title: %w", err)
	}
	log.Printf("Docx title: %q", title)

	pdfFile, err := os.CreateTemp("", "exiftoolpwn*.pdf")
	if err != nil {
		return fmt.Errorf("create pdf file: %w", err)
	}
	defer os.Remove(pdfFile.Name()) // clean up
	pdfFile.Close()

	if err := docx2pdf(docxFile.Name(), pdfFile.Name()); err != nil {
		return fmt.Errorf("convert docx to pdf: %w", err)
	}

	if err := setExifTitle(pdfFile.Name(), title); err != nil {
		return fmt.Errorf("set title: %w", err)
	}

	pdfFile, err = os.Open(pdfFile.Name()) // reopen file fo read
	if err != nil {
		return fmt.Errorf("reopen pdf file: %w", err)
	}

	if _, err := io.Copy(w, pdfFile); err != nil {
		return fmt.Errorf("write output file: %w", err)
	}

	return nil
}

func docx2pdf(src, dest string) error {
	// unoconv -f pdf -o sample.pdf sample.docx
	cmd := exec.Command("unoconv", "-f", "pdf", "-o", dest, src)
	var errb bytes.Buffer
	cmd.Stderr = &errb
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("unoconv: %w (stderr: %s)", err, errb.String())
	}
	return nil
}

func readDocxTitle(filename string) (string, error) {
	zipReader, err := zip.OpenReader(filename)
	if err != nil {
		return "", fmt.Errorf("open docx file: %w", err)
	}
	defer zipReader.Close()

	var coreFile *zip.File
	for _, f := range zipReader.File {
		if f.Name == "docProps/core.xml" {
			coreFile = f
			break
		}
	}

	if coreFile == nil {
		return "", fmt.Errorf("docProps/core.xml file not found")
	}

	type sourceCoreProps struct {
		XMLName xml.Name `xml:"http://schemas.openxmlformats.org/package/2006/metadata/core-properties coreProperties"`
		Title   string   `xml:"http://purl.org/dc/elements/1.1/ title,omitempty"`
	}

	rc, err := coreFile.Open()
	if err != nil {
		return "", fmt.Errorf("open docProps/core.xml file: %w", err)
	}
	defer rc.Close()

	var props sourceCoreProps
	if err := xml.NewDecoder(rc).Decode(&props); err != nil {
		return "", fmt.Errorf("parse props xml: %w", err)
	}

	return props.Title, nil
}
