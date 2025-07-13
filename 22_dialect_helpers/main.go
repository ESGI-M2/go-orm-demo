package main

import (
	"fmt"
	"log"

	"github.com/ESGI-M2/GO/orm/builder"
	"github.com/ESGI-M2/GO/orm/factory"
)

func main() {
	cfg := builder.NewConfigBuilder().
		WithDialect(factory.MySQL).
		WithEnvFile("../shared/env.sample").
		FromEnv()

	orm := builder.NewSimpleORM().WithConfigBuilder(cfg)
	if err := orm.Connect(); err != nil {
		log.Fatalf("connect: %v", err)
	}
	defer orm.Close()

	d := orm.GetORM().GetDialect()

	// Show helper strings (they would be embedded in SQL normally)
	fmt.Println("Random function:", d.GetRandomFunction())
	fmt.Println("Date function:", d.GetDateFunction())
	fmt.Println("JSON Extract:", d.GetJSONExtract())
	// Example full-text search clause
	ft := d.FullTextSearch("content", "demo")
	fmt.Println("FullTextSearch sample:", ft)
}
