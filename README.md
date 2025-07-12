# Go ORM Demo

This project demonstrates the usage and features of the [ESGI-M2/GO](https://github.com/ESGI-M2/GO) ORM library with a comprehensive test suite in `main.go`.

## Prerequisites

- Go 1.24 or later
- Docker and Docker Compose (for local MySQL database)
- Access to [github.com/ESGI-M2/GO](https://github.com/ESGI-M2/GO) (imported as a Go module)

## Setup

1. **Clone this repository** (or copy the demo folder):
   ```sh
   git clone <this-repo-url>
   cd go-orm-demo
   ```

2. **Start the MySQL database using Docker Compose:**
   ```sh
   docker compose up -d
   ```
   This will start a MySQL server and Adminer (web UI at http://localhost:8080).

3. **Configure environment variables (optional):**
   By default, the demo expects:
   - Host: `localhost`
   - Port: `3306`
   - Database: `orm`
   - Username: `user`
   - Password: `password`
   
   You can edit the `.env` file or update the connection config in `main.go` if needed.

4. **Install Go dependencies:**
   ```sh
   go mod tidy
   ```

5. **Run the test suite:**
   ```sh
   go run main.go
   ```

## What Does the Test Suite Do?

The `main.go` file is a comprehensive test suite that exercises all major features of the ORM library. It is organized into several test functions:

### Test Sections

#### 🔧 **Basic CRUD Operations**
- Create, Read, Update, Delete, and Count users
- Demonstrates unique constraints and default values
- Shows automatic timestamp handling

#### 🔍 **Query Builder**
- Complex queries with `WHERE`, `ORDER BY`, `LIMIT`, `EXISTS`, and `COUNT`
- Shows how to use the fluent query API
- Demonstrates chaining multiple conditions

#### 💾 **Transactions**
- Demonstrates ACID-compliant transactions
- Shows both successful commits and rollbacks on error
- Tests transaction isolation and error handling

#### 📊 **Advanced Queries**
- Raw SQL queries (e.g., user statistics with joins and aggregates)
- Aggregate functions (e.g., average age)
- Complex WHERE conditions with multiple operators

#### 📦 **Bulk Operations**
- Bulk creation and querying of users
- Demonstrates efficient batch processing
- Shows how to handle large datasets

#### ⚠️ **Error Handling**
- Tests unique constraint violations
- Handles not found and invalid query errors
- Demonstrates proper error recovery

#### 🚀 **Advanced Features** (NEW!)
- **Caching**: Query result caching with TTL
- **Eager Loading**: Loading related data efficiently
- **Pagination**: Offset and cursor-based pagination
- **Batch Operations**: Create, update, and delete multiple records
- **Chunking**: Process large datasets in chunks
- **Each Processing**: Iterate over records with callbacks
- **Increment/Decrement**: Atomic field updates
- **Soft Deletes**: Mark records as deleted without removing them
- **Advanced Query Features**:
  - OR conditions with multiple criteria
  - Raw WHERE conditions with custom SQL
  - BETWEEN conditions for range queries
  - NULL/NOT NULL conditions
  - LIKE conditions for pattern matching
  - DISTINCT queries for unique results
  - FOR UPDATE locks for concurrent access

### Advanced Features in Detail

#### 📦 **Caching**
```go
// Cache query results for 300 seconds
cacheQuery := orm.Query(&User{}).Cache(300).Where("is_active", "=", true)
results, err := cacheQuery.Find()
```

#### 🔗 **Eager Loading**
```go
// Load users with their related posts
users, err := repo.FindAllWithRelations("posts")
```

#### 📄 **Pagination**
```go
// Offset pagination
paginatedQuery := orm.Query(&User{}).Limit(5).Offset(0)
results, err := paginatedQuery.Find()
```

#### 📦 **Batch Operations**
```go
// Create multiple users at once
users := []interface{}{user1, user2, user3}
err := repo.BatchCreate(users)
```

#### 🔄 **Chunking**
```go
// Process users in chunks of 100
err := repo.Chunk(100, func(chunk []interface{}) error {
    // Process each chunk
    return nil
})
```

#### 📈 **Increment/Decrement**
```go
// Atomic field updates
err := repo.Increment("age", 5)  // Add 5 to age
err := repo.Decrement("age", 2)  // Subtract 2 from age
```

#### 🗑️ **Soft Deletes**
```go
// Mark as deleted (keeps record in database)
err := repo.SoftDelete(user)
// Restore soft-deleted record
err := repo.Restore(user)
```

#### 🔍 **Advanced Query Features**
```go
// OR conditions
orQuery := orm.Query(&User{}).Where("age", ">", 30).WhereOr(
    WhereCondition{Field: "is_active", Operator: "=", Value: true},
    WhereCondition{Field: "age", Operator: "<", Value: 25},
)

// Raw WHERE conditions
rawQuery := orm.Query(&User{}).WhereRaw("age > ? AND is_active = ?", 25, true)

// BETWEEN conditions
betweenQuery := orm.Query(&User{}).WhereBetween("age", 20, 40)

// NULL conditions
notNullQuery := orm.Query(&User{}).WhereNotNull("email")

// LIKE conditions
likeQuery := orm.Query(&User{}).WhereLike("name", "%John%")

// DISTINCT query
distinctQuery := orm.Query(&User{}).Distinct()

// FOR UPDATE lock
forUpdateQuery := orm.Query(&User{}).ForUpdate()
```

### Output Example

When you run the demo, you will see output like:
```
🚀 Starting Go ORM Demo - Comprehensive Test Suite
============================================================
✅ Connected to MySQL successfully
✅ Tables created successfully

📝 Testing Basic CRUD Operations
----------------------------------------
✅ Created user: John Doe ...
✅ Found user: John Doe ...
✅ Updated user age to: ...
✅ Deleted user with ID: ...
✅ Total users in database: ...

🚀 Testing Advanced Features
----------------------------------------
📦 Testing Caching Features
✅ Cache query returned 5 results

📄 Testing Pagination
✅ Pagination returned 5 results

📦 Testing Batch Operations
✅ Batch created 3 users

🔄 Testing Chunking
✅ Processing chunk with 2 items

📈 Testing Increment/Decrement
✅ Age incremented by 5
✅ Age decremented by 2

🔍 Testing Advanced Query Features
✅ OR query returned 3 results
✅ Raw WHERE query returned 2 results
✅ BETWEEN query returned 4 results
✅ NOT NULL query returned 5 results
✅ LIKE query returned 1 results
✅ DISTINCT query returned 5 results
✅ FOR UPDATE query returned 5 results

🎉 All tests completed successfully!
```

## Library Features Overview

The ESGI-M2/GO ORM library provides:

### ✅ **Core Features**
- **CRUD Operations**: Create, Read, Update, Delete
- **Query Builder**: Fluent API for complex queries
- **Transactions**: ACID-compliant database transactions
- **Model Registration**: Automatic table creation and schema management
- **Connection Pooling**: Efficient database connection management

### ✅ **Advanced Features**
- **Caching**: Query result caching with configurable TTL
- **Eager Loading**: Load related data efficiently
- **Pagination**: Both offset and cursor-based pagination
- **Batch Operations**: Create, update, and delete multiple records
- **Soft Deletes**: Mark records as deleted without removing them
- **Chunking**: Process large datasets in manageable chunks
- **Increment/Decrement**: Atomic field updates
- **Advanced Queries**: OR conditions, raw SQL, BETWEEN, NULL checks, LIKE, DISTINCT, locks

### ✅ **Database Support**
- **MySQL**: Full support with MySQL-specific optimizations
- **Mock Dialect**: For testing and development
- **Extensible**: Easy to add support for other databases

### ✅ **Developer Experience**
- **Type Safety**: Strongly typed models and queries
- **Error Handling**: Comprehensive error handling and recovery
- **Testing**: Built-in mock dialect for unit testing
- **Documentation**: Comprehensive examples and documentation

## How to Add More Tests
- Add new models or test cases in `main.go`
- Use the ORM's repository and query builder APIs as shown in the examples
- Extend the test suite with your specific use cases

## Troubleshooting
- If you see a panic or error, check that your MySQL server is running and the credentials match
- If you change the database config, update both Docker Compose and the connection config in `main.go`
- For unique constraint errors, the test suite intentionally triggers and catches these to demonstrate error handling
- Make sure you're using the latest version of the library: `v0.0.3-dev`

## Resources
- [ESGI-M2/GO ORM Library](https://github.com/ESGI-M2/GO)
- [Go Documentation](https://golang.org/doc/)
- [Adminer UI](http://localhost:8080) (for inspecting the database)

---

**This demo is designed to be a living test and showcase for the ORM library. Feel free to extend it with your own models and use cases!** 