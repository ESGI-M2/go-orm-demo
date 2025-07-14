package main

import (
	"fmt"
	"log"

	"go-orm-demo/shared"

	ormcore "github.com/ESGI-M2/GO/orm"

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
	shared.SeedAdvancedUsers(repo)

	q := orm.GetORM().Query(&shared.User{}).
		WhereOr(
			ormcore.WhereCondition{Field: "name", Operator: "LIKE", Value: "%a%"},
			ormcore.WhereCondition{Field: "name", Operator: "LIKE", Value: "%e%"},
		).
		WhereBetween("id", 1, 100000).
		WhereNotNull("email").
		Distinct().
		Limit(3)

	res, err := q.Find()
	if err != nil {
		log.Fatalf("advanced query: %v", err)
	}
	fmt.Printf("advanced query results: %v\n", res)

	shared.Pretty("advanced query results", res)

	inCount, _ := orm.GetORM().Query(&shared.User{}).WhereIn("name", []interface{}{"Anna", "Eve"}).Count()
	fmt.Printf("IN count (Anna, Eve): %d\n", inCount)

	rawRes, _ := orm.GetORM().Query(&shared.User{}).WhereRaw("name LIKE ?", "%r%").Find()
	shared.Pretty("raw where users with 'r'", rawRes)
}
