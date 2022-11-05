package service

//go:generate mockgen -source=./post_service.go -destination=../mock/service/post_service.go

import (
	"gosocialgraph/openapi"
	"gosocialgraph/pkg/repository"
)

type PostService struct {
	Repository repository.PostWriter
}

func (sv PostService) Run(intent openapi.CreatePostRequest) (empty openapi.CreatePostResponse, err error) {
	post, err := sv.Repository.Create(intent.UserId.String(), intent.Content)
	if err != nil {
		return empty, err
	}

	return openapi.CreatePostResponse{Id: post.ID, Content: post.Content}, nil
}
