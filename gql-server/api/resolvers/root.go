package resolvers

import "github.com/aeon/gql-server/api"

type Resolver struct {
}

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }

func (r *Resolver) Mutation() api.MutationResolver { return &mutationResolver{r} }
func (r *Resolver) Query() api.QueryResolver       { return &queryResolver{r} }
