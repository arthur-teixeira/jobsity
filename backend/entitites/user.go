package entitites

type User struct {
	Id       int    `field:"id"`
	Name     string `field:"name"`
	Email    string `field:"email"`
	Password []byte `field:"password"`
	Salt     []byte `field:"salt"`
}
