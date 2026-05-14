package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
)

const defaultSourceURL = "https://api.loteriasapi.com/api/v1/results/euromillones/latest"

type config struct {
	dataDir   string
	sourceURL string
	apiKey    string
	fixture   string
	dryRun    bool
}

type yearEntry struct {
	DrawNumber int
	Date       string
	Numbers    []int
	Stars      []int
	ElMillon   string
}

type apiEnvelope struct {
	Data apiResult `json:"data"`
}

type apiResult struct {
	DrawDate    string       `json:"drawDate"`
	DrawNumber  int          `json:"drawNumber"`
	Combination []int        `json:"combination"`
	ResultData  apiExtraData `json:"resultData"`
	ElMillon    string       `json:"elMillon"`
}

type apiExtraData struct {
	Stars    []int  `json:"estrellas"`
	ElMillon string `json:"elMillon"`
}

type latestDraw struct {
	DrawDate   time.Time
	DrawNumber int
	Numbers    []int
	Stars      []int
	ElMillon   string
}

func main() {
	if err := run(parseFlags(os.Args[1:])); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func parseFlags(args []string) config {
	cfg := config{
		dataDir:   envOrDefault("DRAWS_DATA_DIR", "data"),
		sourceURL: envOrDefault("UPDATE_DRAWS_URL", defaultSourceURL),
		apiKey:    os.Getenv("LOTERIAS_API_KEY"),
		fixture:   os.Getenv("UPDATE_DRAWS_FIXTURE"),
	}

	fs := flag.NewFlagSet("update-draws", flag.ExitOnError)
	fs.StringVar(&cfg.dataDir, "data-dir", cfg.dataDir, "directory containing yearly draw JSON files")
	fs.StringVar(&cfg.sourceURL, "source-url", cfg.sourceURL, "latest Euromillones result API URL")
	fs.StringVar(&cfg.apiKey, "api-key", cfg.apiKey, "API key sent as X-API-Key")
	fs.StringVar(&cfg.fixture, "fixture", cfg.fixture, "read API response JSON from a local file instead of the network")
	fs.BoolVar(&cfg.dryRun, "dry-run", false, "validate the update without writing files")
	_ = fs.Parse(args)
	return cfg
}

func run(cfg config) error {
	draw, err := fetchLatestDraw(cfg)
	if err != nil {
		return err
	}

	year := draw.DrawDate.Year()
	path := filepath.Join(cfg.dataDir, fmt.Sprintf("%d.json", year))
	entries, err := readYearEntries(path)
	if err != nil {
		return err
	}

	updated, changed, err := addLatestDraw(entries, draw)
	if err != nil {
		return err
	}
	if !changed {
		fmt.Printf("Sorteo %s ya existe en %s\n", draw.DrawDate.Format("2006-01-02"), path)
		return nil
	}
	added := updated[len(updated)-1]
	if cfg.dryRun {
		fmt.Printf("Dry run: se añadiria el sorteo %d de %s en %s\n", added.DrawNumber, draw.DrawDate.Format("2006-01-02"), path)
		return nil
	}

	if err := os.MkdirAll(cfg.dataDir, 0755); err != nil {
		return err
	}
	if err := writeYearEntries(path, updated); err != nil {
		return err
	}
	fmt.Printf("Añadido sorteo %d de %s en %s\n", added.DrawNumber, draw.DrawDate.Format("2006-01-02"), path)
	return nil
}

func fetchLatestDraw(cfg config) (latestDraw, error) {
	var body []byte
	var err error
	if cfg.fixture != "" {
		body, err = os.ReadFile(cfg.fixture)
		if err != nil {
			return latestDraw{}, err
		}
	} else {
		if strings.TrimSpace(cfg.apiKey) == "" {
			return latestDraw{}, errors.New("falta LOTERIAS_API_KEY; configura el secret en GitHub o usa -fixture para probar en local")
		}
		body, err = fetchURL(cfg.sourceURL, cfg.apiKey)
		if err != nil {
			return latestDraw{}, err
		}
	}
	return parseLatestDraw(body)
}

func fetchURL(sourceURL string, apiKey string) ([]byte, error) {
	request, err := http.NewRequest(http.MethodGet, sourceURL, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Set("X-API-Key", apiKey)
	request.Header.Set("Accept", "application/json")

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	if response.StatusCode < 200 || response.StatusCode > 299 {
		return nil, fmt.Errorf("la API devolvio HTTP %d: %s", response.StatusCode, strings.TrimSpace(string(body)))
	}
	return body, nil
}

func parseLatestDraw(body []byte) (latestDraw, error) {
	var envelope apiEnvelope
	if err := json.Unmarshal(body, &envelope); err != nil {
		return latestDraw{}, err
	}
	result := envelope.Data
	if result.DrawDate == "" && len(result.Combination) == 0 {
		if err := json.Unmarshal(body, &result); err != nil {
			return latestDraw{}, err
		}
	}

	drawDate, err := parseAPIDate(result.DrawDate)
	if err != nil {
		return latestDraw{}, err
	}
	numbers := sortedCopy(result.Combination)
	stars := sortedCopy(result.ResultData.Stars)
	if err := validateDrawValues(numbers, stars); err != nil {
		return latestDraw{}, err
	}

	elMillon := result.ElMillon
	if elMillon == "" {
		elMillon = result.ResultData.ElMillon
	}
	return latestDraw{
		DrawDate:   drawDate,
		DrawNumber: result.DrawNumber,
		Numbers:    numbers,
		Stars:      stars,
		ElMillon:   strings.TrimSpace(elMillon),
	}, nil
}

func parseAPIDate(value string) (time.Time, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return time.Time{}, errors.New("la respuesta no contiene drawDate")
	}
	for _, layout := range []string{time.RFC3339, "2006-01-02"} {
		parsed, err := time.Parse(layout, value)
		if err == nil {
			return time.Date(parsed.Year(), parsed.Month(), parsed.Day(), 0, 0, 0, 0, time.UTC), nil
		}
	}
	return time.Time{}, fmt.Errorf("drawDate invalida %q", value)
}

func readYearEntries(path string) ([]yearEntry, error) {
	content, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return []yearEntry{}, nil
	}
	if err != nil {
		return nil, err
	}
	if len(strings.TrimSpace(string(content))) == 0 {
		return []yearEntry{}, nil
	}

	var raw []struct {
		DrawNumber int    `json:"sorteo"`
		Date       string `json:"fecha"`
		Numbers    []int  `json:"numeros"`
		Stars      []int  `json:"estrellas"`
		ElMillon   string `json:"elMillon"`
	}
	if err := json.Unmarshal(content, &raw); err != nil {
		return nil, err
	}

	entries := make([]yearEntry, 0, len(raw))
	for _, item := range raw {
		entries = append(entries, yearEntry{
			DrawNumber: item.DrawNumber,
			Date:       item.Date,
			Numbers:    sortedCopy(item.Numbers),
			Stars:      sortedCopy(item.Stars),
			ElMillon:   item.ElMillon,
		})
	}
	return entries, nil
}

func addLatestDraw(entries []yearEntry, draw latestDraw) ([]yearEntry, bool, error) {
	newEntry := yearEntry{
		DrawNumber: draw.DrawNumber,
		Date:       formatSpanishDate(draw.DrawDate),
		Numbers:    draw.Numbers,
		Stars:      draw.Stars,
		ElMillon:   draw.ElMillon,
	}

	for _, entry := range entries {
		sameDate := entry.Date == newEntry.Date
		sameDrawNumber := newEntry.DrawNumber != 0 && entry.DrawNumber == newEntry.DrawNumber
		if sameDate || sameDrawNumber {
			if newEntry.DrawNumber == 0 {
				newEntry.DrawNumber = entry.DrawNumber
			}
			if !sameEntry(entry, newEntry) {
				return nil, false, fmt.Errorf("el sorteo %d/%s ya existe con datos diferentes", newEntry.DrawNumber, newEntry.Date)
			}
			return entries, false, nil
		}
	}

	if newEntry.DrawNumber == 0 {
		newEntry.DrawNumber = nextDrawNumber(entries)
	}
	return append(entries, newEntry), true, nil
}

func writeYearEntries(path string, entries []yearEntry) error {
	content := formatYearEntries(entries)
	tempPath := path + ".tmp"
	if err := os.WriteFile(tempPath, content, 0644); err != nil {
		return err
	}
	return os.Rename(tempPath, path)
}

func formatYearEntries(entries []yearEntry) []byte {
	var buffer bytes.Buffer
	buffer.WriteString("[\n")
	for index, entry := range entries {
		buffer.WriteString("  {\n")
		fmt.Fprintf(&buffer, "    \"sorteo\": %d,\n", entry.DrawNumber)
		fmt.Fprintf(&buffer, "    \"fecha\": %q,\n", entry.Date)
		fmt.Fprintf(&buffer, "    \"numeros\": %s,\n", formatIntArray(entry.Numbers))
		fmt.Fprintf(&buffer, "    \"estrellas\": %s", formatIntArray(entry.Stars))
		if entry.ElMillon != "" {
			fmt.Fprintf(&buffer, ",\n    \"elMillon\": %q\n", entry.ElMillon)
		} else {
			buffer.WriteString("\n")
		}
		buffer.WriteString("  }")
		if index < len(entries)-1 {
			buffer.WriteString(",")
		}
		buffer.WriteString("\n")
	}
	buffer.WriteString("]\n")
	return buffer.Bytes()
}

func formatIntArray(values []int) string {
	parts := make([]string, len(values))
	for index, value := range values {
		parts[index] = strconv.Itoa(value)
	}
	return "[" + strings.Join(parts, ", ") + "]"
}

func formatSpanishDate(value time.Time) string {
	months := []string{"ene", "feb", "mar", "abr", "may", "jun", "jul", "ago", "sep", "oct", "nov", "dic"}
	return fmt.Sprintf("%d-%s", value.Day(), months[int(value.Month())-1])
}

func nextDrawNumber(entries []yearEntry) int {
	maxDrawNumber := 0
	for _, entry := range entries {
		if entry.DrawNumber > maxDrawNumber {
			maxDrawNumber = entry.DrawNumber
		}
	}
	return maxDrawNumber + 1
}

func validateDrawValues(numbers []int, stars []int) error {
	if len(numbers) != 5 {
		return fmt.Errorf("se esperaban 5 numeros, llegaron %d", len(numbers))
	}
	if len(stars) != 2 {
		return fmt.Errorf("se esperaban 2 estrellas, llegaron %d", len(stars))
	}
	for _, number := range numbers {
		if number < 1 || number > 50 {
			return fmt.Errorf("numero fuera de rango: %d", number)
		}
	}
	for _, star := range stars {
		if star < 1 || star > 12 {
			return fmt.Errorf("estrella fuera de rango: %d", star)
		}
	}
	return nil
}

func sameEntry(a yearEntry, b yearEntry) bool {
	return a.DrawNumber == b.DrawNumber &&
		a.Date == b.Date &&
		intsEqual(sortedCopy(a.Numbers), sortedCopy(b.Numbers)) &&
		intsEqual(sortedCopy(a.Stars), sortedCopy(b.Stars)) &&
		strings.TrimSpace(a.ElMillon) == strings.TrimSpace(b.ElMillon)
}

func sortedCopy(values []int) []int {
	copyValues := append([]int{}, values...)
	sort.Ints(copyValues)
	return copyValues
}

func intsEqual(a []int, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for index := range a {
		if a[index] != b[index] {
			return false
		}
	}
	return true
}

func envOrDefault(name string, fallback string) string {
	value := strings.TrimSpace(os.Getenv(name))
	if value == "" {
		return fallback
	}
	return value
}
