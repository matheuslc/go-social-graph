package service

import (
	"fmt"
	"gosocialgraph/pkg/repository"

	"github.com/google/uuid"
)

const ContentMaxLenght = 777

type CreatePostIntent struct {
	UserId  uuid.UUID `json:"user_id"`
	Content string    `json:"content"`
}

type CreatePostResponse struct {
	Id      uuid.UUID `json:"id"`
	Content string    `json:"content"`
}

type PostService struct {
	Repository repository.PostWriter
}

func NewCreatePostIntent(userId uuid.UUID, content string) (CreatePostIntent, error) {
	if len(content) > ContentMaxLenght {
		return CreatePostIntent{}, fmt.Errorf("post max-lenght exceed. Limit: %d", ContentMaxLenght)
	}

	return CreatePostIntent{
		UserId:  userId,
		Content: content,
	}, nil
}

func (sv PostService) Run(intent CreatePostIntent) (CreatePostResponse, error) {
	post, err := sv.Repository.Create(intent.UserId.String(), intent.Content)

	if err != nil {
		return CreatePostResponse{}, err
	}

	return CreatePostResponse{
		Id:      post.ID,
		Content: post.Content,
	}, nil
}
