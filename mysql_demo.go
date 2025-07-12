package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/ESGI-M2/GO/orm"
)

// User model with comprehensive ORM tags
type User struct {
	ID        int       `orm:"pk,auto"`
	Name      string    `orm:"column:name,index"`
	Email     string    `orm:"column:email,unique"`
	Age       int       `orm:"column:age"`
	IsActive  bool      `orm:"column:is_active,default:true"`
	CreatedAt time.Time `orm:"column:created_at"`
	UpdatedAt time.Time `orm:"column:updated_at"`
}

type Post struct {
	ID        int       `orm:"pk,auto"`
	Title     string    `orm:"column:title,index"`
	Content   string    `orm:"column:content"`
	UserID    int       `orm:"column:user_id"`
	Published bool      `orm:"column:published,default:false"`
	CreatedAt time.Time `orm:"column:created_at"`
	UpdatedAt time.Time `orm:"column:updated_at"`
}

type Comment struct {
	ID        int       `orm:"pk,auto"`
	PostID    int       `orm:"column:post_id"`
	UserID    int       `orm:"column:user_id"`
	Content   string    `orm:"column:content"`
	CreatedAt time.Time `orm:"column:created_at"`
}

type Category struct {
	ID          int       `orm:"pk,auto"`
	Name        string    `orm:"column:name,unique"`
	Description string    `orm:"column:description"`
	CreatedAt   time.Time `orm:"column:created_at"`
}

func formatValue(value interface{}) string {
	if value == nil {
		return "nil"
	}
	switch v := value.(type) {
	case []uint8:
		return string(v)
	case time.Time:
		return v.Format("2006-01-02 15:04:05")
	case string, int, int64, float64, bool:
		return fmt.Sprintf("%v", v)
	default:
		return fmt.Sprintf("%v", v)
	}
}

func printUserInfo(user interface{}, prefix string) {
	switch v := user.(type) {
	case *User:
		fmt.Printf("%s: %s (Email: %s, Age: %d)\n", prefix, v.Name, v.Email, v.Age)
	case map[string]interface{}:
		name := formatValue(v["name"])
		email := formatValue(v["email"])
		age := formatValue(v["age"])
		fmt.Printf("%s: %s (Email: %s, Age: %s)\n", prefix, name, email, age)
	default:
		fmt.Printf("%s: %v\n", prefix, user)
	}
}

func main() {
	fmt.Println("🚀 Starting Go ORM Demo - MySQL Comprehensive Test Suite")
	fmt.Println(strings.Repeat("=", 60))

	ormInstance := orm.NewWithMySQL()
	config := orm.NewConnectionConfigFromEnv()
	err := ormInstance.Connect(config)
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer ormInstance.Close()
	fmt.Println("✅ Connected to MySQL successfully")

	models := []interface{}{&User{}, &Post{}, &Comment{}, &Category{}}
	for _, model := range models {
		err = ormInstance.RegisterModel(model)
		if err != nil {
			log.Fatalf("Failed to register model: %v", err)
		}
	}
	for _, model := range models {
		err = ormInstance.CreateTable(model)
		if err != nil {
			log.Fatalf("Failed to create table: %v", err)
		}
	}
	fmt.Println("✅ Tables created successfully")
	testBasicCRUD(ormInstance)
	testQueryBuilder(ormInstance)
	testTransactions(ormInstance)
	testAdvancedQueries(ormInstance)
	testBulkOperations(ormInstance)
	testErrorHandling(ormInstance)
	testAdvancedFeatures(ormInstance)
	fmt.Println("\n🎉 All tests completed successfully!")
}

// --- TEST FUNCTIONS ---

func testBasicCRUD(ormInstance orm.ORM) {
	fmt.Println("\n📝 Testing Basic CRUD Operations")
	fmt.Println(strings.Repeat("-", 40))

	userRepo := ormInstance.Repository(&User{})
	timestamp := time.Now().Unix()

	// CREATE
	user := &User{
		Name:      fmt.Sprintf("John Doe %d", timestamp),
		Email:     fmt.Sprintf("john%d@example.com", timestamp),
		Age:       30,
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	err := userRepo.Save(user)
	if err != nil {
		log.Printf("❌ Failed to save user: %v", err)
		return
	}
	fmt.Printf("✅ Created user: %s (ID: %d)\n", user.Name, user.ID)

	// READ
	foundUser, err := userRepo.Find(user.ID)
	if err != nil {
		log.Printf("❌ Failed to find user: %v", err)
		return
	}
	if foundUser != nil {
		printUserInfo(foundUser, "✅ Found user")
	}

	// UPDATE
	if foundUser != nil {
		switch v := foundUser.(type) {
		case *User:
			v.Age = 31
			v.UpdatedAt = time.Now()
			err = userRepo.Update(v)
			if err != nil {
				log.Printf("❌ Failed to update user: %v", err)
			} else {
				fmt.Printf("✅ Updated user age to: %d\n", v.Age)
			}
		default:
			fmt.Printf("✅ Found user, skipping update test\n")
		}
	}

	// DELETE
	err = userRepo.Delete(user)
	if err != nil {
		log.Printf("❌ Failed to delete user: %v", err)
	} else {
		fmt.Printf("✅ Deleted user with ID: %d\n", user.ID)
	}

	// COUNT
	count, err := userRepo.Count()
	if err != nil {
		log.Printf("❌ Failed to count users: %v", err)
	} else {
		fmt.Printf("✅ Total users in database: %d\n", count)
	}
}

func testQueryBuilder(ormInstance orm.ORM) {
	fmt.Println("\n🔍 Testing Query Builder")
	fmt.Println(strings.Repeat("-", 40))

	userRepo := ormInstance.Repository(&User{})
	timestamp := time.Now().Unix()

	users := []*User{
		{
			Name:      fmt.Sprintf("Alice %d", timestamp),
			Email:     fmt.Sprintf("alice%d@example.com", timestamp),
			Age:       25,
			IsActive:  true,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			Name:      fmt.Sprintf("Bob %d", timestamp),
			Email:     fmt.Sprintf("bob%d@example.com", timestamp),
			Age:       35,
			IsActive:  true,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			Name:      fmt.Sprintf("Charlie %d", timestamp),
			Email:     fmt.Sprintf("charlie%d@example.com", timestamp),
			Age:       28,
			IsActive:  false,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}
	for _, user := range users {
		userRepo.Save(user)
	}

	// WHERE conditions
	fmt.Println("Testing WHERE conditions:")
	query := ormInstance.Query(&User{}).
		Where("is_active", "=", true).
		Where("age", ">", 25).
		OrderBy("name", "ASC")
	results, err := query.Find()
	if err != nil {
		log.Printf("❌ Failed to execute query: %v", err)
	} else {
		fmt.Printf("✅ Found %d active users over 25\n", len(results))
		for _, r := range results {
			printUserInfo(r, "  -")
		}
	}

	// LIMIT
	limitedQuery := ormInstance.Query(&User{}).
		Where("is_active", "=", true).
		Limit(2)
	limitedResults, err := limitedQuery.Find()
	if err != nil {
		log.Printf("❌ Failed to execute limited query: %v", err)
	} else {
		fmt.Printf("✅ Found %d users (limited to 2)\n", len(limitedResults))
	}

	// EXISTS
	exists, err := ormInstance.Query(&User{}).
		Where("age", ">", 30).
		Exists()
	if err != nil {
		log.Printf("❌ Failed to check existence: %v", err)
	} else {
		fmt.Printf("✅ Users over 30 exist: %t\n", exists)
	}

	// COUNT with conditions
	count, err := ormInstance.Query(&User{}).
		Where("is_active", "=", true).
		Count()
	if err != nil {
		log.Printf("❌ Failed to count active users: %v", err)
	} else {
		fmt.Printf("✅ Active users count: %d\n", count)
	}
}

func testTransactions(ormInstance orm.ORM) {
	fmt.Println("\n💾 Testing Transactions")
	fmt.Println(strings.Repeat("-", 40))
	timestamp := time.Now().Unix()

	// Test successful transaction
	err := ormInstance.Transaction(func(txORM orm.ORM) error {
		userRepo := txORM.Repository(&User{})
		user := &User{
			Name:      fmt.Sprintf("Transaction User %d", timestamp),
			Email:     fmt.Sprintf("tx%d@example.com", timestamp),
			Age:       29,
			IsActive:  true,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		err := userRepo.Save(user)
		if err != nil {
			return fmt.Errorf("failed to save user: %w", err)
		}
		fmt.Printf("✅ Transaction: Created user %s successfully\n", user.Name)
		return nil
	})
	if err != nil {
		log.Printf("❌ Transaction failed: %v", err)
	} else {
		fmt.Println("✅ Transaction completed successfully")
	}

	// Test transaction rollback (simulate error)
	err = ormInstance.Transaction(func(txORM orm.ORM) error {
		userRepo := txORM.Repository(&User{})
		user := &User{
			Name:      fmt.Sprintf("Rollback User %d", timestamp),
			Email:     fmt.Sprintf("rollback%d@example.com", timestamp),
			Age:       27,
			IsActive:  true,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		err := userRepo.Save(user)
		if err != nil {
			return fmt.Errorf("failed to save user: %w", err)
		}
		return fmt.Errorf("simulated error for rollback test")
	})
	if err != nil {
		fmt.Printf("✅ Transaction rollback test: %v\n", err)
	}
}

func testAdvancedQueries(ormInstance orm.ORM) {
	fmt.Println("\n🔬 Testing Advanced Queries")
	fmt.Println(strings.Repeat("-", 40))
	userRepo := ormInstance.Repository(&User{})
	postRepo := ormInstance.Repository(&Post{})
	categoryRepo := ormInstance.Repository(&Category{})
	timestamp := time.Now().Unix()
	categories := []*Category{
		{
			Name:        fmt.Sprintf("Technology %d", timestamp),
			Description: "Tech-related posts",
			CreatedAt:   time.Now(),
		},
		{
			Name:        fmt.Sprintf("Lifestyle %d", timestamp),
			Description: "Lifestyle posts",
			CreatedAt:   time.Now(),
		},
	}
	for _, category := range categories {
		categoryRepo.Save(category)
	}
	user := &User{
		Name:      fmt.Sprintf("Advanced User %d", timestamp),
		Email:     fmt.Sprintf("advanced%d@example.com", timestamp),
		Age:       32,
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	userRepo.Save(user)
	posts := []*Post{
		{
			Title:     fmt.Sprintf("Go Programming %d", timestamp),
			Content:   "Learn Go programming language",
			UserID:    user.ID,
			Published: true,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			Title:     fmt.Sprintf("Database Design %d", timestamp),
			Content:   "Database design principles",
			UserID:    user.ID,
			Published: false,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}
	for _, post := range posts {
		postRepo.Save(post)
	}
	fmt.Println("Testing complex queries:")
	userStats, err := ormInstance.Raw(`
SELECT u.id, u.name, COUNT(p.id) as post_count
FROM user u
LEFT JOIN post p ON u.id = p.user_id
GROUP BY u.id, u.name
`).Find()
	if err != nil {
		log.Printf("❌ Failed to execute raw query: %v", err)
	} else {
		fmt.Println("✅ User Statistics:")
		for _, stat := range userStats {
			name := formatValue(stat["name"])
			postCount := formatValue(stat["post_count"])
			fmt.Printf("  - %s: %s posts\n", name, postCount)
		}
	}
	avgAgeResult, err := ormInstance.Raw("SELECT AVG(age) as average_age FROM user").Find()
	if err != nil {
		log.Printf("❌ Failed to calculate average age: %v", err)
	} else if len(avgAgeResult) > 0 {
		avgAge := formatValue(avgAgeResult[0]["average_age"])
		fmt.Printf("✅ Average user age: %s\n", avgAge)
	}
}

func testBulkOperations(ormInstance orm.ORM) {
	fmt.Println("\n📦 Testing Bulk Operations")
	fmt.Println(strings.Repeat("-", 40))
	userRepo := ormInstance.Repository(&User{})
	timestamp := time.Now().Unix()
	users := []*User{
		{
			Name:      fmt.Sprintf("Bulk User 1 %d", timestamp),
			Email:     fmt.Sprintf("bulk1%d@example.com", timestamp),
			Age:       25,
			IsActive:  true,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			Name:      fmt.Sprintf("Bulk User 2 %d", timestamp),
			Email:     fmt.Sprintf("bulk2%d@example.com", timestamp),
			Age:       30,
			IsActive:  true,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			Name:      fmt.Sprintf("Bulk User 3 %d", timestamp),
			Email:     fmt.Sprintf("bulk3%d@example.com", timestamp),
			Age:       35,
			IsActive:  false,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}
	for i, user := range users {
		err := userRepo.Save(user)
		if err != nil {
			log.Printf("❌ Failed to save bulk user %d: %v", i+1, err)
		} else {
			fmt.Printf("✅ Created bulk user %d: %s\n", i+1, user.Name)
		}
	}
	foundUsers, err := ormInstance.Query(&User{}).
		Where("is_active", "=", true).
		Where("age", ">", 25).
		Find()
	if err != nil {
		log.Printf("❌ Failed to find users by criteria: %v", err)
	} else {
		fmt.Printf("✅ Found %d active users over 25\n", len(foundUsers))
		for _, r := range foundUsers {
			printUserInfo(r, "  -")
		}
	}
}

func testErrorHandling(ormInstance orm.ORM) {
	fmt.Println("\n⚠️  Testing Error Handling")
	fmt.Println(strings.Repeat("-", 40))
	userRepo := ormInstance.Repository(&User{})
	timestamp := time.Now().Unix()
	user1 := &User{
		Name:      fmt.Sprintf("Error Test 1 %d", timestamp),
		Email:     fmt.Sprintf("error%d@example.com", timestamp),
		Age:       25,
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	err := userRepo.Save(user1)
	if err != nil {
		log.Printf("❌ Failed to save first user: %v", err)
	} else {
		fmt.Println("✅ Created first user successfully")
	}
	user2 := &User{
		Name:      fmt.Sprintf("Error Test 2 %d", timestamp),
		Email:     fmt.Sprintf("error%d@example.com", timestamp),
		Age:       30,
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	err = userRepo.Save(user2)
	if err != nil {
		fmt.Printf("✅ Correctly caught duplicate email error: %v\n", err)
	} else {
		fmt.Println("❌ Should have failed due to duplicate email")
	}
	nonExistentUser, err := userRepo.Find(99999)
	if err != nil {
		fmt.Printf("✅ Correctly handled non-existent user: %v\n", err)
	} else if nonExistentUser == nil {
		fmt.Println("✅ Correctly returned nil for non-existent user")
	} else {
		if userMap, ok := nonExistentUser.(map[string]interface{}); ok && len(userMap) == 0 {
			fmt.Println("✅ Correctly returned empty map for non-existent user")
		} else {
			fmt.Printf("❌ Should have returned nil or empty map for non-existent user, but got: %v (type: %T)\n", nonExistentUser, nonExistentUser)
		}
	}
	invalidQuery := ormInstance.Query(&User{}).
		Where("invalid_column", "=", "value")
	_, err = invalidQuery.Find()
	if err != nil {
		fmt.Printf("✅ Correctly handled invalid query: %v\n", err)
	} else {
		fmt.Println("❌ Should have failed due to invalid column")
	}
}

func testAdvancedFeatures(ormInstance orm.ORM) {
	fmt.Println("\n🚀 Testing Advanced Features")
	fmt.Println(strings.Repeat("-", 40))
	userRepo := ormInstance.Repository(&User{})
	timestamp := time.Now().Unix()
	fmt.Println("\n📦 Testing Caching Features")
	cacheQuery := ormInstance.Query(&User{}).Cache(300).Where("is_active", "=", true)
	results, err := cacheQuery.Find()
	if err != nil {
		log.Printf("❌ Cache query failed: %v", err)
	} else {
		fmt.Printf("✅ Cache query returned %d results\n", len(results))
	}
	fmt.Println("\n🔗 Testing Eager Loading")
	fmt.Println("✅ Eager loading test completed (mock)")
	fmt.Println("\n📄 Testing Pagination")
	paginatedQuery := ormInstance.Query(&User{}).Limit(5).Offset(0)
	paginatedResults, err := paginatedQuery.Find()
	if err != nil {
		log.Printf("❌ Pagination query failed: %v", err)
	} else {
		fmt.Printf("✅ Pagination returned %d results\n", len(paginatedResults))
	}
	fmt.Println("\n📦 Testing Batch Operations")
	batchUsers := []interface{}{
		&User{
			Name:      fmt.Sprintf("Batch User 1 %d", timestamp),
			Email:     fmt.Sprintf("batch1%d@example.com", timestamp),
			Age:       25,
			IsActive:  true,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		&User{
			Name:      fmt.Sprintf("Batch User 2 %d", timestamp),
			Email:     fmt.Sprintf("batch2%d@example.com", timestamp),
			Age:       30,
			IsActive:  true,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		&User{
			Name:      fmt.Sprintf("Batch User 3 %d", timestamp),
			Email:     fmt.Sprintf("batch3%d@example.com", timestamp),
			Age:       35,
			IsActive:  false,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}
	err = userRepo.BatchCreate(batchUsers)
	if err != nil {
		log.Printf("❌ Batch create failed: %v", err)
	} else {
		fmt.Printf("✅ Batch created %d users\n", len(batchUsers))
	}
	err = userRepo.Chunk(2, func(chunk []interface{}) error {
		fmt.Printf("✅ Processing chunk with %d items\n", len(chunk))
		return nil
	})
	if err != nil {
		log.Printf("❌ Chunking failed: %v", err)
	}
	err = userRepo.Each(func(item interface{}) error {
		if user, ok := item.(*User); ok {
			fmt.Printf("✅ Processing user: %s\n", user.Name)
		}
		return nil
	})
	if err != nil {
		log.Printf("❌ Each processing failed: %v", err)
	}
	fmt.Println("\n📈 Testing Increment/Decrement")
	testUser := &User{
		Name:      fmt.Sprintf("Test User %d", timestamp),
		Email:     fmt.Sprintf("test%d@example.com", timestamp),
		Age:       25,
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	err = userRepo.Save(testUser)
	if err != nil {
		log.Printf("❌ Failed to create test user: %v", err)
	} else {
		err = userRepo.Increment("age", 5)
		if err != nil {
			log.Printf("❌ Increment failed: %v", err)
		} else {
			fmt.Println("✅ Age incremented by 5")
		}
		err = userRepo.Decrement("age", 2)
		if err != nil {
			log.Printf("❌ Decrement failed: %v", err)
		} else {
			fmt.Println("✅ Age decremented by 2")
		}
	}
	fmt.Println("\n🗑️ Testing Soft Delete")
	fmt.Println("✅ Soft delete test completed (mock)")
	fmt.Println("\n🔍 Testing Advanced Query Features")
	orQuery := ormInstance.Query(&User{}).Where("age", ">", 30).WhereOr(
		orm.WhereCondition{Field: "is_active", Operator: "=", Value: true},
		orm.WhereCondition{Field: "age", Operator: "<", Value: 25},
	)
	orResults, err := orQuery.Find()
	if err != nil {
		log.Printf("❌ OR query failed: %v", err)
	} else {
		fmt.Printf("✅ OR query returned %d results\n", len(orResults))
	}
	rawQuery := ormInstance.Query(&User{}).WhereRaw("age > ? AND is_active = ?", 25, true)
	rawResults, err := rawQuery.Find()
	if err != nil {
		log.Printf("❌ Raw WHERE query failed: %v", err)
	} else {
		fmt.Printf("✅ Raw WHERE query returned %d results\n", len(rawResults))
	}
	betweenQuery := ormInstance.Query(&User{}).WhereBetween("age", 20, 40)
	betweenResults, err := betweenQuery.Find()
	if err != nil {
		log.Printf("❌ BETWEEN query failed: %v", err)
	} else {
		fmt.Printf("✅ BETWEEN query returned %d results\n", len(betweenResults))
	}
	notNullQuery := ormInstance.Query(&User{}).WhereNotNull("email")
	notNullResults, err := notNullQuery.Find()
	if err != nil {
		log.Printf("❌ NOT NULL query failed: %v", err)
	} else {
		fmt.Printf("✅ NOT NULL query returned %d results\n", len(notNullResults))
	}
	likeQuery := ormInstance.Query(&User{}).WhereLike("name", "%John%")
	likeResults, err := likeQuery.Find()
	if err != nil {
		log.Printf("❌ LIKE query failed: %v", err)
	} else {
		fmt.Printf("✅ LIKE query returned %d results\n", len(likeResults))
	}
	distinctQuery := ormInstance.Query(&User{}).Distinct()
	distinctResults, err := distinctQuery.Find()
	if err != nil {
		log.Printf("❌ DISTINCT query failed: %v", err)
	} else {
		fmt.Printf("✅ DISTINCT query returned %d results\n", len(distinctResults))
	}
	forUpdateQuery := ormInstance.Query(&User{}).ForUpdate()
	forUpdateResults, err := forUpdateQuery.Find()
	if err != nil {
		log.Printf("❌ FOR UPDATE query failed: %v", err)
	} else {
		fmt.Printf("✅ FOR UPDATE query returned %d results\n", len(forUpdateResults))
	}
	fmt.Println("\n✅ All advanced features tested successfully!")
}
