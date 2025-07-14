package shared

import (
	"fmt"
	"time"
)

func SeedBasicUsers(repo interface{ Save(interface{}) error }) {
	names := []string{"Alice", "Bob", "Charlie"}
	for _, n := range names {
		u := &User{Name: n, Email: fmt.Sprintf("%s_%d@example.com", n, time.Now().Unix()), CreatedAt: time.Now()}
		_ = repo.Save(u)
	}
}

func SeedAdvancedUsers(repo interface{ Save(interface{}) error }) {
	entries := []struct {
		Name string
		Age  int
	}{
		{"Anna", 22}, {"Brian", 30}, {"Clara", 27}, {"Derek", 19}, {"Eve", 25},
	}
	for _, e := range entries {
		u := &User{Name: e.Name, Email: fmt.Sprintf("%s_%d@example.com", e.Name, time.Now().UnixNano()), Age: e.Age, CreatedAt: time.Now()}
		_ = repo.Save(u)
	}
}

func SeedBulkUsers(repo interface{ Save(interface{}) error }, n int) {
	for i := 0; i < n; i++ {
		name := fmt.Sprintf("User_%d", (time.Now().UnixNano()%1e6)+int64(i))
		u := &User{Name: name, Email: fmt.Sprintf("%s@example.com", name), CreatedAt: time.Now()}
		_ = repo.Save(u)
	}
}

func SeedPosts(repo interface{ Save(interface{}) error }, userID int) {
	posts := []Post{
		{Title: "First", Content: "first", UserID: userID, CreatedAt: time.Now()},
		{Title: "Second", Content: "second", UserID: userID, CreatedAt: time.Now()},
	}
	for _, p := range posts {
		p := p
		_ = repo.Save(&p)
	}
}
