package validators

import (
	"errors"
	"sort"
	"time"

	"euromillones/internal/models"
)

type DrawInput struct {
	DrawDate string `json:"drawDate"`
	N1       int    `json:"n1"`
	N2       int    `json:"n2"`
	N3       int    `json:"n3"`
	N4       int    `json:"n4"`
	N5       int    `json:"n5"`
	Star1    int    `json:"star1"`
	Star2    int    `json:"star2"`
}

func NormalizeDraw(input DrawInput) (models.Draw, error) {
	if input.DrawDate == "" {
		return models.Draw{}, errors.New("la fecha del sorteo es obligatoria")
	}

	drawDate, err := time.Parse("2006-01-02", input.DrawDate)
	if err != nil {
		return models.Draw{}, errors.New("la fecha debe tener formato YYYY-MM-DD")
	}

	numbers := []int{input.N1, input.N2, input.N3, input.N4, input.N5}
	stars := []int{input.Star1, input.Star2}

	if err := validateRangeAndUniqueness(numbers, 1, 50, "números principales"); err != nil {
		return models.Draw{}, err
	}
	if err := validateRangeAndUniqueness(stars, 1, 12, "estrellas"); err != nil {
		return models.Draw{}, err
	}

	sort.Ints(numbers)
	sort.Ints(stars)

	return models.Draw{
		DrawDate: drawDate,
		N1:       numbers[0],
		N2:       numbers[1],
		N3:       numbers[2],
		N4:       numbers[3],
		N5:       numbers[4],
		Star1:    stars[0],
		Star2:    stars[1],
	}, nil
}

func validateRangeAndUniqueness(values []int, min int, max int, label string) error {
	seen := map[int]bool{}
	for _, value := range values {
		if value < min || value > max {
			return errors.New(label + " fuera de rango")
		}
		if seen[value] {
			return errors.New(label + " repetidos en el mismo sorteo")
		}
		seen[value] = true
	}

	return nil
}
