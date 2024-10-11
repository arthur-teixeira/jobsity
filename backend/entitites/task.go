package entitites

type Task struct {
	Id          int    `field:"id" json:"id"`
	Title       string `field:"title" json:"title"`
	IsCompleted bool   `field:"is_completed" json:"isCompleted"`
}
