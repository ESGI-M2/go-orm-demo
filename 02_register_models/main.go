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
		FromEnv().
		WithAutoCreateDatabase()

	orm := builder.NewSimpleORM().
		WithConfigBuilder(cfg).
		RegisterModels(&shared.User{}, &shared.Post{})

	if err := orm.Connect(); err != nil {
		log.Fatalf("connect: %v", err)
	}
	defer orm.Close()

	shared.Pretty("status", "models registered and migrated successfully")
}
