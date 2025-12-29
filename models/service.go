package models

type Service struct {
	ID       int     `db:"id" json:"id"`
	Name     string  `db:"name" json:"name"`
	Price    float64 `db:"price" json:"price"`
	Duration string  `db:"duration" json:"duration"`
}
