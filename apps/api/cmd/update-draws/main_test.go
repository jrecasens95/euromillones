package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRunAddsLatestDrawFromFixture(t *testing.T) {
	tempDir := t.TempDir()
	dataDir := filepath.Join(tempDir, "data")
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		t.Fatalf("mkdir data: %v", err)
	}
	if err := os.WriteFile(filepath.Join(dataDir, "2026.json"), []byte(`[
  {
    "sorteo": 38,
    "fecha": "12-may",
    "numeros": [4, 26, 32, 35, 36],
    "estrellas": [5, 7],
    "elMillon": "BZN58201"
  }
]
`), 0644); err != nil {
		t.Fatalf("write year file: %v", err)
	}

	fixture := filepath.Join(tempDir, "latest.json")
	if err := os.WriteFile(fixture, []byte(`{
  "success": true,
  "data": {
    "drawDate": "2026-05-15",
    "combination": [50, 1, 23, 8, 41],
    "resultData": {
      "estrellas": [12, 2],
      "elMillon": "BCD12345"
    }
  }
}`), 0644); err != nil {
		t.Fatalf("write fixture: %v", err)
	}

	if err := run(config{dataDir: dataDir, fixture: fixture}); err != nil {
		t.Fatalf("run returned error: %v", err)
	}

	content, err := os.ReadFile(filepath.Join(dataDir, "2026.json"))
	if err != nil {
		t.Fatalf("read updated year file: %v", err)
	}
	updated := string(content)
	for _, expected := range []string{
		`"sorteo": 39`,
		`"fecha": "15-may"`,
		`"numeros": [1, 8, 23, 41, 50]`,
		`"estrellas": [2, 12]`,
		`"elMillon": "BCD12345"`,
	} {
		if !strings.Contains(updated, expected) {
			t.Fatalf("updated file does not contain %s:\n%s", expected, updated)
		}
	}
}

func TestRunIsIdempotentWhenDrawAlreadyExists(t *testing.T) {
	tempDir := t.TempDir()
	dataDir := filepath.Join(tempDir, "data")
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		t.Fatalf("mkdir data: %v", err)
	}
	yearPath := filepath.Join(dataDir, "2026.json")
	original := `[
  {
    "sorteo": 39,
    "fecha": "15-may",
    "numeros": [1, 8, 23, 41, 50],
    "estrellas": [2, 12],
    "elMillon": "BCD12345"
  }
]
`
	if err := os.WriteFile(yearPath, []byte(original), 0644); err != nil {
		t.Fatalf("write year file: %v", err)
	}

	fixture := filepath.Join(tempDir, "latest.json")
	if err := os.WriteFile(fixture, []byte(`{
  "data": {
    "drawDate": "2026-05-15",
    "combination": [1, 8, 23, 41, 50],
    "resultData": {
      "estrellas": [2, 12],
      "elMillon": "BCD12345"
    }
  }
}`), 0644); err != nil {
		t.Fatalf("write fixture: %v", err)
	}

	if err := run(config{dataDir: dataDir, fixture: fixture}); err != nil {
		t.Fatalf("run returned error: %v", err)
	}

	content, err := os.ReadFile(yearPath)
	if err != nil {
		t.Fatalf("read year file: %v", err)
	}
	if string(content) != original {
		t.Fatalf("expected file to remain unchanged:\n%s", content)
	}
}

func TestRunAcceptsCombinationWithNumbersAndStars(t *testing.T) {
	tempDir := t.TempDir()
	dataDir := filepath.Join(tempDir, "data")
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		t.Fatalf("mkdir data: %v", err)
	}

	fixture := filepath.Join(tempDir, "latest.json")
	if err := os.WriteFile(fixture, []byte(`{
  "success": true,
  "data": {
    "id": "69b07e36e19fc394c9ce118b",
    "game": {
      "slug": "euromillones",
      "name": "Euromillones"
    },
    "drawId": "1310702040",
    "drawDate": "2026-05-19",
    "dayOfWeek": "martes",
    "year": 2026,
    "status": "COMPLETED",
    "jackpot": "9200000000",
    "jackpotFormatted": "92.000.000,00 €",
    "combination": [12, 20, 45, 38, 2, 5, 2],
    "resultData": {
      "estrellas": [2, 5],
      "millon": {
        "gameid": "LAMI",
        "relsorteoid_asociado": "1310711040",
        "combinacion": "CGM77305",
        "importe": 1000000,
        "activo": "1"
      },
      "escrutinioMillon": [],
      "combinacionRaw": " 02 - 12 - 20 - 38 - 45 - 02 - 05"
    },
    "prizes": [],
    "documents": []
  },
  "timestamp": "2026-05-21T10:52:09.008Z"
}`), 0644); err != nil {
		t.Fatalf("write fixture: %v", err)
	}

	if err := run(config{dataDir: dataDir, fixture: fixture}); err != nil {
		t.Fatalf("run returned error: %v", err)
	}

	content, err := os.ReadFile(filepath.Join(dataDir, "2026.json"))
	if err != nil {
		t.Fatalf("read updated year file: %v", err)
	}
	updated := string(content)
	for _, expected := range []string{
		`"fecha": "19-may"`,
		`"numeros": [2, 12, 20, 38, 45]`,
		`"estrellas": [2, 5]`,
		`"elMillon": "CGM77305"`,
	} {
		if !strings.Contains(updated, expected) {
			t.Fatalf("updated file does not contain %s:\n%s", expected, updated)
		}
	}
}

func TestRunFailsWhenExistingDrawHasDifferentValues(t *testing.T) {
	tempDir := t.TempDir()
	dataDir := filepath.Join(tempDir, "data")
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		t.Fatalf("mkdir data: %v", err)
	}
	if err := os.WriteFile(filepath.Join(dataDir, "2026.json"), []byte(`[
  {
    "sorteo": 39,
    "fecha": "15-may",
    "numeros": [1, 8, 23, 41, 50],
    "estrellas": [2, 12],
    "elMillon": "BCD12345"
  }
]
`), 0644); err != nil {
		t.Fatalf("write year file: %v", err)
	}

	fixture := filepath.Join(tempDir, "latest.json")
	if err := os.WriteFile(fixture, []byte(`{
  "data": {
    "drawDate": "2026-05-15",
    "combination": [1, 8, 23, 41, 49],
    "resultData": {
      "estrellas": [2, 12],
      "elMillon": "BCD12345"
    }
  }
}`), 0644); err != nil {
		t.Fatalf("write fixture: %v", err)
	}

	err := run(config{dataDir: dataDir, fixture: fixture})
	if err == nil || !strings.Contains(err.Error(), "datos diferentes") {
		t.Fatalf("expected different values error, got %v", err)
	}
}
