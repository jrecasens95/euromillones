package repositories

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"euromillones/internal/models"
)

var ErrDrawNotFound = errors.New("sorteo no encontrado")

type PaginatedDraws struct {
	Draws []models.Draw `json:"draws"`
	Total int64         `json:"total"`
	Page  int           `json:"page"`
	Limit int           `json:"limit"`
}

type yearDraw struct {
	DrawNumber int    `json:"sorteo"`
	Date       string `json:"fecha"`
	Numbers    []int  `json:"numeros"`
	Stars      []int  `json:"estrellas"`
	DrawYear   int    `json:"-"`
}

func (d *yearDraw) UnmarshalJSON(data []byte) error {
	type rawYearDraw struct {
		DrawNumber json.RawMessage `json:"sorteo"`
		Date       string          `json:"fecha"`
		Numbers    []int           `json:"numeros"`
		Stars      []int           `json:"estrellas"`
	}

	var raw rawYearDraw
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	drawNumber, drawYear, err := parseYearDrawNumber(raw.DrawNumber)
	if err != nil {
		return err
	}

	d.DrawNumber = drawNumber
	d.Date = raw.Date
	d.Numbers = raw.Numbers
	d.Stars = raw.Stars
	d.DrawYear = drawYear
	return nil
}

type DrawRepository struct {
	dataDir string
	mu      sync.RWMutex
	draws   []models.Draw
	years   map[int]bool
}

func NewDrawRepository(dataDir string) (*DrawRepository, error) {
	repo := &DrawRepository{dataDir: dataDir}
	if err := repo.load(); err != nil {
		return nil, err
	}
	return repo, nil
}

func (r *DrawRepository) load() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if err := os.MkdirAll(r.dataDir, 0755); err != nil {
		return err
	}

	files, err := filepath.Glob(filepath.Join(r.dataDir, "*.json"))
	if err != nil {
		return err
	}

	r.draws = []models.Draw{}
	r.years = map[int]bool{}
	for _, file := range files {
		year, err := yearFromFile(file)
		if err != nil {
			continue
		}
		r.years[year] = true
		draws, err := readYearFile(file, year)
		if err != nil {
			return err
		}
		r.draws = append(r.draws, draws...)
	}

	r.sortLocked()
	return nil
}

func (r *DrawRepository) List(page int, limit int) (PaginatedDraws, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	total := int64(len(r.draws))
	start := (page - 1) * limit
	if start > len(r.draws) {
		start = len(r.draws)
	}
	end := start + limit
	if end > len(r.draws) {
		end = len(r.draws)
	}

	draws := append([]models.Draw{}, r.draws[start:end]...)
	return PaginatedDraws{Draws: draws, Total: total, Page: page, Limit: limit}, nil
}

func (r *DrawRepository) All() ([]models.Draw, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return append([]models.Draw{}, r.draws...), nil
}

func (r *DrawRepository) Find(id uint) (models.Draw, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, draw := range r.draws {
		if draw.ID == id {
			return draw, nil
		}
	}
	return models.Draw{}, ErrDrawNotFound
}

func (r *DrawRepository) Create(draw models.Draw) (models.Draw, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if err := r.ensureDateAvailableLocked(draw.DrawDate, 0); err != nil {
		return models.Draw{}, err
	}

	draw.DrawNumber = r.nextDrawNumberLocked(draw.DrawDate.Year())
	draw.ID = idFor(draw.DrawDate.Year(), draw.DrawNumber)
	now := time.Now()
	draw.CreatedAt = now
	draw.UpdatedAt = now

	r.draws = append(r.draws, draw)
	r.sortLocked()
	return draw, r.persistLocked()
}

func (r *DrawRepository) Update(id uint, draw models.Draw) (models.Draw, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	index := r.findIndexLocked(id)
	if index == -1 {
		return models.Draw{}, ErrDrawNotFound
	}
	if err := r.ensureDateAvailableLocked(draw.DrawDate, id); err != nil {
		return models.Draw{}, err
	}

	current := r.draws[index]
	oldYear := current.DrawDate.Year()
	newYear := draw.DrawDate.Year()
	if oldYear != newYear {
		current.DrawNumber = r.nextDrawNumberLocked(newYear)
		current.ID = idFor(newYear, current.DrawNumber)
	}

	current.DrawDate = draw.DrawDate
	current.N1 = draw.N1
	current.N2 = draw.N2
	current.N3 = draw.N3
	current.N4 = draw.N4
	current.N5 = draw.N5
	current.Star1 = draw.Star1
	current.Star2 = draw.Star2
	current.UpdatedAt = time.Now()

	r.draws[index] = current
	r.sortLocked()
	return current, r.persistLocked()
}

func (r *DrawRepository) Delete(id uint) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	index := r.findIndexLocked(id)
	if index == -1 {
		return ErrDrawNotFound
	}
	r.draws = append(r.draws[:index], r.draws[index+1:]...)
	return r.persistLocked()
}

func readYearFile(file string, year int) ([]models.Draw, error) {
	content, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}
	if len(strings.TrimSpace(string(content))) == 0 {
		return []models.Draw{}, nil
	}

	rawDraws := []yearDraw{}
	if err := json.Unmarshal(content, &rawDraws); err != nil {
		return nil, fmt.Errorf("%s: %w", file, err)
	}

	draws := make([]models.Draw, 0, len(rawDraws))
	for index, raw := range rawDraws {
		drawDate, err := parseYearDrawDate(raw.Date, year)
		if err != nil {
			return nil, fmt.Errorf("%s sorteo %d: %w", file, raw.DrawNumber, err)
		}
		drawYear := drawDate.Year()
		if raw.DrawYear != 0 {
			drawYear = raw.DrawYear
		} else if index == 0 && !dateHasExplicitYear(raw.Date) && drawDate.Month() == time.December {
			drawYear = year - 1
		}
		if drawYear != drawDate.Year() {
			drawDate = time.Date(drawYear, drawDate.Month(), drawDate.Day(), 0, 0, 0, 0, time.UTC)
		}
		if len(raw.Numbers) != 5 || len(raw.Stars) != 2 {
			return nil, fmt.Errorf("%s sorteo %d: estructura de numeros/estrellas invalida", file, raw.DrawNumber)
		}

		numbers := append([]int{}, raw.Numbers...)
		stars := append([]int{}, raw.Stars...)
		sort.Ints(numbers)
		sort.Ints(stars)

		draws = append(draws, models.Draw{
			ID:         idFor(drawYear, raw.DrawNumber),
			DrawNumber: raw.DrawNumber,
			DrawDate:   drawDate,
			N1:         numbers[0],
			N2:         numbers[1],
			N3:         numbers[2],
			N4:         numbers[3],
			N5:         numbers[4],
			Star1:      stars[0],
			Star2:      stars[1],
		})
	}

	return draws, nil
}

func (r *DrawRepository) persistLocked() error {
	grouped := map[int][]yearDraw{}
	for _, draw := range r.draws {
		year := draw.DrawDate.Year()
		r.years[year] = true
		grouped[year] = append(grouped[year], yearDraw{
			DrawNumber: draw.DrawNumber,
			Date:       draw.DrawDate.Format("02-01-2006"),
			Numbers:    draw.Numbers(),
			Stars:      draw.Stars(),
		})
	}

	for year := range r.years {
		draws := grouped[year]
		sort.Slice(draws, func(i, j int) bool {
			return draws[i].DrawNumber < draws[j].DrawNumber
		})
		content, err := json.MarshalIndent(draws, "", "  ")
		if err != nil {
			return err
		}

		path := filepath.Join(r.dataDir, fmt.Sprintf("%d.json", year))
		tempPath := path + ".tmp"
		if err := os.WriteFile(tempPath, content, 0644); err != nil {
			return err
		}
		if err := os.Rename(tempPath, path); err != nil {
			return err
		}
	}

	return nil
}

func (r *DrawRepository) ensureDateAvailableLocked(drawDate time.Time, excludeID uint) error {
	for _, draw := range r.draws {
		if draw.ID != excludeID && sameDate(draw.DrawDate, drawDate) {
			return errors.New("ya existe un sorteo con esa fecha")
		}
	}
	return nil
}

func (r *DrawRepository) findIndexLocked(id uint) int {
	for index, draw := range r.draws {
		if draw.ID == id {
			return index
		}
	}
	return -1
}

func (r *DrawRepository) nextDrawNumberLocked(year int) int {
	maxDrawNumber := 0
	for _, draw := range r.draws {
		if draw.DrawDate.Year() == year && draw.DrawNumber > maxDrawNumber {
			maxDrawNumber = draw.DrawNumber
		}
	}
	return maxDrawNumber + 1
}

func (r *DrawRepository) sortLocked() {
	sort.Slice(r.draws, func(i, j int) bool {
		return r.draws[i].DrawDate.After(r.draws[j].DrawDate)
	})
}

func yearFromFile(file string) (int, error) {
	name := strings.TrimSuffix(filepath.Base(file), filepath.Ext(file))
	return strconv.Atoi(name)
}

func parseYearDrawNumber(raw json.RawMessage) (int, int, error) {
	var number int
	if err := json.Unmarshal(raw, &number); err == nil {
		return number, 0, nil
	}

	var value string
	if err := json.Unmarshal(raw, &value); err != nil {
		return 0, 0, fmt.Errorf("sorteo invalido %s", string(raw))
	}

	value = strings.TrimSpace(value)
	parts := strings.Split(value, "/")
	if len(parts) == 2 {
		drawYear, err := strconv.Atoi(parts[0])
		if err != nil {
			return 0, 0, fmt.Errorf("sorteo invalido %q", value)
		}
		drawNumber, err := strconv.Atoi(parts[1])
		if err != nil {
			return 0, 0, fmt.Errorf("sorteo invalido %q", value)
		}
		return drawNumber, drawYear, nil
	}

	drawNumber, err := strconv.Atoi(value)
	if err != nil {
		return 0, 0, fmt.Errorf("sorteo invalido %q", value)
	}
	return drawNumber, 0, nil
}

func parseYearDrawDate(value string, fallbackYear int) (time.Time, error) {
	value = normalizeSpanishMonthDate(value)

	for _, layout := range []string{"02-01-2006", "2-1-2006", "02-01-06"} {
		parsed, err := time.Parse(layout, value)
		if err == nil {
			return parsed, nil
		}
	}

	for _, layout := range []string{"02-01", "2-1"} {
		parsed, err := time.Parse(layout, value)
		if err == nil {
			return time.Date(fallbackYear, parsed.Month(), parsed.Day(), 0, 0, 0, 0, time.UTC), nil
		}
	}

	return time.Time{}, fmt.Errorf("fecha invalida %q", value)
}

func dateHasExplicitYear(value string) bool {
	parts := strings.Split(strings.TrimSpace(value), "-")
	return len(parts) == 3
}

func normalizeSpanishMonthDate(value string) string {
	value = strings.TrimSpace(strings.ToLower(value))
	parts := strings.Split(value, "-")
	months := map[string]string{
		"ene": "01",
		"feb": "02",
		"mar": "03",
		"abr": "04",
		"may": "05",
		"jun": "06",
		"jul": "07",
		"ago": "08",
		"sep": "09",
		"oct": "10",
		"nov": "11",
		"dic": "12",
	}
	for index, part := range parts {
		if month, ok := months[part]; ok {
			parts[index] = month
		}
	}
	return strings.Join(parts, "-")
}

func idFor(year int, drawNumber int) uint {
	return uint(year*10000 + drawNumber)
}

func sameDate(a time.Time, b time.Time) bool {
	return a.Format("2006-01-02") == b.Format("2006-01-02")
}
