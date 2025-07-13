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

	user := &shared.User{
		Name:      "Alice",
		Email:     fmt.Sprintf("alice_%d@example.com", time.Now().Unix()),
		CreatedAt: time.Now(),
	}

	if err := repo.Save(user); err != nil {
		log.Fatalf("save: %v", err)
	}
	shared.Pretty("created user", user)

	found, err := repo.Find(user.ID)
	if err != nil {
		log.Fatalf("find: %v", err)
	}
	shared.Pretty("found user", found)

	user.Name = "Alice Updated"
	if err := repo.Update(user); err != nil {
		log.Fatalf("update: %v", err)
	}
	fmt.Println("updated user name")

	cnt, _ := repo.Count()
	fmt.Printf("user count %d\n", cnt)

	if err := repo.Delete(user); err != nil {
		log.Fatalf("delete: %v", err)
	}
	fmt.Println("user deleted")
}
