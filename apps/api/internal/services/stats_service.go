package services

import (
	"sort"

	"euromillones/internal/models"
)

type FrequencyStat struct {
	Value int `json:"value"`
	Count int `json:"count"`
}

type FrequenciesResponse struct {
	Numbers []FrequencyStat `json:"numbers"`
	Stars   []FrequencyStat `json:"stars"`
}

type PositionFrequency struct {
	Position string          `json:"position"`
	Values   []FrequencyStat `json:"values"`
}

type PositionsResponse struct {
	Numbers []PositionFrequency `json:"numbers"`
	Stars   []PositionFrequency `json:"stars"`
}

type HotColdResponse struct {
	HotNumbers  []FrequencyStat `json:"hotNumbers"`
	ColdNumbers []FrequencyStat `json:"coldNumbers"`
	HotStars    []FrequencyStat `json:"hotStars"`
	ColdStars   []FrequencyStat `json:"coldStars"`
}

type DelayStat struct {
	Value int `json:"value"`
	Delay int `json:"delay"`
}

type DelaysResponse struct {
	Numbers []DelayStat `json:"numbers"`
	Stars   []DelayStat `json:"stars"`
}

type PairStat struct {
	A     int `json:"a"`
	B     int `json:"b"`
	Count int `json:"count"`
}

type DashboardStats struct {
	TotalDraws         int            `json:"totalDraws"`
	MostFrequentNumber *FrequencyStat `json:"mostFrequentNumber"`
	MostFrequentStar   *FrequencyStat `json:"mostFrequentStar"`
	MostDelayedNumber  *DelayStat     `json:"mostDelayedNumber"`
	LastDrawDate       string         `json:"lastDrawDate"`
}

type StatsService struct {
	drawService *DrawService
}

func NewStatsService(drawService *DrawService) *StatsService {
	return &StatsService{drawService: drawService}
}

func (s *StatsService) Frequencies() (FrequenciesResponse, error) {
	draws, err := s.drawService.All()
	if err != nil {
		return FrequenciesResponse{}, err
	}
	numberCounts, starCounts := buildFrequencyMaps(draws)

	return FrequenciesResponse{
		Numbers: frequencySlice(numberCounts, 1, 50),
		Stars:   frequencySlice(starCounts, 1, 12),
	}, nil
}

func (s *StatsService) Positions() (PositionsResponse, error) {
	draws, err := s.drawService.All()
	if err != nil {
		return PositionsResponse{}, err
	}

	numberPositions := []PositionFrequency{}
	for index, name := range []string{"n1", "n2", "n3", "n4", "n5"} {
		counts := emptyCounts(1, 50)
		for _, draw := range draws {
			counts[draw.Numbers()[index]]++
		}
		numberPositions = append(numberPositions, PositionFrequency{Position: name, Values: frequencySlice(counts, 1, 50)})
	}

	starPositions := []PositionFrequency{}
	for index, name := range []string{"star1", "star2"} {
		counts := emptyCounts(1, 12)
		for _, draw := range draws {
			counts[draw.Stars()[index]]++
		}
		starPositions = append(starPositions, PositionFrequency{Position: name, Values: frequencySlice(counts, 1, 12)})
	}

	return PositionsResponse{Numbers: numberPositions, Stars: starPositions}, nil
}

func (s *StatsService) HotCold(limit int) (HotColdResponse, error) {
	if limit < 1 {
		limit = 10
	}
	frequencies, err := s.Frequencies()
	if err != nil {
		return HotColdResponse{}, err
	}

	return HotColdResponse{
		HotNumbers:  topByCount(frequencies.Numbers, limit, true),
		ColdNumbers: topByCount(frequencies.Numbers, limit, false),
		HotStars:    topByCount(frequencies.Stars, 5, true),
		ColdStars:   topByCount(frequencies.Stars, 5, false),
	}, nil
}

func (s *StatsService) Delays() (DelaysResponse, error) {
	draws, err := s.drawService.All()
	if err != nil {
		return DelaysResponse{}, err
	}

	numberDelays := buildDelays(draws, 1, 50, false)
	starDelays := buildDelays(draws, 1, 12, true)
	return DelaysResponse{Numbers: numberDelays, Stars: starDelays}, nil
}

func (s *StatsService) Pairs(limit int) ([]PairStat, error) {
	if limit < 1 {
		limit = 20
	}
	draws, err := s.drawService.All()
	if err != nil {
		return nil, err
	}

	counts := map[[2]int]int{}
	for _, draw := range draws {
		numbers := draw.Numbers()
		for i := 0; i < len(numbers); i++ {
			for j := i + 1; j < len(numbers); j++ {
				counts[[2]int{numbers[i], numbers[j]}]++
			}
		}
	}

	pairs := []PairStat{}
	for pair, count := range counts {
		pairs = append(pairs, PairStat{A: pair[0], B: pair[1], Count: count})
	}
	sort.Slice(pairs, func(i, j int) bool {
		if pairs[i].Count == pairs[j].Count {
			if pairs[i].A == pairs[j].A {
				return pairs[i].B < pairs[j].B
			}
			return pairs[i].A < pairs[j].A
		}
		return pairs[i].Count > pairs[j].Count
	})
	if len(pairs) > limit {
		pairs = pairs[:limit]
	}

	return pairs, nil
}

func (s *StatsService) Dashboard() (DashboardStats, error) {
	draws, err := s.drawService.All()
	if err != nil {
		return DashboardStats{}, err
	}
	frequencies, err := s.Frequencies()
	if err != nil {
		return DashboardStats{}, err
	}
	delays, err := s.Delays()
	if err != nil {
		return DashboardStats{}, err
	}

	var lastDrawDate string
	if len(draws) > 0 {
		lastDrawDate = draws[0].DrawDate.Format("2006-01-02")
	}

	return DashboardStats{
		TotalDraws:         len(draws),
		MostFrequentNumber: firstFrequency(topByCount(frequencies.Numbers, 1, true)),
		MostFrequentStar:   firstFrequency(topByCount(frequencies.Stars, 1, true)),
		MostDelayedNumber:  firstDelay(topDelays(delays.Numbers, 1)),
		LastDrawDate:       lastDrawDate,
	}, nil
}

func buildFrequencyMaps(draws []models.Draw) (map[int]int, map[int]int) {
	numberCounts := emptyCounts(1, 50)
	starCounts := emptyCounts(1, 12)
	for _, draw := range draws {
		for _, number := range draw.Numbers() {
			numberCounts[number]++
		}
		for _, star := range draw.Stars() {
			starCounts[star]++
		}
	}
	return numberCounts, starCounts
}

func emptyCounts(min int, max int) map[int]int {
	counts := map[int]int{}
	for value := min; value <= max; value++ {
		counts[value] = 0
	}
	return counts
}

func frequencySlice(counts map[int]int, min int, max int) []FrequencyStat {
	stats := []FrequencyStat{}
	for value := min; value <= max; value++ {
		stats = append(stats, FrequencyStat{Value: value, Count: counts[value]})
	}
	return stats
}

func topByCount(values []FrequencyStat, limit int, descending bool) []FrequencyStat {
	copyValues := append([]FrequencyStat{}, values...)
	sort.Slice(copyValues, func(i, j int) bool {
		if copyValues[i].Count == copyValues[j].Count {
			return copyValues[i].Value < copyValues[j].Value
		}
		if descending {
			return copyValues[i].Count > copyValues[j].Count
		}
		return copyValues[i].Count < copyValues[j].Count
	})
	if len(copyValues) > limit {
		return copyValues[:limit]
	}
	return copyValues
}

func buildDelays(draws []models.Draw, min int, max int, stars bool) []DelayStat {
	delays := []DelayStat{}
	for value := min; value <= max; value++ {
		delay := len(draws)
		for index, draw := range draws {
			values := draw.Numbers()
			if stars {
				values = draw.Stars()
			}
			if contains(values, value) {
				delay = index
				break
			}
		}
		delays = append(delays, DelayStat{Value: value, Delay: delay})
	}
	return delays
}

func topDelays(values []DelayStat, limit int) []DelayStat {
	copyValues := append([]DelayStat{}, values...)
	sort.Slice(copyValues, func(i, j int) bool {
		if copyValues[i].Delay == copyValues[j].Delay {
			return copyValues[i].Value < copyValues[j].Value
		}
		return copyValues[i].Delay > copyValues[j].Delay
	})
	if len(copyValues) > limit {
		return copyValues[:limit]
	}
	return copyValues
}

func firstFrequency(values []FrequencyStat) *FrequencyStat {
	if len(values) == 0 {
		return nil
	}
	return &values[0]
}

func firstDelay(values []DelayStat) *DelayStat {
	if len(values) == 0 {
		return nil
	}
	return &values[0]
}

func contains(values []int, target int) bool {
	for _, value := range values {
		if value == target {
			return true
		}
	}
	return false
}
