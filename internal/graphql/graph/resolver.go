package graph

import (
	"github.com/bertpersyn/posology-graphql/internal/graphql/graph/model"
	"github.com/bertpersyn/posology-graphql/internal/posology"
	samparser "github.com/bertpersyn/posology-graphql/internal/sam"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.
//go:generate go run github.com/99designs/gqlgen

type Resolver struct {
	Sam                 *samparser.Service
	Posology            *posology.Posology
	Substances          []*model.Substance
	Medicines           []*model.Medicine
	PharmaceuticalForms []*model.PharmaceuticalForm
}

func (r *Resolver) Init() {
	s := make([]*model.Substance, len(r.Sam.GetSubstances()))
	i := 0
	for _, v := range r.Sam.GetSubstances() {
		s[i] = &model.Substance{
			Name: v.Name,
			Code: v.Code,
		}
		i++
	}
	r.Substances = s

	m := make([]*model.Medicine, len(r.Sam.GetActualMedicinalProducts()))
	i = 0
	for _, v := range r.Sam.GetActualMedicinalProducts() {
		//posologyNote := ""
		//if note, found :=  r.Sam.GetVmpPosologyNotes()[v.VmpCode]; found {
		//	posologyNote = note
		//}
		m[i] = &model.Medicine{
			Code: v.Code,
			Name: v.OfficialName,
			//PosologyNote : &posologyNote,
			Ingredient: &model.Ingredient{
				From: v.AmpComponent.RealActualIngredient.From,
				PharmaceuticalForm: &model.PharmaceuticalForm{
					Name: r.Sam.GetPharmaceuticalForms()[v.AmpComponent.PharmaceuticalFormCode].Name,
					Code: v.AmpComponent.PharmaceuticalFormCode,
				},
				Substance: &model.Substance{
					Name: r.Sam.GetSubstances()[v.AmpComponent.RealActualIngredient.SubstanceCode].Name,
					Code: v.AmpComponent.RealActualIngredient.SubstanceCode,
				},
				Strength: &model.Strength{
					Unit:  v.AmpComponent.RealActualIngredient.Strength.Unit,
					Value: v.AmpComponent.RealActualIngredient.Strength.Value,
				},
			},
		}
		i++
	}
	r.Medicines = m
}

func (r *Resolver) parsePFs() {
	s := make([]*model.PharmaceuticalForm, len(r.Sam.GetPharmaceuticalForms()))
	i := 0
	for _, v := range r.Sam.GetPharmaceuticalForms() {
		s[i] = &model.PharmaceuticalForm{
			Name: v.Name,
			Code: v.Code,
		}
		i++
	}
	r.PharmaceuticalForms = s
}

func (r *Resolver) calcReqToMap(calcRequest []*model.CalcArg) map[string]interface{} {
	m := map[string]interface{}{}
	for _, calcRequest := range calcRequest {
		m[calcRequest.Name] = calcRequest.Value
	}
	return m
}
