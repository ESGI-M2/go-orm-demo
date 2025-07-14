package main

import (
	"fmt"
	"log"

	"go-orm-demo/shared"

	"github.com/ESGI-M2/GO/orm/builder"
	"github.com/ESGI-M2/GO/orm/factory"
)

func main() {
	// Dialect factory
	df := factory.NewDialectFactory()
	shared.Pretty("available dialects", df.GetAvailableDialects())

	// ConfigBuilder workflow
	cfgBuilder := builder.NewConfigBuilder().
		WithDialect(factory.MySQL).
		WithEnvFile("../shared/env.sample").
		FromEnv().
		WithAutoCreateDatabase()

	cfg, dialectType, autoCreate, err := cfgBuilder.Build()
	if err != nil {
		log.Fatalf("config build: %v", err)
	}
	shared.Pretty("built config", map[string]interface{}{"config": cfg, "dialect": dialectType, "autoCreate": autoCreate})

	// QuickSetupFromEnv helper
	quickORM, err := builder.QuickSetupFromEnv("mysql", &shared.User{})
	if err != nil {
		log.Fatalf("quick setup: %v", err)
	}
	fmt.Printf("quick setup dialect: %s\n", quickORM.GetDialectType())
	quickORM.Close()
}
