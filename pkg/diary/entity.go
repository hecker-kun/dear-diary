package diary

type Entry struct {
	ID        uint32 `json:"id"`
	Content   string `json:"content"`
	CreatedAt string `json:"created_at"`
	Author 	  string `json:"author"`
}
