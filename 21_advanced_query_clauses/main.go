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

	// Seed some rows
	_ = repo.Save(&shared.User{Name: "A", Email: fmt.Sprintf("a_%d@example.com", time.Now().UnixNano()), Age: 10, CreatedAt: time.Now()})
	_ = repo.Save(&shared.User{Name: "B", Email: fmt.Sprintf("b_%d@example.com", time.Now().UnixNano()), Age: 20, CreatedAt: time.Now()})

	// Build two queries and union them
	q1 := orm.GetORM().Query(&shared.User{}).Select("id", "name").Where("age", "<", 15)
	q2 := orm.GetORM().Query(&shared.User{}).Select("id", "name").Where("age", ">", 15)

	union := q1.UnionAll(q2)
	results, _ := union.Find()
	shared.Pretty("union all", results)

	// ForUpdate lock example (no real tx here)
	lockQ := orm.GetORM().Query(&shared.User{}).Where("id", "=", 1).ForUpdate()
	lockRows, _ := lockQ.Find()
	shared.Pretty("for update", lockRows)
}
