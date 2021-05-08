package resolvers_test

import (
	"log"
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
	config := config.Load("../.env-testing")
	database := store.NewDB(config)

	err := store.ExecuteSQLFile(database, "../scripts/create_database.sql")
	if err != nil {
		log.Fatal(err)
	}
	err = store.ExecuteSQLFile(database, "../scripts/seed_database.sql")
	if err != nil {
		log.Fatal(err)
	}

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
	//t.Logf("response: %+v", resp)
	for _, obj := range resp.GetStatsForUser.Stats {
		t.Log("Sport name:", obj.SportName)
	}
}
