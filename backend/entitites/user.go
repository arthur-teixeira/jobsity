package entitites

type User struct {
	Id       int    `field:"id"`
	Email    string `field:"email"`
	Password []byte `field:"password"`
	Salt     []byte `field:"salt"`
}
