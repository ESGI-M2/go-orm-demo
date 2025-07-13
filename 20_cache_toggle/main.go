package main

import (
	"fmt"
	"log"
	"time"

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

	// Insert a row (skip error handling for brevity)
	_ = repo.Save(&shared.User{Name: "Cached", Email: fmt.Sprintf("cache_%d@example.com", time.Now().UnixNano()), Age: 20, CreatedAt: time.Now()})

	// Query with cache enabled (TTL 60s)
	qb := orm.GetORM().Query(&shared.User{}).Cache(60)
	first, _ := qb.Find()
	shared.Pretty("first fetch (cached)", first)

	// Now disable cache explicitly
	qb2 := orm.GetORM().Query(&shared.User{}).WithoutCache()
	second, _ := qb2.Find()
	shared.Pretty("second fetch (no cache)", second)
}
