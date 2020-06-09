package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"regexp"

	"github.com/bertpersyn/posology-graphql/internal/graphql/graph/generated"
	"github.com/bertpersyn/posology-graphql/internal/graphql/graph/model"
)

func (r *queryResolver) Medicines(ctx context.Context) ([]*model.Medicine, error) {
	return r.Resolver.Medicines, nil
}

func (r *queryResolver) Substances(ctx context.Context) ([]*model.Substance, error) {
	return r.Resolver.Substances, nil
}

func (r *queryResolver) Search(ctx context.Context, filter *model.Filter) ([]*model.Medicine, error) {
	nameRegex := new(regexp.Regexp)
	if filter.Name == nil {
		nameRegex = regexp.MustCompile(`.*`)
	} else {
		nameRegex = regexp.MustCompile(*filter.Name)
	}
	substanceRegex := new(regexp.Regexp)
	if filter.Substance == nil {
		substanceRegex = regexp.MustCompile(`.*`)
	} else {
		substanceRegex = regexp.MustCompile(*filter.Substance)
	}
	result := []*model.Medicine{}
	for _, m := range r.Resolver.Medicines {
		if nameRegex.MatchString(m.Name) && substanceRegex.MatchString(m.Ingredient.Substance.Name) {
			result = append(result, m)
		}
	}

	return result, nil
}

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }
