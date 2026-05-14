package models

import "time"

type Draw struct {
	ID         uint      `json:"id"`
	DrawNumber int       `json:"drawNumber"`
	DrawDate   time.Time `json:"drawDate"`
	N1         int       `json:"n1"`
	N2         int       `json:"n2"`
	N3         int       `json:"n3"`
	N4         int       `json:"n4"`
	N5         int       `json:"n5"`
	Star1      int       `json:"star1"`
	Star2      int       `json:"star2"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
}

func (d Draw) Numbers() []int {
	return []int{d.N1, d.N2, d.N3, d.N4, d.N5}
}

func (d Draw) Stars() []int {
	return []int{d.Star1, d.Star2}
}
