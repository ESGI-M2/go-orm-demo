package main

import (
	"fmt"
	"log"
	"time"

	"github.com/ESGI-M2/GO/orm/builder"
	"github.com/ESGI-M2/GO/orm/factory"
)

// Simple User model for factory demo
type FactoryUser struct {
	ID        int       `orm:"pk,auto"`
	Name      string    `orm:"column:name"`
	Email     string    `orm:"column:email,unique"`
	Age       int       `orm:"column:age"`
	CreatedAt time.Time `orm:"column:created_at"`
}

func main() {
	fmt.Println("🏭 Go ORM Factory Pattern Demo")
	fmt.Println("===============================")

	// Method 1: Using SimpleORM Builder Pattern
	fmt.Println("\n📚 Method 1: SimpleORM Builder Pattern")
	mysqlORM := builder.NewSimpleORM().
		WithMySQL().                  // Set dialect to MySQL
		WithEnvConfig().              // Load config from environment
		WithAutoCreateDatabase().     // Auto-create database if needed
		RegisterModel(&FactoryUser{}) // Register our model

	err := mysqlORM.Connect()
	if err != nil {
		log.Printf("❌ Failed to connect to MySQL: %v", err)
	} else {
		fmt.Printf("✅ Connected to MySQL using SimpleORM\n")

		// Basic operation
		repo := mysqlORM.GetORM().Repository(&FactoryUser{})
		user := &FactoryUser{
			Name:      "Factory User",
			Email:     fmt.Sprintf("factory%d@example.com", time.Now().Unix()),
			Age:       25,
			CreatedAt: time.Now(),
		}
		if err := repo.Save(user); err == nil {
			fmt.Printf("✅ Created user: %s (ID: %d)\n", user.Name, user.ID)
		}

		mysqlORM.Close()
	}

	// Method 2: Using QuickSetup
	fmt.Println("\n⚡ Method 2: QuickSetup Helper")
	quickORM, err := builder.QuickSetupFromEnv("postgresql", &FactoryUser{})
	if err != nil {
		log.Printf("❌ QuickSetup failed: %v", err)
	} else {
		fmt.Printf("✅ Connected to PostgreSQL using QuickSetup\n")

		// Basic operation
		repo := quickORM.GetORM().Repository(&FactoryUser{})
		user := &FactoryUser{
			Name:      "Quick User",
			Email:     fmt.Sprintf("quick%d@example.com", time.Now().Unix()),
			Age:       30,
			CreatedAt: time.Now(),
		}
		if err := repo.Save(user); err == nil {
			fmt.Printf("✅ Created user: %s (ID: %d)\n", user.Name, user.ID)
		}

		quickORM.Close()
	}

	// Method 3: Using Factory and ConfigBuilder
	fmt.Println("\n🔧 Method 3: Manual Factory & ConfigBuilder")

	// Create dialect using factory
	dialectFactory := factory.NewDialectFactory()
	_, err = dialectFactory.Create(factory.MySQL)
	if err != nil {
		log.Printf("❌ Failed to create dialect: %v", err)
		return
	}
	fmt.Printf("✅ Created dialect: %s\n", factory.MySQL)

	// Build configuration
	config, _, autoCreate, err := builder.NewConfigBuilder().
		WithDialect(factory.MySQL).
		WithHost("localhost").
		WithPort(3306).
		WithDatabase("orm").
		WithCredentials("user", "password").
		Build()

	if err != nil {
		log.Printf("❌ Failed to build config: %v", err)
		return
	}

	fmt.Printf("✅ Built config: %s@%s:%d/%s (AutoCreate: %v)\n",
		config.Username, config.Host, config.Port, config.Database, autoCreate)

	// Show available dialects
	fmt.Println("\n🗂️  Available Dialects:")
	for _, d := range dialectFactory.GetAvailableDialects() {
		supported := dialectFactory.IsSupported(d)
		fmt.Printf("  - %s (supported: %v)\n", d, supported)
	}

	fmt.Println("\n🎉 Factory Pattern Demo Complete!")
	fmt.Println("\nKey Benefits:")
	fmt.Println("  ✨ Fluent builder pattern for easy configuration")
	fmt.Println("  🔌 Automatic database creation")
	fmt.Println("  📦 Bulk model registration")
	fmt.Println("  ⚡ Quick setup helpers")
	fmt.Println("  🏭 Flexible dialect factory")
	fmt.Println("  🔧 Environment-based configuration")
}
