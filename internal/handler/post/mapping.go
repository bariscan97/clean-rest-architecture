package post

import (
	"github.com/bariscan97/clean-rest-architecture/internal/domains"
)

func CreateReqToDomain(p CreatePostReq) *domains.Post {
	return &domains.Post{
		Title:   p.Title,
		Content: p.Content,
	}
}

func toCreatePostRes(p *domains.Post) CreatePostRes {
	return CreatePostRes{
		ID:       p.ID,
		UserID:   p.UserID,
		ParentID: p.ParentID,
		Title:    p.Title,
		Content:  p.Content,
		UpdateAt: p.UpdateAt,
		CreateAt: p.CreateAt,
	}
}

func ListPostRes(posts []*domains.PostManyToMany) []FetchPostRes {
	var ListPosts []FetchPostRes

	for _, post := range posts {
		ListPosts = append(ListPosts, FetchPostRes{
			ID:       post.ID,
			UserID:   post.UserID,
			ParentID: post.ParentID,
			UserName: post.UserName,
			Title:    post.Title,
			Content:  post.Content,
			UserImg:  post.UserImg,
			UpdateAt: post.UpdateAt,
			CreateAt: post.CreateAt,
		})
	}

	return ListPosts
}
