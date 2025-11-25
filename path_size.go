package code

import (
	"fmt"
	"os"
)

// GetSize возвращает размер файла или директории в формате "<размер>\t<путь>"
// Если путь — файл, возвращает его размер.
// Если директория — суммирует размеры файлов первого уровня.
func GetSize(path string, recursive, human, all bool) (string, error) {
	// Проверяем существование пути
	stat, err := os.Stat(path)
	if err != nil {
		return "", fmt.Errorf("не удалось получить информацию о пути: %w", err)
	}

	var size int64

	// Если это файл
	if !stat.IsDir() {
		size = stat.Size()
	} else {
		// Если это директория - суммируем размеры файлов первого уровня
		size, err = getDirSize(path)
		if err != nil {
			return "", err
		}
	}

	// Возвращаем результат в формате: <размер>\t<путь>
	return fmt.Sprintf("%d\t%s", size, path), nil
}

// getDirSize суммирует размеры файлов в директории (только первый уровень)
func getDirSize(dirPath string) (int64, error) {
	var totalSize int64

	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return 0, fmt.Errorf("не удалось прочитать директорию: %w", err)
	}

	for _, entry := range entries {
		// Пропускаем директории, берём только файлы первого уровня
		if !entry.IsDir() {
			info, err := entry.Info()
			if err != nil {
				continue // Пропускаем файлы, информацию о которых не получили
			}
			totalSize += info.Size()
		}
	}

	return totalSize, nil
}

