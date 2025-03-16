package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"project_sem/internal/model"
	"project_sem/internal/zip"
)

type response struct {
	TotalItems      int     `json:"total_items"`
	TotalCategories int     `json:"total_categories"`
	TotalPrice      float64 `json:"total_price"`
}

func (h *Handler) POSTHandler(w http.ResponseWriter, r *http.Request) {

	// Получаем размер ZIP-файла для создания нового zip.Reader
	zipFileSize := r.ContentLength

	// Считываем данные в буфер
	buf := new(bytes.Buffer)
	if _, err := io.Copy(buf, r.Body); err != nil {
		fmt.Println("error reading response body: %w", err)
		return
	}

	// Вызываем функцию разархивации
	records, err := zip.Unzip(bytes.NewReader(buf.Bytes()), "destination", zipFileSize)
	if err != nil {
		fmt.Printf("Ошибка разархивации: %v\n", err)
		return
	} else {
		fmt.Println("Разархивация завершена успешно!")
	}

	var resp response
	var unicCategory map[string]string
	for _, record := range records {
		product, err := model.NewProduct(record)
		if err != nil {
			fmt.Println("error creating product: %w", err)
			continue
		}
		if err = h.db.Create(product); err != nil {
			fmt.Println("error creating product db: %w", err)
			continue
		}

		resp.TotalItems++
		if _, ok := unicCategory[product.Category]; !ok {
			resp.TotalCategories++
		}
	}

	allProducts, err := h.db.Get()
	if err != nil {
		fmt.Println("error getting products: %w", err)
		return
	}
	for _, product := range allProducts {
		resp.TotalPrice = resp.TotalPrice + product.Price
	}

	marshal, err := json.Marshal(resp)
	if err != nil {
		fmt.Println("error marshalling response: %w", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(marshal)

}
