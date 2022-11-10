package rest

import (
	"gosocialgraph/openapi"
	"gosocialgraph/pkg/entity"
)

func MapUserPostsToOpenAPI(posts []entity.UserPost) (openapiResponse []openapi.UserPost, err error) {
	for _, post := range posts {
		var parent *openapi.UserPost

		if post.Post.Parent != nil {
			parent = &openapi.UserPost{
				Post: openapi.Post{
					Content:   post.Post.Parent.Post.Content,
					CreatedAt: post.Post.Parent.Post.CreatedAt,
					Id:        post.Post.Parent.Post.ID,
				},
				User: openapi.User{
					Id:        post.Post.Parent.User.ID,
					CreatedAt: post.Post.Parent.User.CreatedAt,
					Username:  post.Post.Parent.User.Username,
				},
			}
		}

		openapiUserPost := openapi.UserPost{
			User: openapi.User{
				Id:        post.User.ID,
				CreatedAt: post.User.CreatedAt,
				Username:  post.User.Username,
			},
			Post: openapi.Post{
				Id:        post.Post.ID,
				Content:   post.Post.Content,
				CreatedAt: post.Post.CreatedAt,
				Parent:    parent,
			},
		}

		openapiResponse = append(openapiResponse, openapiUserPost)
	}

	return openapiResponse, nil
}
