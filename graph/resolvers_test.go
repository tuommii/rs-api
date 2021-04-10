package resolvers_test

import (
	"testing"

	"github.com/99designs/gqlgen/client"
	"github.com/99designs/gqlgen/graphql/handler"
	"miikka.xyz/rs/config"
	"miikka.xyz/rs/graph/generated"
	"miikka.xyz/rs/graph/model"
	"miikka.xyz/rs/graph/resolvers"
	"miikka.xyz/rs/store"
)

func TestResolvers(t *testing.T) {
	config := config.Load("../.env")
	database := store.NewDB(config)
	defer database.Close()
	r := resolvers.NewRootResolvers(database)
	c := client.New(handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: r})))
	var resp struct {
		GetStatsForUser struct {
			Stats []model.StatsForSport
		}
	}
	q := `query {
		getStatsForUser(username: "tuommii") {
			stats {
				sportName
			}
		}
	}
`
	c.MustPost(q, &resp)
	t.Logf("response: %+v", resp)
}
