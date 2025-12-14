package user

type User struct {
	ID          uint
	Name        string
	Email       string
	PhoneNumber string
	Role        Role
	Password    string
}