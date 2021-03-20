package resolvers

import (
	"miikka.xyz/rs/graph/generated"
	"miikka.xyz/rs/store"
)

// Resolver is for all dependencies, like db
type Resolver struct {
	db *store.DB
}

// NewRootResolvers helps inject data to resolver
func NewRootResolvers(db *store.DB) *Resolver {
	r := &Resolver{db: db}
	return r
}

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
