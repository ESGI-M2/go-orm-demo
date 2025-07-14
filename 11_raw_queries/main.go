package main

import (
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

	ormBuilder := builder.NewSimpleORM().
		WithConfigBuilder(cfg).
		RegisterModels(&shared.User{})

	if err := ormBuilder.Connect(); err != nil {
		log.Fatalf("connect: %v", err)
	}
	defer ormBuilder.Close()

	orm := ormBuilder.GetORM()
	repo := orm.Repository(&shared.User{})
	shared.SeedBulkUsers(repo, 3)

	// raw select
	rows, err := orm.Raw("SELECT id, name, age FROM user ORDER BY id DESC LIMIT 5").Find()
	if err != nil {
		log.Fatalf("raw select: %v", err)
	}
	shared.Pretty("latest 5 users via raw SQL", rows)

	// aggregate
	avgAgeRes, _ := orm.Raw("SELECT AVG(age) as avg_age FROM user").Find()
	shared.Pretty("average age", avgAgeRes)
}
