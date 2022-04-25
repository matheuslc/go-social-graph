package timeline

import "gosocialgraph/pkg/post"

type AllService struct {
	Repository Reader
}

type AllPostResponse struct {
	Posts []post.UserPost `json:"posts"`
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
