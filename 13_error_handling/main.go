package main

import (
	"errors"
	"fmt"
	"log"
	"time"

	"go-orm-demo/shared"

	"github.com/ESGI-M2/GO/orm/builder"
	"github.com/ESGI-M2/GO/orm/factory"
	mysql "github.com/go-sql-driver/mysql"
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

	// 1. Duplicate email error
	email := fmt.Sprintf("dup_%d@example.com", time.Now().UnixNano())
	u1 := &shared.User{Name: "Dup1", Email: email, CreatedAt: time.Now()}
	u2 := &shared.User{Name: "Dup2", Email: email, CreatedAt: time.Now()}

	_ = repo.Save(u1)
	err := repo.Save(u2)
	if err != nil {
		fmt.Printf("duplicate save error captured: %v\n", err)
	}

	// 2. Find non-existent user
	notFound, err := repo.Find(999999)
	if err != nil {
		fmt.Printf("find non-existent error: %v\n", err)
	} else if notFound == nil {
		fmt.Println("find non-existent returned nil (as expected)")
	}

	// 3. Invalid column query
	_, err = orm.GetORM().Query(&shared.User{}).Where("invalid_col", "=", 1).Find()
	if err != nil {
		fmt.Printf("invalid column error: %v\n", err)
	}

	// pretty print current users with dup email to confirm only one exists
	rows, _ := orm.GetORM().Query(&shared.User{}).Where("email", "=", email).Find()
	shared.Pretty("rows with duplicate email", rows)

	// human-friendly summary
	if err != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) {
			fmt.Println("MySQL error code:", mysqlErr.Number)
		}
	}
}
