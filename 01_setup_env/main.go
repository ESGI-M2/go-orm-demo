package main

import (
	"log"

	"go-orm-demo/shared"

	"github.com/ESGI-M2/GO/orm/builder"
	"github.com/ESGI-M2/GO/orm/factory"
)

func main() {
	cfgBuilder := builder.NewConfigBuilder().
		WithDialect(factory.MySQL).
		WithEnvFile("../shared/env.sample").
		FromEnv().
		WithAutoCreateDatabase()

	ormBuilder := builder.NewSimpleORM().
		WithConfigBuilder(cfgBuilder)

	// Connect to the database
	if err := ormBuilder.Connect(); err != nil {
		log.Fatalf("❌ Connection failed: %v", err)
	}
	defer ormBuilder.Close()

	shared.Pretty("status", "✅ Environment loaded and database connection established!")
}
