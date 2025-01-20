package domain

// representation user data that use google as sign-in method
type UserGl struct {
	Id              string
	Name            string
	Email           string
	IsEmailVerified bool

	StreamKey string
}
