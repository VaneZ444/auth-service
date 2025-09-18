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
	ID       int64
	Email    string
	Nickname string
	Hash     string
	Role     Role
	Status   Status
}

func NewUser(email, nickname, hashedPassword string, role Role) *User {
	return &User{
		Email:    email,
		Nickname: nickname,
		Hash:     hashedPassword,
		Role:     role,
		Status:   StatusActive,
	}
}
