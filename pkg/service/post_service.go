package service

//go:generate mockgen -source=./post_service.go -destination=../mock/service/post_service.go

import (
	"context"
	"gosocialgraph/pkg/entity"
	"gosocialgraph/pkg/repository"

	"github.com/google/uuid"
)

type PostService struct {
	Repository repository.PostWriter
}

func (sv PostService) Run(ctx context.Context, userID uuid.UUID, content string) (post entity.Post, err error) {
	post, err = sv.Repository.Create(ctx, userID.String(), content)
	if err != nil {
		return post, err
	}

	return post, nil
}
