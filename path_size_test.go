package code

import (
	"testing"
)

// TestGetSize_File проверяет возврат размера для одного файла
func TestGetSize_File(t *testing.T) {
	testFile := "testdata/65012_melt.log"

	result, err := GetPathSize(testFile, false, false, false)
	if err != nil {
		t.Fatalf("GetPathSize вернул ошибку: %v", err)
	}

	expectedSize := "1863B"
	if result != expectedSize {
		t.Errorf("Ожидается размер %s, получено: %s", expectedSize, result)
	}
}

// TestGetSize_AnotherFile проверяет размер другого файла
func TestGetSize_AnotherFile(t *testing.T) {
	testFile := "testdata/65049_melt.log"

	result, err := GetPathSize(testFile, false, false, false)
	if err != nil {
		t.Fatalf("GetPathSize вернул ошибку: %v", err)
	}

	expectedSize := "4129B"
	if result != expectedSize {
		t.Errorf("Ожидается размер %s, получено: %s", expectedSize, result)
	}
}

// TestGetSize_DirectoryFirstLevel проверяет суммирование файлов первого уровня в LogViewer
func TestGetSize_DirectoryFirstLevel(t *testing.T) {
	testDir := "testdata/LogViewer"

	result, err := GetPathSize(testDir, false, false, false)
	if err != nil {
		t.Fatalf("GetPathSize вернул ошибку: %v", err)
	}

	// В LogViewer первого уровня три XML файла: 348 + 65 + 155 = 568 байт
	expectedSize := "568B"
	if result != expectedSize {
		t.Errorf("Ожидается размер %s, получено: %s", expectedSize, result)
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
		{1536 * 1024, "1.5MB"},
		{1024 * 1024 * 1024, "1.0GB"},
		{1024 * 1024 * 1024 * 2, "2.0GB"},
		{1024 * 1024 * 1024 * 1024, "1.0TB"},
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
	result, err := GetPathSize("testdata/nonexistent_file.txt", false, false, false)

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
		{"..", true},
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
	resultWith, err := GetPathSize("testdata", false, false, true)
	if err != nil {
		t.Fatalf("GetPathSize вернул ошибку: %v", err)
	}

	// В testdata со скрытыми файлами: 5992 + 14 + 14 = 6020
	expectedWith := "6020B"

	if resultWith != expectedWith {
		t.Errorf("С флагом all ожидается %s, получено: %s", expectedWith, resultWith)
	}
}

// TestGetSize_DirectoryWithoutHidden проверяет, что скрытые файлы игнорируются по умолчанию
func TestGetSize_DirectoryWithoutHidden(t *testing.T) {
	result, err := GetPathSize("testdata/LogViewer", false, false, false)
	if err != nil {
		t.Fatalf("GetPathSize вернул ошибку: %v", err)
	}

	// В LogViewer первого уровня три XML файла: 348 + 65 + 155 = 568 байт
	expectedSize := "568B"
	if result != expectedSize {
		t.Errorf("Ожидается размер %s, получено: %s", expectedSize, result)
	}
}

// TestGetSize_Recursive проверяет рекурсивный подсчёт папки testdata/nested
func TestGetSize_Recursive(t *testing.T) {
	// Без флага recursive - только первый уровень (file1.txt = 6 байт)
	resultNonRecursive, err := GetPathSize("testdata/nested", false, false, false)
	if err != nil {
		t.Fatalf("GetPathSize вернул ошибку: %v", err)
	}

	// С флагом recursive - все файлы (file1.txt 6 + file2.txt 6 + file3.txt 6 = 18 байт)
	resultRecursive, err := GetPathSize("testdata/nested", true, false, false)
	if err != nil {
		t.Fatalf("GetPathSize вернул ошибку: %v", err)
	}

	// Без recursive: только file1.txt (6 байт)
	expectedNonRecursive := "6B"
	if resultNonRecursive != expectedNonRecursive {
		t.Errorf("Без recursive ожидается %s, получено: %s", expectedNonRecursive, resultNonRecursive)
	}

	// С recursive: все файлы (18 байт)
	expectedRecursive := "18B"
	if resultRecursive != expectedRecursive {
		t.Errorf("С recursive ожидается %s, получено: %s", expectedRecursive, resultRecursive)
	}
}

// TestGetSize_HumanFormat проверяет форматирование размера для человека
func TestGetSize_HumanFormat(t *testing.T) {
	result, err := GetPathSize("testdata/65049_melt.log", false, true, false)
	if err != nil {
		t.Fatalf("GetPathSize вернул ошибку: %v", err)
	}

	expectedSize := "4.0KB"
	if result != expectedSize {
		t.Errorf("Ожидается размер %s, получено: %s", expectedSize, result)
	}
}

// TestGetSize_RecursiveWithHuman проверяет рекурсивный подсчёт в человекочитаемом формате
func TestGetSize_RecursiveWithHuman(t *testing.T) {
	result, err := GetPathSize("testdata/nested", true, true, false)
	if err != nil {
		t.Fatalf("GetPathSize вернул ошибку: %v", err)
	}

	// 18 байт = 18B
	expectedSize := "18B"
	if result != expectedSize {
		t.Errorf("Ожидается размер %s, получено: %s", expectedSize, result)
	}
}
