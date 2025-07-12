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

- **Basic CRUD Operations**
  - Create, Read, Update, Delete, and Count users
  - Demonstrates unique constraints and default values

- **Query Builder**
  - Complex queries with `WHERE`, `ORDER BY`, `LIMIT`, `EXISTS`, and `COUNT`
  - Shows how to use the fluent query API

- **Transactions**
  - Demonstrates ACID-compliant transactions
  - Shows both successful commits and rollbacks on error

- **Advanced Queries**
  - Raw SQL queries (e.g., user statistics with joins and aggregates)
  - Aggregate functions (e.g., average age)

- **Bulk Operations**
  - Bulk creation and querying of users

- **Error Handling**
  - Tests unique constraint violations
  - Handles not found and invalid query errors

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
...
```

## How to Add More Tests
- Add new models or test cases in `main.go`.
- Use the ORM's repository and query builder APIs as shown in the examples.

## Troubleshooting
- If you see a panic or error, check that your MySQL server is running and the credentials match.
- If you change the database config, update both Docker Compose and the connection config in `main.go`.
- For unique constraint errors, the test suite intentionally triggers and catches these to demonstrate error handling.

## Resources
- [ESGI-M2/GO ORM Library](https://github.com/ESGI-M2/GO)
- [Go Documentation](https://golang.org/doc/)
- [Adminer UI](http://localhost:8080) (for inspecting the database)

---

**This demo is designed to be a living test and showcase for the ORM library. Feel free to extend it with your own models and use cases!** 