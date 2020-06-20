package posology

import (
	"fmt"
	"github.com/bertpersyn/posology-graphql/internal/sam/types"
	"github.com/martinlindhe/unit"
	"github.com/mitchellh/mapstructure"
	"strings"
)

type Dosage struct {
	Strength *Strength
	Period   *Period
}

type Strength struct {
	Value float64
	Unit  string
}

type Period struct {
	Value float64
	Cron  string
}

type PK interface {
	GetID() string
}

type ParacetamolDosageIdentifier struct {
	Risk                       bool
	Adult                      bool
	MaxWeight                  float64
	RouteOfAdministrationCodes map[string]bool
}

type ParacetamolCalcRequest struct {
	Risk                      bool
	Adult                     bool
	Weight                    float64
	RouteOfAdministrationCode string
}

func (p *ParacetamolDosageIdentifier) Fits(req *ParacetamolCalcRequest) bool {
	_, e := p.RouteOfAdministrationCodes[req.RouteOfAdministrationCode]
	if req.Risk == p.Risk && req.Adult == p.Adult && (req.Weight < p.MaxWeight || p.MaxWeight == -1) && e {
		return true
	}
	return false
}

func (p *ParacetamolDosageIdentifier) GetID() string {
	return fmt.Sprintf("paracetamol_%v_%v_%v_%v", p.Risk, p.Adult, p.MaxWeight, p.RouteOfAdministrationCodes)
}

type Posology struct {
	data map[string]map[PK]*Dosage
}

func New() *Posology {
	p := &Posology{
		data: map[string]map[PK]*Dosage{},
	}
	p.initParacetamol()
	return p
}

var weightMap = map[string]unit.Mass{
	"mg": unit.Milligram,
	"dg": unit.Decigram,
	"cg": unit.Centigram,
	"g":  unit.Gram,
	"kg": unit.Kilogram,
}

var calcMap = map[unit.Mass]func(grams unit.Mass) float64{
	unit.Milligram: func(g unit.Mass) float64 { return g.Milligrams() },
	unit.Decigram:  func(g unit.Mass) float64 { return g.Decigrams() },
	unit.Centigram: func(g unit.Mass) float64 { return g.Centigrams() },
	unit.Gram:      func(g unit.Mass) float64 { return g.Grams() },
	unit.Kilogram:  func(g unit.Mass) float64 { return g.Kilograms() },
}

//todo: put in paracetamol pkg
func (p *Posology) initParacetamol() {
	peros := map[string]bool{"57": true, "66": true}
	parenteraal := map[string]bool{"49": true}
	pk1 := &ParacetamolDosageIdentifier{
		Risk:                       false,
		Adult:                      true,
		MaxWeight:                  50,
		RouteOfAdministrationCodes: peros, //oraal
	}
	pk2 := &ParacetamolDosageIdentifier{
		Risk:                       false,
		Adult:                      false,
		MaxWeight:                  50,
		RouteOfAdministrationCodes: peros, //oraal
	}
	pk3 := &ParacetamolDosageIdentifier{
		Risk:                       false,
		Adult:                      true,
		MaxWeight:                  -1,
		RouteOfAdministrationCodes: peros, //oraal
	}
	pk4 := &ParacetamolDosageIdentifier{
		Risk:                       false,
		Adult:                      true,
		MaxWeight:                  -1,
		RouteOfAdministrationCodes: parenteraal, //intraveneus
	}

	d1 := &Dosage{
		Strength: &Strength{
			Value: 15,
			Unit:  "mg/kg",
		},
		Period: &Period{
			Value: 4,
			Cron:  "0 */6 * * *", //4 times a day = “At minute 0 past every 6th hour.”
		},
	}

	d3 := &Dosage{
		Strength: &Strength{
			Value: 1,
			Unit:  "g",
		},
		Period: &Period{
			Value: 4,
			Cron:  "0 */6 * * *", //4 times a day = “At minute 0 past every 6th hour.”
		},
	}

	paracetamolPosologyMap := map[PK]*Dosage{
		pk1: d1,
		pk2: d1,
		pk3: d3,
		pk4: d3,
	}

	p.data["4"] = paracetamolPosologyMap
}

func (p *Posology) getParacetamolDosage(req *ParacetamolCalcRequest) (bool, *Dosage) {
	if data, exists := p.data["4"]; exists {
		for x, y := range data {
			x2 := x.(*ParacetamolDosageIdentifier)
			if x2.Fits(req) {
				return true, y
			}
		}
	}

	return false, nil
}

func (p *Posology) GetDosage(amp *types.ActualMedicinalProduct, calcRequest map[string]interface{}) (bool, *Dosage, error) {
	sc := amp.AmpComponent.RealActualIngredient.SubstanceCode
	switch sc {
	case "4":
		var paracetamolCalcRequest ParacetamolCalcRequest
		config := &mapstructure.DecoderConfig{
			WeaklyTypedInput: true,
			Result:           &paracetamolCalcRequest,
		}
		decoder, err := mapstructure.NewDecoder(config)
		if err != nil {
			return false, nil, err
		}
		err = decoder.Decode(calcRequest)
		if err != nil {
			return false, nil, err
		}
		paracetamolCalcRequest.RouteOfAdministrationCode = amp.AmpComponent.RouteOfAdministrationCode
		found, dosage := p.getParacetamolDosage(&paracetamolCalcRequest)

		if !found {
			return found, dosage, nil
		}
		nwDosage := &Dosage{
			Strength: p.calcDosageStrength(paracetamolCalcRequest.Weight,
				&Strength{
					Value: amp.AmpComponent.RealActualIngredient.Strength.Value,
					Unit:  amp.AmpComponent.RealActualIngredient.Strength.Unit,
				}, dosage.Strength),
			Period: dosage.Period,
		}
		return found, nwDosage, nil
	}

	return false, nil, fmt.Errorf("substance %v has no posology atm", sc)
}

//todo: inject pos calculations using method
//todo: question -> is Dosage the correct object to be returned?

func (p *Posology) calcDosageStrength(weightKg float64, medicineStrength *Strength, posStrength *Strength) *Strength {
	medNom, medDenom := p.getNomDenom(medicineStrength)
	posNom, posDenom := p.getNomDenom(posStrength)

	posUnit := weightMap[posNom]
	medUnit := weightMap[medNom]

	posValueGram := posStrength.Value * posUnit.Grams()
	medValueGram := medicineStrength.Value * medUnit.Grams()
	dosage := 0.0
	if posDenom != "" {
		dosage = (posValueGram * (unit.Mass(weightKg).Grams() / weightMap[posDenom].Grams())) / medValueGram
	} else {
		dosage = posValueGram / medValueGram
	}
	strUnit := medDenom
	if strUnit == "" {
		strUnit = "pieces"
	}
	return &Strength{
		Value: dosage,
		Unit:  strUnit,
	}
}

func (p *Posology) getNomDenom(s *Strength) (string, string) {
	sSplit := strings.Split(s.Unit, "/")
	nom := sSplit[0]
	denom := ""
	if len(sSplit) == 2 {
		denom = sSplit[1]
	}
	return nom, denom
}
