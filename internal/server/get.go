package server

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"net/http"
	"project_sem/internal/zip"
	"time"
)

func (h *Handler) GetHandler(w http.ResponseWriter, r *http.Request) {

	get, err := h.db.Get()
	if err != nil {
		fmt.Println("error db %w", err)
		return
	}

	buf := new(bytes.Buffer)

	csvWriter := csv.NewWriter(buf)

	for _, product := range get {
		str := []string{fmt.Sprintf("%d", product.ID),
			product.Name,
			product.Category,
			fmt.Sprintf("%.f", product.Price),
			product.Data.Format(time.DateOnly)}
		err := csvWriter.Write(str)
		if err != nil {
			fmt.Println("error csv writer")
			continue
		}
	}

	csvWriter.Flush()

	reader := io.Reader(buf)

	zipFiles, err := zip.ZipFiles(reader)
	if err != nil {
		fmt.Println("error zip files %w", err)
		return
	}

	w.Write(zipFiles.Bytes())

}
