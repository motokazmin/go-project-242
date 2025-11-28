package code

import (
	"fmt"
	"os"
	"path/filepath"
)

// getSize возвращает размер файла или директории
// Внутренняя функция для получения размера в байтах
func getSize(path string, recursive, all bool) (int64, error) {
	// Проверяем существование пути
	stat, err := os.Stat(path)
	if err != nil {
		return 0, fmt.Errorf("не удалось получить информацию о пути: %w", err)
	}

	var size int64

	// Если это файл
	if !stat.IsDir() {
		size = stat.Size()
	} else {
		size, err = getDirSize(path, all, recursive)
		if err != nil {
			return 0, err
		}
	}

	return size, nil
}

// GetPathSize возвращает размер файла или директории в виде строки
// Если путь — файл, возвращает его размер.
// Если директория — суммирует размеры файлов первого уровня.
// Если human == true, размер форматируется в человекочитаемый вид (например "2.0KB").
// Если human == false, размер возвращается в байтах (например "2048B").
func GetPathSize(path string, recursive, human, all bool) (string, error) {
	// Получаем размер
	size, err := getSize(path, recursive, all)
	if err != nil {
		return "", err
	}

	// Форматируем размер
	sizeStr := FormatSize(size, human)

	// Возвращаем размер как строку
	return sizeStr, nil
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

// getDirSize суммирует размеры файлов в директории.
// Если 'recursive' == false, она считает только файлы первого уровня в 'dirPath', игнорируя поддиректории.
// Если 'recursive' == true, она рекурсивно суммирует размеры всех файлов в 'dirPath' и всех подпапках.
// Если 'all' == false, функция пропускает скрытые файлы и директории на каждом уровне обхода.
func getDirSize(dirPath string, all bool, recursive bool) (int64, error) {
	var totalSize int64

	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return 0, fmt.Errorf("не удалось прочитать директорию %s: %w", dirPath, err)
	}

	for _, entry := range entries {
		// Пропускаем скрытые файлы и директории, если all == false
		if !all && IsHidden(entry.Name()) {
			continue
		}

		fullPath := filepath.Join(dirPath, entry.Name())

		if entry.IsDir() {
			if recursive {
				// Рекурсивно подсчитываем размер поддиректории
				size, err := getDirSize(fullPath, all, recursive)
				if err != nil {
					// Пропускаем недоступные директории
					continue
				}
				totalSize += size
			}

		} else {
			// Добавляем размер файла
			info, err := entry.Info()
			if err != nil {
				// Пропускаем файлы, информацию о которых не получили
				continue
			}
			totalSize += info.Size()
		}
	}

	return totalSize, nil
}
