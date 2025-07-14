package main

import (
	"fmt"
	"log"

	"go-orm-demo/shared"

	"github.com/ESGI-M2/GO/orm/builder"
	"github.com/ESGI-M2/GO/orm/factory"
)

func main() {
	cfg := builder.NewConfigBuilder().
		WithDialect(factory.MySQL).
		WithEnvFile("../shared/env.sample").
		FromEnv()

	orm := builder.NewSimpleORM().
		WithConfigBuilder(cfg).
		RegisterModels(&shared.User{})

	if err := orm.Connect(); err != nil {
		log.Fatalf("connect: %v", err)
	}
	defer orm.Close()

	repo := orm.GetORM().Repository(&shared.User{})
	shared.SeedBulkUsers(repo, 10)

	// Cached query (TTL 60 seconds)
	cachedQB := orm.GetORM().Query(&shared.User{}).
		Cache(60).
		OrderBy("id", "DESC").
		Limit(5)
	cachedRes, err := cachedQB.Find()
	if err != nil {
		log.Fatalf("cache query: %v", err)
	}
	shared.Pretty("cached users (limit 5)", cachedRes)

	// Pagination using limit/offset
	perPage := 3
	for page := 0; page < 3; page++ {
		offset := page * perPage
		res, err := orm.GetORM().Query(&shared.User{}).
			OrderBy("id", "ASC").
			Limit(perPage).
			Offset(offset).
			Find()
		if err != nil {
			log.Fatalf("paginate: %v", err)
		}
		shared.Pretty(fmt.Sprintf("page %d", page+1), res)
	}
}
