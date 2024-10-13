package entitites

type Task struct {
	Id          int    `field:"id" json:"id"`
  UserId      int    `field:"user_id" json:"-"`
	Title       string `field:"title" json:"title"`
	IsCompleted bool   `field:"is_completed" json:"isCompleted"`
}
