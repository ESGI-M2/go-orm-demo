package main

import (
	"fmt"
	"log"
	"time"

	"go-orm-demo/shared"

	"github.com/ESGI-M2/GO/orm/builder"
	"github.com/ESGI-M2/GO/orm/factory"
)

func runMySQL() {
	cfg := builder.NewConfigBuilder().
		WithDialect(factory.MySQL).
		WithEnvFile("../shared/env.sample").
		FromEnv()

	orm := builder.NewSimpleORM().
		WithConfigBuilder(cfg).
		RegisterModels(&shared.User{})

	if err := orm.Connect(); err != nil {
		log.Printf("mysql connect failed: %v", err)
		return
	}
	defer orm.Close()

	fmt.Println("dialect:", orm.GetDialectType())
	repo := orm.GetORM().Repository(&shared.User{})
	u := &shared.User{Name: "Multi", Email: fmt.Sprintf("multi_%d@example.com", time.Now().UnixNano()), CreatedAt: time.Now()}
	_ = repo.Save(u)
	shared.Pretty("mysql user", u)
}

func runMock() {
	orm := builder.NewSimpleORM().
		WithDialect(factory.Mock).
		RegisterModels(&shared.User{})

	// empty config OK for mock dialect
	if err := orm.Connect(); err != nil {
		log.Fatalf("mock connect: %v", err)
	}
	defer orm.Close()

	fmt.Println("dialect:", orm.GetDialectType())
	repo := orm.GetORM().Repository(&shared.User{})
	u := &shared.User{Name: "Mocker", Email: "mock@example.com", CreatedAt: time.Now()}
	_ = repo.Save(u)
	rows, _ := repo.Find(u.ID)
	shared.Pretty("mock find", rows)
}

func runPostgres() {
	cfg := builder.NewConfigBuilder().
		WithDialect(factory.Postgres).
		WithEnvFile("../shared/env.sample"). // expects POSTGRES_* vars if running
		FromEnv()

	orm := builder.NewSimpleORM().
		WithConfigBuilder(cfg).
		RegisterModels(&shared.User{})

	if err := orm.Connect(); err != nil {
		log.Printf("postgres connect failed: %v", err)
		return
	}
	defer orm.Close()

	fmt.Println("dialect:", orm.GetDialectType())
	repo := orm.GetORM().Repository(&shared.User{})
	u := &shared.User{Name: "PG", Email: fmt.Sprintf("pg_%d@example.com", time.Now().UnixNano()), CreatedAt: time.Now()}
	_ = repo.Save(u)
	shared.Pretty("postgres user", u)
}

func main() {
	fmt.Println("--- MySQL ---")
	runMySQL()

	fmt.Println("\n--- PostgreSQL ---")
	runPostgres()

	fmt.Println("\n--- Mock Dialect ---")
	runMock()
}
