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

	// recreate table to include Age column if it was added later
	_ = orm.GetORM().DropTable(&shared.User{})
	_ = orm.GetORM().CreateTable(&shared.User{})

	repo := orm.GetORM().Repository(&shared.User{})

	// create test user
	u := &shared.User{Name: "IncDec", Email: fmt.Sprintf("incdec_%d@example.com", time.Now().UnixNano()), Age: 20, CreatedAt: time.Now()}
	_ = repo.Save(u)

	first, _ := repo.Find(u.ID)
	shared.Pretty("initial user", first)

	_ = repo.Increment("age", 5)
	afterInc, _ := repo.Find(u.ID)
	shared.Pretty("after +5", afterInc)

	_ = repo.Decrement("age", 2)
	afterDec, _ := repo.Find(u.ID)
	shared.Pretty("after -2", afterDec)
}
