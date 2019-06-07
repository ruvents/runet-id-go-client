package models

type User struct {
	ID         uint32
	Email      string
	Phone      string
	FirstName  string
	LastName   string
	FatherName string
	Birthday   string
	Work       UserWork
	Photo      UserPhoto
	Status     EventStatus
	Attributes map[string]string
}

type UserWork struct {
	Company  string
	Position string
}

func (user User) GetAttribute(name string) string {
	return user.Attributes[name]
}

func (user User) AddAttribute(name string, value string) {
	user.Attributes[name] = value
}

func (user User) HasAttribute(name string) bool {
	value, found := user.Attributes[name]
	return found && value != ""
}
