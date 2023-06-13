package entities

// User struct for authboss
type User struct {
	ID string `bson:"_id"`

	Firstname string
	Lastname  string
	Company   string

	// Auth
	Email    string
	Password string
}
