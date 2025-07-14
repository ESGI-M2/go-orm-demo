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

	_ = orm.GetORM().DropTable(&shared.User{})
	_ = orm.GetORM().CreateTable(&shared.User{})

	repo := orm.GetORM().Repository(&shared.User{})

	user := &shared.User{Name: "Soft", Email: fmt.Sprintf("soft_%d@example.com", time.Now().UnixNano()), Age: 30, CreatedAt: time.Now()}
	if err := repo.Save(user); err != nil {
		log.Printf("save err: %v", err)
	}
	shared.Pretty("saved user", user)

	_ = repo.SoftDelete(user)
	trashed, _ := repo.FindTrashed()
	shared.Pretty("trashed list", trashed)

	_ = repo.Restore(user)
	restored, _ := repo.Find(user.ID)
	shared.Pretty("restored user", restored)

	_ = repo.ForceDelete(user)
	cnt, _ := repo.Exists(user.ID)
	fmt.Printf("exists after force delete: %v\n", cnt)
}
