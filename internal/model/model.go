package model

import (
	"fmt"
	"strconv"
	"time"
)

type Product struct {
	ID       int
	Name     string
	Category string
	Price    float64
	Data     time.Time
}

func NewProduct(mass []string) (*Product, error) {
	id, err := strconv.Atoi(mass[0])
	if err != nil {
		return nil, fmt.Errorf("error id")
	}
	name := mass[1]
	category := mass[2]
	price, err := strconv.ParseFloat(mass[3], 64)
	if err != nil {
		return nil, fmt.Errorf("error price")
	}
	data, err := time.Parse(time.DateOnly, mass[4])
	if err != nil {
		return nil, fmt.Errorf("error data")
	}
	return &Product{
		ID:       id,
		Name:     name,
		Category: category,
		Price:    price,
		Data:     data,
	}, nil
}
