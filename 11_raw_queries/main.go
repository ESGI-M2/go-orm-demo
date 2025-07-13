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

func seed(orm interfaces.ORM) {
	repo := orm.Repository(&shared.User{})
	for i := 0; i < 3; i++ {
		name := fmt.Sprintf("Raw_%d", time.Now().UnixNano()%1e6+int64(i))
		_ = repo.Save(&shared.User{Name: name, Email: fmt.Sprintf("%s@example.com", name), Age: 20 + i, CreatedAt: time.Now()})
	}
}

func main() {
	cfg := builder.NewConfigBuilder().
		WithDialect(factory.MySQL).
		WithEnvFile("../shared/env.sample").
		FromEnv()

	ormBuilder := builder.NewSimpleORM().
		WithConfigBuilder(cfg).
		RegisterModels(&shared.User{})

	if err := ormBuilder.Connect(); err != nil {
		log.Fatalf("connect: %v", err)
	}
	defer ormBuilder.Close()

	orm := ormBuilder.GetORM()
	seed(orm)

	// raw select
	rows, err := orm.Raw("SELECT id, name, age FROM user ORDER BY id DESC LIMIT 5").Find()
	if err != nil {
		log.Fatalf("raw select: %v", err)
	}
	shared.Pretty("latest 5 users via raw SQL", rows)

	// aggregate
	avgAgeRes, _ := orm.Raw("SELECT AVG(age) as avg_age FROM user").Find()
	shared.Pretty("average age", avgAgeRes)
}
