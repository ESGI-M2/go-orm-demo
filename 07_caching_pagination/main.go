package main

import (
	"fmt"
	"log"
	"time"

	"go-orm-demo/shared"

	ormcore "github.com/ESGI-M2/GO/orm"
	"github.com/ESGI-M2/GO/orm/builder"
	"github.com/ESGI-M2/GO/orm/factory"
)

func seed(repo ormcore.Repository) {
	// Insert 10 users if less than 10 exist
	count, _ := repo.Count()
	if count >= 10 {
		return
	}
	for i := 0; i < 10; i++ {
		name := fmt.Sprintf("User_%d", (time.Now().UnixNano()%1e6)+int64(i))
		u := &shared.User{Name: name, Email: fmt.Sprintf("%s@example.com", name), CreatedAt: time.Now()}
		_ = repo.Save(u)
	}
}

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
	seed(repo)

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
