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
		RegisterModels(&shared.User{}, &shared.Post{})

	if err := orm.Connect(); err != nil {
		log.Fatalf("connect: %v", err)
	}
	defer orm.Close()

	// Recreate tables for a clean run
	_ = orm.GetORM().DropTable(&shared.Post{})
	_ = orm.GetORM().DropTable(&shared.User{})
	_ = orm.GetORM().CreateTable(&shared.User{})
	_ = orm.GetORM().CreateTable(&shared.Post{})

	userRepo := orm.GetORM().Repository(&shared.User{})
	postRepo := orm.GetORM().Repository(&shared.Post{})

	u := &shared.User{Name: "Eager", Email: fmt.Sprintf("eager_%d@example.com", time.Now().UnixNano()), Age: 28, CreatedAt: time.Now()}
	if err := userRepo.Save(u); err != nil {
		log.Printf("save user err: %v", err)
	}

	p1 := &shared.Post{Title: "First", Content: "first", UserID: u.ID, CreatedAt: time.Now()}
	p2 := &shared.Post{Title: "Second", Content: "second", UserID: u.ID, CreatedAt: time.Now()}
	_ = postRepo.Save(p1)
	_ = postRepo.Save(p2)

	// Eager load posts when fetching user
	result, err := userRepo.FindWithRelations(u.ID, "Posts")
	if err != nil {
		log.Printf("FindWithRelations err: %v", err)
	}
	shared.Pretty("user with posts", result)

	// WithCount example – count posts per user
	qb := orm.GetORM().Query(&shared.User{}).WithCount("Posts")
	users, err := qb.Find()
	if err != nil {
		log.Printf("WithCount err: %v", err)
	}
	shared.Pretty("users with posts_count", users)
}
