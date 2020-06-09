package types

import "time"

type Cfg struct {
	XML XML `yaml:"XML"`
}

type XML struct {
	RefPath string `yaml:"RefPath"`
	AmpPath string `yaml:"AmpPath"`
}

type Substance struct {
	Code string
	Name string
}

type PharmaceuticalForm struct {
	Code string
	Name string
}

type ActualMedicinalProduct struct {
	OfficialName     string
	From             time.Time
	AmpComponent AmpComponent
	VmpCode      int
	Code         string
	Ampp Ampp
}

type AmpComponent struct {
	RealActualIngredient   RealActualIngredient
	PharmaceuticalFormCode string
	VmpComponentCode       int
}

type Ampp struct {
	PosologyNotes []string
}

type RealActualIngredient struct {
	Type                  string
	Strength              Quantity
	SubstanceCode         string
	From                  time.Time
}

type Quantity struct {
	Value float64
	Unit  string
}
