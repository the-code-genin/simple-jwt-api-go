package users

type Users interface {
	GetOneById(id int) (*User, error)
	GetOneByEmail(email string) (*User, error)

	Create(user *User) (*User, error)

	EmailTaken(email string) (bool, error)
}
