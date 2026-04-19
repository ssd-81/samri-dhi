package personas

import (
	_ "embed"
	"encoding/json"
	"fd-credit-score/internal/scoring"
)

//go:embed data/personas.json
var personasJSON []byte

type Persona struct {
	ID        string            `json:"id"`
	Name      string            `json:"name"`
	Age       int               `json:"age"`
	Occupation string          `json:"occupation"`
	City      string            `json:"city"`
	Deposits  []scoring.Deposit `json:"deposits"`
}

type Personas struct {
	Personas []Persona `json:"personas"`
}

func GetAll() ([]Persona, error) {
	var p Personas
	err := json.Unmarshal(personasJSON, &p)
	if err != nil {
		return nil, err
	}
	return p.Personas, nil
}

func GetByID(id string) (*Persona, error) {
	all, err := GetAll()
	if err != nil {
		return nil, err
	}
	for _, p := range all {
		if p.ID == id {
			return &p, nil
		}
	}
	return nil, nil
}