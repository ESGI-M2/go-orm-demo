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

	// prepare bulk users
	var users []interface{}
	for i := 0; i < 5; i++ {
		n := fmt.Sprintf("Bulk_%d", time.Now().UnixNano()%1e6+int64(i))
		users = append(users, &shared.User{Name: n, Email: fmt.Sprintf("%s@example.com", n), CreatedAt: time.Now()})
	}

	if err := repo.BatchCreate(users); err != nil {
		log.Fatalf("batch create: %v", err)
	}
	shared.Pretty("batch created users", users)

	// chunk through all users in batches of 4
	err := repo.Chunk(4, func(chunk []interface{}) error {
		shared.Pretty("processing chunk", chunk)
		return nil
	})
	if err != nil {
		log.Fatalf("chunk iterate: %v", err)
	}
}
