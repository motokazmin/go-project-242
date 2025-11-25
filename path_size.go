package code

import (
	"fmt"
	"os"
)

// GetPathSize возвращает размер файла или директории в формате "<размер>\t<путь>"
// Если путь — файл, возвращает его размер.
// Если директория — суммирует размеры файлов первого уровня.
// Если human == true, размер форматируется в человекочитаемый вид.
func GetPathSize(path string, recursive, human, all bool) (string, error) {
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
		// Если это директория
		if recursive {
			// Рекурсивно суммируем все файлы и поддиректории
			size, err = getDirSizeRecursive(path, all)
		} else {
			// Суммируем размеры файлов только первого уровня
			size, err = getDirSize(path, all)
		}
		if err != nil {
			return "", err
		}
	}

	// Форматируем размер
	sizeStr := FormatSize(size, human)

	// Возвращаем результат в формате: <размер>\t<путь>
	return fmt.Sprintf("%s\t%s", sizeStr, path), nil
}

// IsHidden проверяет, является ли файл или директория скрытыми
// Файл считается скрытым, если его имя начинается с точки (.)
func IsHidden(filename string) bool {
	if len(filename) == 0 {
		return false
	}
	return filename[0] == '.'
}

// FormatSize форматирует размер байт в удобный вид
// Если human == false, возвращает строку вида "123B"
// Если human == true, конвертирует в человекочитаемый формат
// (единицы: B, KB, MB, GB, TB, PB, EB)
func FormatSize(bytes int64, human bool) string {
	if !human {
		return fmt.Sprintf("%dB", bytes)
	}

	// Единицы и их размеры
	units := []string{"B", "KB", "MB", "GB", "TB", "PB", "EB"}
	size := float64(bytes)

	for _, unit := range units {
		if size < 1024 {
			// Форматируем с 1 десятичным знаком для размеров >= 1KB
			if unit == "B" {
				return fmt.Sprintf("%.0f%s", size, unit)
			}
			return fmt.Sprintf("%.1f%s", size, unit)
		}
		size /= 1024
	}

	// На случай если размер больше EB
	return fmt.Sprintf("%.1f%s", size*1024, units[len(units)-1])
}

// getDirSize суммирует размеры файлов в директории (только первый уровень)
// Если all == false, пропускает скрытые файлы и директории
func getDirSize(dirPath string, all bool) (int64, error) {
	var totalSize int64

	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return 0, fmt.Errorf("не удалось прочитать директорию: %w", err)
	}

	for _, entry := range entries {
		// Пропускаем скрытые файлы и директории, если all == false
		if !all && IsHidden(entry.Name()) {
			continue
		}

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

// getDirSizeRecursive рекурсивно суммирует размеры всех файлов в директории и подпапках
// Если all == false, пропускает скрытые файлы и директории
func getDirSizeRecursive(dirPath string, all bool) (int64, error) {
	var totalSize int64

	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return 0, fmt.Errorf("не удалось прочитать директорию: %w", err)
	}

	for _, entry := range entries {
		// Пропускаем скрытые файлы и директории, если all == false
		if !all && IsHidden(entry.Name()) {
			continue
		}

		fullPath := dirPath + "/" + entry.Name()

		if entry.IsDir() {
			// Рекурсивно подсчитываем размер поддиректории
			size, err := getDirSizeRecursive(fullPath, all)
			if err != nil {
				continue // Пропускаем недоступные директории
			}
			totalSize += size
		} else {
			// Добавляем размер файла
			info, err := entry.Info()
			if err != nil {
				continue // Пропускаем файлы, информацию о которых не получили
			}
			totalSize += info.Size()
		}
	}

	return totalSize, nil
}
