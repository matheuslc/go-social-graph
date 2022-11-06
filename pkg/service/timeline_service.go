package service

//go:generate mockgen -source=./timeline_service.go -destination=../mock/service/timeline_service.go

import (
	"gosocialgraph/openapi"
	"gosocialgraph/pkg/repository"

	"github.com/google/uuid"
)

// TimelineServive defines the struct holder for the service and its dependencies
type TimelineServive struct {
	Repository repository.TimelineReader
}

// Run executes the use case
func (sv TimelineServive) Run(userID uuid.UUID) (empty openapi.TimelineResponse, err error) {
	response, err := sv.Repository.TimelineFor(userID)
	if err != nil {
		return empty, err
	}

	openapiResponse := []openapi.UserPost{}
	for _, post := range response {
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

	return openapi.TimelineResponse{Posts: &openapiResponse}, nil
}
