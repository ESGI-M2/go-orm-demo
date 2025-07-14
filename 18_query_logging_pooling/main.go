package main

import (
	"fmt"
	"log"
	"time"

	"go-orm-demo/shared"

	"github.com/ESGI-M2/GO/orm/builder"
	"github.com/ESGI-M2/GO/orm/core/interfaces"
	"github.com/ESGI-M2/GO/orm/factory"
)

func main() {
	// Build config with pool settings
	cfg := builder.NewConfigBuilder().
		WithDialect(factory.MySQL).
		WithConnectionPool(10, 5).
		WithEnvFile("../shared/env.sample").
		FromEnv()

	orm := builder.NewSimpleORM().
		WithConfigBuilder(cfg).
		RegisterModels(&shared.User{})

	if err := orm.Connect(); err != nil {
		log.Fatalf("connect: %v", err)
	}
	defer orm.Close()

	// Enable query log
	orm.GetORM().EnableQueryLog()

	repo := orm.GetORM().Repository(&shared.User{})

	// Simple query to generate logs
	u := &shared.User{Name: "Logger", Email: fmt.Sprintf("logger_%d@example.com", time.Now().UnixNano()), Age: 40, CreatedAt: time.Now()}
	_ = repo.Save(u)
	_, _ = repo.Find(u.ID)

	// Fetch logs
	if logger, ok := orm.GetORM().(interfaces.QueryLogger); ok {
		logs := logger.GetLogs()
		shared.Pretty("query logs", logs)
		logger.ClearLogs()
	}

	// Disable logging and run another query
	orm.GetORM().DisableQueryLog()
	_, _ = repo.Count()

	fmt.Println("query logging disabled – no logs should appear")
}
