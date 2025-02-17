package main

import (
	"fmt"
	"html"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/barasher/go-exiftool"
)

func pdfHandler(w http.ResponseWriter, r *http.Request) error {
	pdfFile, err := os.CreateTemp("", "exiftoolpwn*.pdf")
	if err != nil {
		return fmt.Errorf("create pdf file: %w", err)
	}
	defer os.Remove(pdfFile.Name()) // clean up

	if _, err = io.Copy(pdfFile, r.Body); err != nil {
		return fmt.Errorf("read input data: %w", err)
	}
	if err := pdfFile.Close(); err != nil {
		return fmt.Errorf("close input file: %w", err)
	}

	title := r.URL.Query().Get("title")
	if title == "" {
		return fmt.Errorf("empty 'title' param")
	}
	title = html.EscapeString(title)
	log.Printf("Query title: %q", title)

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

func setExifTitle(filename string, title string) error {
	et, err := exiftool.NewExiftool()
	if err != nil {
		return fmt.Errorf("init exiftool: %w", err)
	}
	defer et.Close()

	metadata := []exiftool.FileMetadata{{
		File: filename,
		Fields: map[string]any{
			"title": title,
		},
	}}
	et.WriteMetadata(metadata)
	if err := metadata[0].Err; err != nil {
		return fmt.Errorf("write metadata: %w", err)
	}

	return nil
}
