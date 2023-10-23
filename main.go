package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type Performance struct {
	Type string `json:"type"`
	Audience int `json:"audience"`
}

type PerformanceBehaviour interface {
	Amount() float64
}

type TragedyPerformance struct {
	Performance
}

func (t *TragedyPerformance) Amount() float64 {
	return float64(t.Audience) * 100
}

type ComedyPerformance struct {
	Performance
}

func (c *ComedyPerformance) Amount() float64 {
	return float64(c.Audience) * 150
}

func NewPerformance(base Performance) PerformanceBehaviour {
	switch base.Type {
	case "tragedy":
		return &TragedyPerformance{base}
	case "comedy":
		return &ComedyPerformance{base}
	default:
		panic("unknow show type")
	}
}


type Invoice struct {
	Performances []PerformanceBehaviour `json:"-"`
}

func (i *Invoice) UnmarshalJSON(data []byte) error {
	// First the "data" is serialized into the data type
	type tempType struct {
		Performances []Performance `json:"performances"`
	}
	var temp tempType
	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}

	// Each Performances data is converted into corresponding behaviour
	// and explicityly pushed into Performances slice
	for _, performance := range temp.Performances {
		performanceBehaviour := NewPerformance(performance)
		i.Performances = append(i.Performances, performanceBehaviour)
	}

	return nil
}

func (i *Invoice) TotalAmount() float64 {
	var amount float64
	for _, performance := range i.Performances {
		amount += performance.Amount()
	}
	return amount
}

func main() {
	file, err := os.Open("invoices.json")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	var invoices []Invoice
	if err = json.NewDecoder(file).Decode(&invoices); err!= nil {
		panic(err)
	}

	fmt.Printf("Total amount: ₹%.2f\n", invoices[0].TotalAmount())
}