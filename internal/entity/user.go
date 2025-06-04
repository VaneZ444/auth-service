package entity

type Role string

const (
	RoleUser  Role = "user"
	RoleAdmin Role = "admin"
)

type Status string

const (
	StatusActive Status = "active"
	StatusBanned Status = "banned"
)

type User struct {
	ID     int64
	Email  string
	Hash   string // Это уже хеш!
	Role   Role
	Status Status
}

func NewUser(email, hashedPassword string, role Role) *User {
	return &User{
		Email:  email,
		Hash:   hashedPassword,
		Role:   role,
		Status: StatusActive,
	}
}
