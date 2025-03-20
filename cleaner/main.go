package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func main() {
	// Папки для обработки (относительные пути от папки проекта)
	folders := []string{
		"../pyakin/txt",
		"../zaznobin/txt",
	}

	// Обработка каждой папки
	for _, folder := range folders {
		processFolder(folder)
	}

	fmt.Println("Обработка завершена.")
}

func processFolder(folderPath string) {
	// Определяем папку для сохранения результатов
	outputFolder := filepath.Dir(folderPath)

	// Создаем файлы для сохранения очищенного корпуса
	outputFileWithNewlines, err := os.Create(filepath.Join(outputFolder, "cleaned_corpus.txt"))
	if err != nil {
		fmt.Println("Ошибка при создании файла:", err)
		return
	}
	defer outputFileWithNewlines.Close()

	outputFileSingleLine, err := os.Create(filepath.Join(outputFolder, "cleaned_corpus_single_line.txt"))
	if err != nil {
		fmt.Println("Ошибка при создании файла:", err)
		return
	}
	defer outputFileSingleLine.Close()

	// Регулярное выражение для поиска строк времени
	timeRegex := regexp.MustCompile(`^\d{1,2}:\d{2}(:\d{2})?$`)

	// Чтение всех файлов в папке
	files, err := os.ReadDir(folderPath)
	if err != nil {
		fmt.Println("Ошибка при чтении папки:", err)
		return
	}

	// Обработка каждого файла
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".txt") {
			filePath := filepath.Join(folderPath, file.Name())
			processFile(filePath, outputFileWithNewlines, outputFileSingleLine, timeRegex)
		}
	}

	fmt.Printf("Обработка папки %s завершена.\n", folderPath)
}

func processFile(filePath string, outputFileWithNewlines, outputFileSingleLine *os.File, timeRegex *regexp.Regexp) {
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Ошибка при открытии файла:", err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		// Если строка не соответствует формату времени, обрабатываем её
		if !timeRegex.MatchString(line) {
			// Запись в файл с переносами строк
			_, err := outputFileWithNewlines.WriteString(line + "\n")
			if err != nil {
				fmt.Println("Ошибка при записи в файл:", err)
				return
			}

			// Запись в файл с одной строкой
			_, err = outputFileSingleLine.WriteString(line + " ")
			if err != nil {
				fmt.Println("Ошибка при записи в файл:", err)
				return
			}
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Ошибка при чтении файла:", err)
	}
}
