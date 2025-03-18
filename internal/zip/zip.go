package zip

import (
	"archive/zip"
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func Unzip(zipReader io.ReaderAt, destination string, zipFileSize int64) ([][]string, error) {
	// Открываем ZIP-архив
	zipFile, err := zip.NewReader(zipReader, zipFileSize)
	if err != nil {
		return nil, fmt.Errorf("не удалось открыть zip-файл: %v", err)
	}
	var allRecords [][]string
	// Перебираем файлы в архиве
	for _, file := range zipFile.File {
		// Получаем полное имя файла
		filePath := filepath.Join(destination, file.Name)

		// Проверяем, является ли файл директорией, и создаем её, если это так
		if file.FileInfo().IsDir() {
			os.MkdirAll(filePath, os.ModePerm)
			continue
		}
		if !strings.Contains(file.FileInfo().Name(), ".csv") {
			continue
		}

		reader, err := file.Open()
		if err != nil {
			return nil, fmt.Errorf("ошибка при открытии файла")
		}

		readAll, err := csv.NewReader(reader).ReadAll()
		if err != nil {
			return nil, fmt.Errorf("read data error %w", err)
		}
		allRecords = append(allRecords, readAll[1:]...) // без header
	}

	return allRecords, nil
}

// ZipFiles создает ZIP-архив из указанных файлов
func ZipFiles(r io.Reader) (*bytes.Buffer, error) {
	buf := new(bytes.Buffer)

	zipWriter := zip.NewWriter(buf)
	defer zipWriter.Close()

	// Создаем заголовок для добавляемого файла в архив
	writer, err := zipWriter.Create("data.csv")
	if err != nil {
		return nil, fmt.Errorf("ошибка создания записи для файла в ZIP: %v", err)
	}

	// Копируем содержимое файла в ZIP
	_, err = io.Copy(writer, r)
	if err != nil {
		return nil, fmt.Errorf("ошибка копирования файла %s в ZIP: %v", err)
	}

	return buf, nil
}
