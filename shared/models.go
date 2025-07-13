package shared

import "time"

type User struct {
	ID        int        `table:"users" orm:"pk,auto"`
	Name      string     `orm:"column:name"`
	Email     string     `orm:"column:email,unique"`
	Age       int        `orm:"column:age"`
	CreatedAt time.Time  `orm:"column:created_at"`
	DeletedAt *time.Time `orm:"column:deleted_at,soft"`
	Posts     []Post     `orm:"relation:one_to_many,fk:user_id"`
}

type Post struct {
	ID        int       `orm:"pk,auto"`
	Title     string    `orm:"column:title"`
	Content   string    `orm:"column:content"`
	UserID    int       `orm:"column:user_id"`
	CreatedAt time.Time `orm:"column:created_at"`
}
