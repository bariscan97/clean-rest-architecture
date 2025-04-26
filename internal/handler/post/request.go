package post

type CreatePostReq struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

type UpdatePostReq struct {
	Title   *string `json:"title,omitempty"`
	Content *string `json:"content,omitempty"`
}