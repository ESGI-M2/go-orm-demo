package main

import (
	"fmt"
	"log"
	"time"

	"go-orm-demo/shared"

	"github.com/ESGI-M2/GO/orm/builder"
	"github.com/ESGI-M2/GO/orm/factory"
)

// Simple model to illustrate dynamic migration
type Comment struct {
	ID        int       `orm:"pk,auto"`
	Body      string    `orm:"column:body"`
	UserID    int       `orm:"column:user_id"`
	CreatedAt time.Time `orm:"column:created_at"`
}

func main() {
	cfg := builder.NewConfigBuilder().
		WithDialect(factory.MySQL).
		WithEnvFile("../shared/env.sample").
		FromEnv()

	orm := builder.NewSimpleORM().
		WithConfigBuilder(cfg).
		RegisterModels(&shared.User{}, &Comment{})

	if err := orm.Connect(); err != nil {
		log.Fatalf("connect: %v", err)
	}
	defer orm.Close()

	// Drop tables if they exist
	_ = orm.GetORM().DropTable(&Comment{})
	_ = orm.GetORM().DropTable(&shared.User{})

	// Use Migrate() to create all missing tables
	if err := orm.GetORM().Migrate(); err != nil {
		log.Fatalf("migrate: %v", err)
	}
	fmt.Println("Migrate completed (tables created)")

	// Demonstrate dropping/creating a single table manually
	_ = orm.GetORM().DropTable(&Comment{})
	fmt.Println("Comment table dropped")

	_ = orm.GetORM().CreateTable(&Comment{})
	fmt.Println("Comment table recreated via CreateTable()")
}
