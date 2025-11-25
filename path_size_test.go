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

	// Проверяем, что размер равен 1863 байтам
	expectedSize := "1863"
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

	// Проверяем, что размер равен 4129 байтам
	expectedSize := "4129"
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
	expectedSize := "568"
	if parts[0] != expectedSize {
		t.Errorf("Ожидается размер %s (348 + 65 + 155), получено: %s", expectedSize, parts[0])
	}

	if parts[1] != testDir {
		t.Errorf("Ожидается путь %s, получено: %s", testDir, parts[1])
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
