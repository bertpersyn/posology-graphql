package graph

import (
	"github.com/bertpersyn/posology-graphql/internal/graphql/graph/model"
	"github.com/bertpersyn/posology-graphql/internal/sam/types"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.
//go:generate go run github.com/99designs/gqlgen

type Resolver struct {
	Sam struct {
		PharmaceuticalForms map[string]*types.PharmaceuticalForm
		Substances          map[string]*types.Substance
		Medicines           map[string]*types.ActualMedicinalProduct
	}
	Substances          []*model.Substance
	Medicines           []*model.Medicine
	PharmaceuticalForms []*model.PharmaceuticalForm
}

func (r *Resolver) Init() {
	s := make([]*model.Substance, len(r.Sam.Substances))
	i := 0
	for _, v := range r.Sam.Substances {
		s[i] = &model.Substance{
			Name: v.Name,
			Code: v.Code,
		}
		i++
	}
	r.Substances = s

	m := make([]*model.Medicine, len(r.Sam.Medicines))
	i = 0
	for _, v := range r.Sam.Medicines {
		m[i] = &model.Medicine{
			Code:          v.Code,
			Name:          v.OfficialName,
			PosologyNotes: v.Ampp.PosologyNotes,
			Ingredient: &model.Ingredient{
				From: v.AmpComponent.RealActualIngredient.From,
				PharmaceuticalForm: &model.PharmaceuticalForm{
					Name: r.Sam.PharmaceuticalForms[v.AmpComponent.PharmaceuticalFormCode].Name,
					Code: v.AmpComponent.PharmaceuticalFormCode,
				},
				Substance: &model.Substance{
					Name: r.Sam.Substances[v.AmpComponent.RealActualIngredient.SubstanceCode].Name,
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
	s := make([]*model.PharmaceuticalForm, len(r.Sam.PharmaceuticalForms))
	i := 0
	for _, v := range r.Sam.PharmaceuticalForms {
		s[i] = &model.PharmaceuticalForm{
			Name: v.Name,
			Code: v.Code,
		}
		i++
	}
	r.PharmaceuticalForms = s
}
