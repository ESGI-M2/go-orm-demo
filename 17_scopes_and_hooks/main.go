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

	// Recreate table
	_ = orm.GetORM().DropTable(&shared.User{})
	_ = orm.GetORM().CreateTable(&shared.User{})

	// --- hooks setup ---
	meta, _ := orm.GetORM().GetMetadata(&shared.User{})
	meta.Hooks = &interfaces.ModelHooks{
		BeforeCreate: []func(interface{}) error{
			func(entity interface{}) error {
				if u, ok := entity.(*shared.User); ok {
					// Ensure name is title case
					u.Name = fmt.Sprintf("%s%s", string(u.Name[0]-32), u.Name[1:])
				}
				return nil
			},
		},
		AfterCreate: []func(interface{}) error{
			func(entity interface{}) error {
				if u, ok := entity.(*shared.User); ok {
					log.Printf("after create hook -> user id=%d", u.ID)
				}
				return nil
			},
		},
	}

	// --- scopes setup ---
	meta.Scopes = map[string]func(interfaces.QueryBuilder) interfaces.QueryBuilder{
		"adults": func(q interfaces.QueryBuilder) interfaces.QueryBuilder {
			return q.Where("age", ">=", 18)
		},
		"name_like": func(q interfaces.QueryBuilder) interfaces.QueryBuilder {
			return q.WhereLike("name", "%E%")
		},
	}

	repo := orm.GetORM().Repository(&shared.User{})

	// Seed data
	users := []interface{}{
		&shared.User{Name: "alice", Email: "alice@example.com", Age: 17, CreatedAt: time.Now()},
		&shared.User{Name: "eve", Email: "eve@example.com", Age: 25, CreatedAt: time.Now()},
		&shared.User{Name: "bob", Email: "bob@example.com", Age: 30, CreatedAt: time.Now()},
	}
	_ = repo.BatchCreate(users)

	// Use scopes
	adults := repo.Scope("adults")
	adultList, _ := adults.FindAll()
	shared.Pretty("adults scope", adultList)

	adultsNamedE := repo.Scope("adults").Scope("name_like")
	filterList, _ := adultsNamedE.FindAll()
	shared.Pretty("adults name like E", filterList)
}
