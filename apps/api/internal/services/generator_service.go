package services

import (
	"errors"
	"math"
	"math/rand"
	"sort"
	"time"
)

type GenerateRequest struct {
	Strategy string `json:"strategy"`
	Count    int    `json:"count"`
}

type GeneratedCombination struct {
	Numbers     []int  `json:"numbers"`
	Stars       []int  `json:"stars"`
	Strategy    string `json:"strategy"`
	Score       int    `json:"score"`
	Explanation string `json:"explanation"`
}

type GenerateResponse struct {
	Combinations []GeneratedCombination `json:"combinations"`
}

type GeneratorService struct {
	stats *StatsService
	rng   *rand.Rand
}

func NewGeneratorService(stats *StatsService) *GeneratorService {
	return &GeneratorService{
		stats: stats,
		rng:   rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

func (s *GeneratorService) Generate(request GenerateRequest) (GenerateResponse, error) {
	if request.Count < 1 {
		request.Count = 1
	}
	if request.Count > 50 {
		request.Count = 50
	}

	strategy := request.Strategy
	if strategy == "" {
		strategy = "balanced"
	}
	if !validStrategy(strategy) {
		return GenerateResponse{}, errors.New("estrategia de generación no soportada")
	}

	frequencies, err := s.stats.Frequencies()
	if err != nil {
		return GenerateResponse{}, err
	}
	delays, err := s.stats.Delays()
	if err != nil {
		return GenerateResponse{}, err
	}

	combinations := []GeneratedCombination{}
	for len(combinations) < request.Count {
		numbers := s.pickNumbers(strategy, frequencies.Numbers, delays.Numbers, 5, 50)
		stars := s.pickNumbers(strategy, frequencies.Stars, delays.Stars, 2, 12)
		sort.Ints(numbers)
		sort.Ints(stars)
		score := scoreCombination(numbers, stars, strategy, frequencies, delays)

		combinations = append(combinations, GeneratedCombination{
			Numbers:     numbers,
			Stars:       stars,
			Strategy:    strategy,
			Score:       score,
			Explanation: explanationFor(strategy),
		})
	}

	return GenerateResponse{Combinations: combinations}, nil
}

func validStrategy(strategy string) bool {
	switch strategy {
	case "hot", "cold", "delayed", "balanced", "random", "anti_human":
		return true
	default:
		return false
	}
}

func (s *GeneratorService) pickNumbers(strategy string, frequencies []FrequencyStat, delays []DelayStat, amount int, max int) []int {
	selected := []int{}
	weights := map[int]int{}

	for _, frequency := range frequencies {
		weights[frequency.Value] = 1
		switch strategy {
		case "hot":
			weights[frequency.Value] += frequency.Count * 3
		case "cold":
			weights[frequency.Value] += maxFrequency(frequencies) - frequency.Count + 1
		case "balanced":
			weights[frequency.Value] += frequency.Count + 1
		}
	}

	for _, delay := range delays {
		switch strategy {
		case "delayed":
			weights[delay.Value] += delay.Delay*4 + 1
		case "balanced":
			weights[delay.Value] += delay.Delay + 1
		}
	}

	for len(selected) < amount {
		candidate := s.weightedPick(weights, max)
		if strategy == "anti_human" && max == 50 && len(selected) >= 3 && candidate <= 31 {
			candidate = 32 + s.rng.Intn(19)
		}
		if !contains(selected, candidate) {
			selected = append(selected, candidate)
		}
	}

	return selected
}

func (s *GeneratorService) weightedPick(weights map[int]int, max int) int {
	total := 0
	for value := 1; value <= max; value++ {
		weight := weights[value]
		if weight < 1 {
			weight = 1
		}
		total += weight
	}

	point := s.rng.Intn(total) + 1
	running := 0
	for value := 1; value <= max; value++ {
		weight := weights[value]
		if weight < 1 {
			weight = 1
		}
		running += weight
		if running >= point {
			return value
		}
	}
	return max
}

func scoreCombination(numbers []int, stars []int, strategy string, frequencies FrequenciesResponse, delays DelaysResponse) int {
	score := 70
	if strategy == "random" {
		score = 55
	}

	odd := 0
	low := 0
	sum := 0
	for _, number := range numbers {
		sum += number
		if number%2 == 1 {
			odd++
		}
		if number <= 25 {
			low++
		}
	}

	if odd == 2 || odd == 3 {
		score += 8
	}
	if low == 2 || low == 3 {
		score += 8
	}
	if sum >= 100 && sum <= 170 {
		score += 6
	}
	score -= humanPatternPenalty(numbers)
	score += int(math.Min(10, float64(averageDelay(numbers, delays.Numbers))))

	if score < 1 {
		return 1
	}
	if score > 100 {
		return 100
	}
	return score
}

func humanPatternPenalty(numbers []int) int {
	penalty := 0
	under32 := 0
	even := 0
	consecutive := 0

	for index, number := range numbers {
		if number <= 31 {
			under32++
		}
		if number%2 == 0 {
			even++
		}
		if index > 0 && number == numbers[index-1]+1 {
			consecutive++
		}
	}

	if under32 >= 4 {
		penalty += 12
	}
	if even == 0 || even == len(numbers) {
		penalty += 15
	}
	if consecutive >= 2 {
		penalty += 16
	}
	if numbers[0] == 1 && numbers[1] == 2 && numbers[2] == 3 {
		penalty += 25
	}

	return penalty
}

func averageDelay(numbers []int, delays []DelayStat) int {
	delayByValue := map[int]int{}
	for _, delay := range delays {
		delayByValue[delay.Value] = delay.Delay
	}
	total := 0
	for _, number := range numbers {
		total += delayByValue[number]
	}
	return total / len(numbers)
}

func maxFrequency(values []FrequencyStat) int {
	max := 0
	for _, value := range values {
		if value.Count > max {
			max = value.Count
		}
	}
	return max
}

func explanationFor(strategy string) string {
	switch strategy {
	case "hot":
		return "Prioriza números y estrellas con mayor frecuencia histórica."
	case "cold":
		return "Prioriza números y estrellas con menor frecuencia histórica."
	case "delayed":
		return "Prioriza valores que llevan más sorteos sin aparecer."
	case "anti_human":
		return "Evita patrones típicos como fechas, secuencias y exceso de números bajos."
	case "random":
		return "Combinación aleatoria válida sin ponderación estadística."
	default:
		return "Combinación equilibrada con mezcla de frecuencia, retraso, pares/impares y bajos/altos."
	}
}
