package code

import (
	"strings"
	"testing"
)

// TestGetSize_File проверяет возврат размера для одного файла
func TestGetSize_File(t *testing.T) {
	// Используем тестовый файл из testdata
	testFile := "testdata/65012_melt.log"

	result, err := GetSize(testFile, false, false, false)
	if err != nil {
		t.Fatalf("GetSize вернул ошибку: %v", err)
	}

	// Проверяем, что результат содержит размер и путь
	parts := strings.Split(result, "\t")
	if len(parts) != 2 {
		t.Fatalf("Результат должен содержать размер и путь, разделённые табуляцией: %s", result)
	}

	// Проверяем, что размер равен 1863B байтам
	expectedSize := "1863B"
	if parts[0] != expectedSize {
		t.Errorf("Ожидается размер %s, получено: %s", expectedSize, parts[0])
	}

	// Проверяем, что путь совпадает
	if parts[1] != testFile {
		t.Errorf("Ожидается путь %s, получено: %s", testFile, parts[1])
	}
}

// TestGetSize_AnotherFile проверяет размер другого файла
func TestGetSize_AnotherFile(t *testing.T) {
	testFile := "testdata/65049_melt.log"

	result, err := GetSize(testFile, false, false, false)
	if err != nil {
		t.Fatalf("GetSize вернул ошибку: %v", err)
	}

	parts := strings.Split(result, "\t")
	if len(parts) != 2 {
		t.Fatalf("Результат должен содержать размер и путь: %s", result)
	}

	// Проверяем, что размер равен 4129B байтам
	expectedSize := "4129B"
	if parts[0] != expectedSize {
		t.Errorf("Ожидается размер %s, получено: %s", expectedSize, parts[0])
	}
}

// TestGetSize_DirectoryFirstLevel проверяет суммирование файлов первого уровня в LogViewer
func TestGetSize_DirectoryFirstLevel(t *testing.T) {
	testDir := "testdata/LogViewer"

	result, err := GetSize(testDir, false, false, false)
	if err != nil {
		t.Fatalf("GetSize вернул ошибку: %v", err)
	}

	parts := strings.Split(result, "\t")
	if len(parts) != 2 {
		t.Fatalf("Результат должен содержать размер и путь: %s", result)
	}

	// В LogViewer первого уровня три файла: 348 + 65 + 155 = 568 байт
	// Папка logs игнорируется, так как это директория
	expectedSize := "568B"
	if parts[0] != expectedSize {
		t.Errorf("Ожидается размер %s (348 + 65 + 155), получено: %s", expectedSize, parts[0])
	}

	if parts[1] != testDir {
		t.Errorf("Ожидается путь %s, получено: %s", testDir, parts[1])
	}
}

// TestFormatSize_Bytes проверяет форматирование байтов без human флага
func TestFormatSize_Bytes(t *testing.T) {
	tests := []struct {
		bytes    int64
		human    bool
		expected string
	}{
		{0, false, "0B"},
		{1, false, "1B"},
		{123, false, "123B"},
		{1023, false, "1023B"},
	}

	for _, test := range tests {
		result := FormatSize(test.bytes, test.human)
		if result != test.expected {
			t.Errorf("FormatSize(%d, %v): ожидается %s, получено %s", test.bytes, test.human, test.expected, result)
		}
	}
}

// TestFormatSize_HumanReadable проверяет форматирование в человекочитаемый вид
func TestFormatSize_HumanReadable(t *testing.T) {
	tests := []struct {
		bytes    int64
		expected string
	}{
		{0, "0B"},
		{512, "512B"},
		{1024, "1.0KB"},
		{1024 * 1024, "1.0MB"},
		{1536 * 1024, "1.5MB"},               // 1.5 MB
		{1024 * 1024 * 1024, "1.0GB"},        // 1 GB
		{1024 * 1024 * 1024 * 2, "2.0GB"},    // 2 GB
		{1024 * 1024 * 1024 * 1024, "1.0TB"}, // 1 TB
	}

	for _, test := range tests {
		result := FormatSize(test.bytes, true)
		if result != test.expected {
			t.Errorf("FormatSize(%d, true): ожидается %s, получено %s", test.bytes, test.expected, result)
		}
	}
}

// TestGetSize_NonExistentPath проверяет ошибку для несуществующего пути
func TestGetSize_NonExistentPath(t *testing.T) {
	result, err := GetSize("testdata/nonexistent_file.txt", false, false, false)

	if err == nil {
		t.Errorf("Ожидается ошибка для несуществующего пути, но получено: %s", result)
	}

	if result != "" {
		t.Errorf("Ожидается пустой результат при ошибке, получено: %s", result)
	}
}

// TestIsHidden проверяет функцию определения скрытых файлов
func TestIsHidden(t *testing.T) {
	tests := []struct {
		filename string
		expected bool
	}{
		{".hidden", true},
		{".config", true},
		{".bashrc", true},
		{"file.txt", false},
		{"README.md", false},
		{"", false},
	}

	for _, test := range tests {
		result := IsHidden(test.filename)
		if result != test.expected {
			t.Errorf("IsHidden(%q): ожидается %v, получено %v", test.filename, test.expected, result)
		}
	}
}

// TestGetSize_WithHiddenFiles проверяет учёт скрытых файлов
func TestGetSize_WithHiddenFiles(t *testing.T) {
	// В testdata есть скрытые файлы (.hidden_file_1, .hidden_file_2)
	// С флагом all они должны быть учтены

	resultWith, err := GetSize("testdata", false, false, true)
	if err != nil {
		t.Fatalf("GetSize вернул ошибку: %v", err)
	}

	// Разделяем результаты
	parts := strings.Split(resultWith, "\t")

	if len(parts) != 2 {
		t.Fatalf("Неправильный формат результата")
	}

	// В testdata со скрытыми файлами: 5992 + 14 + 14 = 6020
	expectedWith := "6020B"

	if parts[0] != expectedWith {
		t.Errorf("С флагом all ожидается %s, получено: %s", expectedWith, parts[0])
	}
}

// TestGetSize_DirectoryWithoutHidden проверяет, что скрытые файлы игнорируются по умолчанию
func TestGetSize_DirectoryWithoutHidden(t *testing.T) {
	// Тестируем с all=false
	result, err := GetSize("testdata/LogViewer", false, false, false)
	if err != nil {
		t.Fatalf("GetSize вернул ошибку: %v", err)
	}

	parts := strings.Split(result, "\t")
	if len(parts) != 2 {
		t.Fatalf("Результат должен содержать размер и путь: %s", result)
	}

	// В LogViewer первого уровня три XML файла: 348 + 65 + 155 = 568 байт
	expectedSize := "568B"
	if parts[0] != expectedSize {
		t.Errorf("Ожидается размер %s, получено: %s", expectedSize, parts[0])
	}
}
