package main

import (
	"fmt"
	"log"
	"time"

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

	// successful transaction
	err := orm.GetORM().Transaction(func(tx ormcore.ORM) error {
		repo := tx.Repository(&shared.User{})
		u := &shared.User{Name: "TxSuccess", Email: fmt.Sprintf("tx_%d@example.com", time.Now().UnixNano()), CreatedAt: time.Now()}
		if err := repo.Save(u); err != nil {
			return err
		}
		shared.Pretty("created in tx", u)
		return nil // commit
	})
	if err != nil {
		log.Fatalf("transaction: %v", err)
	}

	// transaction with rollback
	_ = orm.GetORM().Transaction(func(tx ormcore.ORM) error {
		repo := tx.Repository(&shared.User{})
		u := &shared.User{Name: "TxRollback", Email: fmt.Sprintf("rb_%d@example.com", time.Now().UnixNano()), CreatedAt: time.Now()}
		_ = repo.Save(u)
		return fmt.Errorf("simulated error to rollback")
	})

	// verify rollback user not persisted (should not exist)
	count, _ := orm.GetORM().Query(&shared.User{}).Where("name", "=", "TxRollback").Count()
	fmt.Printf("rollback user count: %d (expect 0)\n", count)
}
