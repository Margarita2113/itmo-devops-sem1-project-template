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
	tx, err := h.db.Begin()
	if err != nil {
		fmt.Println("error begin transaction")
		return
	}
	defer func() {
		err := tx.Rollback()
		if err != nil {
			fmt.Println("error rollback transaction")
		}
	}()

	for _, record := range records {
		product, err := model.NewProduct(record)
		if err != nil {
			fmt.Println("error creating product: %w", err)
			return
		}
		if err = h.db.Create(tx, product); err != nil {
			fmt.Println("error creating product db: %w", err)
			return
		}
		resp.TotalItems++

	}
	category, totalPrice, err := h.db.GetTotalPriceAndUnicCategory(tx)
	if err != nil {
		fmt.Println("error getting total price: %w", err)
		return
	}
	resp.TotalCategories = category
	resp.TotalPrice = totalPrice
	err = tx.Commit()
	if err != nil {
		fmt.Println("error commit transaction", err)
		return
	}

	marshal, err := json.Marshal(resp)
	if err != nil {
		fmt.Println("error marshalling response: %w", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(marshal)

}
