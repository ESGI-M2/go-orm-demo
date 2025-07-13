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

	orm := builder.NewSimpleORM().
		WithConfigBuilder(cfg)

	if err := orm.Connect(); err != nil {
		log.Fatalf("connect: %v", err)
	}

	shared.Pretty("status", "connected to database")

	if err := orm.Close(); err != nil {
		log.Fatalf("close: %v", err)
	}

	shared.Pretty("status", "connection closed cleanly")
}
