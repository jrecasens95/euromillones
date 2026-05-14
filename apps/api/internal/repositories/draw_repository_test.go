package repositories

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestParseYearDrawDateUsesFallbackYearForSpanishMonthWithoutYear(t *testing.T) {
	parsed, err := parseYearDrawDate("29-dic", 2006)
	if err != nil {
		t.Fatalf("parseYearDrawDate returned error: %v", err)
	}

	expected := time.Date(2006, time.December, 29, 0, 0, 0, 0, time.UTC)
	if !parsed.Equal(expected) {
		t.Fatalf("expected %s, got %s", expected, parsed)
	}
}

func TestParseYearDrawDateUsesFallbackYearForUnpaddedSpanishMonthWithoutYear(t *testing.T) {
	parsed, err := parseYearDrawDate("7-ene", 2011)
	if err != nil {
		t.Fatalf("parseYearDrawDate returned error: %v", err)
	}

	expected := time.Date(2011, time.January, 7, 0, 0, 0, 0, time.UTC)
	if !parsed.Equal(expected) {
		t.Fatalf("expected %s, got %s", expected, parsed)
	}
}

func TestParseYearDrawDateParsesSpanishMonthWithYear(t *testing.T) {
	parsed, err := parseYearDrawDate("06-ene-2006", 1999)
	if err != nil {
		t.Fatalf("parseYearDrawDate returned error: %v", err)
	}

	expected := time.Date(2006, time.January, 6, 0, 0, 0, 0, time.UTC)
	if !parsed.Equal(expected) {
		t.Fatalf("expected %s, got %s", expected, parsed)
	}
}

func TestReadYearFileHandlesBoundaryYearDraws(t *testing.T) {
	path := filepath.Join(t.TempDir(), "2015.json")
	content := `[
  {"sorteo": 104, "fecha": "30-dic", "numeros": [6, 18, 39, 50, 44], "estrellas": [11, 8]},
  {"sorteo": 1, "fecha": "2-ene", "numeros": [24, 25, 49, 28, 22], "estrellas": [3, 6]},
  {"sorteo": "2016/001", "fecha": "1-ene", "numeros": [44, 37, 38, 39, 4], "estrellas": [4, 7]}
]`
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("write fixture: %v", err)
	}

	draws, err := readYearFile(path, 2015)
	if err != nil {
		t.Fatalf("readYearFile returned error: %v", err)
	}
	if len(draws) != 3 {
		t.Fatalf("expected 3 draws, got %d", len(draws))
	}

	expectedDates := []time.Time{
		time.Date(2014, time.December, 30, 0, 0, 0, 0, time.UTC),
		time.Date(2015, time.January, 2, 0, 0, 0, 0, time.UTC),
		time.Date(2016, time.January, 1, 0, 0, 0, 0, time.UTC),
	}
	expectedIDs := []uint{idFor(2014, 104), idFor(2015, 1), idFor(2016, 1)}
	for index := range draws {
		if !draws[index].DrawDate.Equal(expectedDates[index]) {
			t.Fatalf("draw %d expected date %s, got %s", index, expectedDates[index], draws[index].DrawDate)
		}
		if draws[index].ID != expectedIDs[index] {
			t.Fatalf("draw %d expected ID %d, got %d", index, expectedIDs[index], draws[index].ID)
		}
	}
}
