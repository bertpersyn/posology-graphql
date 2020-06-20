package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"strings"

	"github.com/bertpersyn/posology-graphql/internal/graphql/graph/generated"
	"github.com/bertpersyn/posology-graphql/internal/graphql/graph/model"
)

func (r *queryResolver) Medicines(ctx context.Context) ([]*model.Medicine, error) {
	return r.Resolver.Medicines, nil
}

func (r *queryResolver) Substances(ctx context.Context) ([]*model.Substance, error) {
	return r.Resolver.Substances, nil
}

func (r *queryResolver) Search(ctx context.Context, searchTerm string) ([]*model.Medicine, error) {
	result := []*model.Medicine{}
	for _, m := range r.Resolver.Medicines {
		if strings.HasPrefix(strings.ToLower(m.Name), strings.ToLower(searchTerm)) || strings.HasPrefix(strings.ToLower(m.Ingredient.Substance.Name), strings.ToLower(searchTerm)) {
			result = append(result, m)
		}
	}

	return result, nil
}

func (r *queryResolver) Posology(ctx context.Context, medicineCode string, calcRequest []*model.CalcArg) (*model.Posology, error) {
	medicine, found := r.Resolver.Sam.GetActualMedicinalProducts()[medicineCode]
	if !found {
		return nil, fmt.Errorf("medicine %v unknown", medicineCode)
	}
	found, dosage, err := r.Resolver.Posology.GetDosage(medicine, r.calcReqToMap(calcRequest))
	if err != nil {
		return nil, err
	}
	if !found {
		return nil, fmt.Errorf("posology for %v not available", medicineCode)
	}

	return &model.Posology{
		Note: "",
		Dosage: &model.Dosage{
			Period: &model.Period{
				Value: dosage.Period.Value,
				Cron:  dosage.Period.Cron,
			},
			Strength: &model.Strength{
				Value: dosage.Strength.Value,
				Unit:  dosage.Strength.Unit,
			},
		},
	}, nil
}

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }
