package service

import (
	"gosocialgraph/pkg/entity"
	"gosocialgraph/pkg/repository"
)

type AllService struct {
	Repository repository.TimelineReader
}

type AllPostResponse struct {
	Posts []entity.UserPost `json:"posts"`
}

func (sv AllService) Run() (AllPostResponse, error) {
	response, err := sv.Repository.All()

	if err != nil {
		return AllPostResponse{}, err
	}

	return AllPostResponse{
		Posts: response,
	}, nil
}
