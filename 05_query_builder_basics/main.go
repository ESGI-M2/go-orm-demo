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
	shared.SeedBasicUsers(repo)

	qb := orm.GetORM().Query(&shared.User{}).
		Where("name", "!=", "").
		OrderBy("id", "DESC").
		Limit(2)

	results, err := qb.Find()
	if err != nil {
		log.Fatalf("query find: %v", err)
	}
	shared.Pretty("latest two users", results)

	count, _ := orm.GetORM().Query(&shared.User{}).Where("name", "LIKE", "%A%").Count()
	fmt.Printf("users with letter A in name: %d\n", count)

	exists, _ := orm.GetORM().Query(&shared.User{}).Where("name", "=", "Bob").Exists()
	fmt.Printf("user Bob exists: %v\n", exists)
}
