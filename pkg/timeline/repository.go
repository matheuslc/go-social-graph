package timeline

import (
	"fmt"
	"gosocialgraph/pkg/post"
	"gosocialgraph/pkg/user"
	"time"

	"github.com/google/uuid"
	"github.com/neo4j/neo4j-go-driver/neo4j"
)

type Repository struct {
	Client neo4j.Driver
}

type Reader interface {
	All() ([]post.UserPost, error)
	TimelineFor(userId uuid.UUID) ([]post.UserPost, error)
	UserPosts(userId uuid.UUID) ([]post.UserPost, error)
}

func (repo Repository) All() ([]post.UserPost, error) {
	var list []post.UserPost

	session, err := repo.Client.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	if err != nil {
		return list, fmt.Errorf("could not create a new session for Create query")
	}

	defer session.Close()

	_, err = session.ReadTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
		result, err := transaction.Run(
			"MATCH (users:User)-[r:TWEET|REPOST]->(posts)-[ow:OWNS]->(us:User) RETURN posts, r, users, us ORDER BY posts.created_at",
			map[string]interface{}{},
		)

		if err != nil {
			return nil, err
		}

		list, err = dataMapper(result)
		if err != nil {
			return nil, err
		}

		return list, nil
	})

	if err != nil {
		return list, err
	}

	return list, nil
}

func (repo Repository) TimelineFor(userId uuid.UUID) ([]post.UserPost, error) {
	var list []post.UserPost
	session, err := repo.Client.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	if err != nil {
		return nil, fmt.Errorf("could not create a new session for Create query")
	}

	defer session.Close()

	_, err = session.ReadTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
		result, err := transaction.Run(
			"MATCH (user:User {uuid: $currentUser})-[:FOLLOW]->(followers)-[r:TWEET|REPOST]->(posts:Post)-[ow:OWNS]->(us:User) RETURN posts, r, followers, us ORDER BY posts.created_at",
			map[string]interface{}{
				"currentUser": userId.String(),
			},
		)

		if err != nil {
			return nil, err
		}

		list, err = dataMapper(result)
		if err != nil {
			return nil, err
		}

		return list, nil
	})

	if err != nil {
		return []post.UserPost{}, err
	}

	return list, nil
}

func (repo Repository) UserPosts(userId uuid.UUID) ([]post.UserPost, error) {
	var list []post.UserPost
	session, err := repo.Client.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	if err != nil {
		return nil, fmt.Errorf("could not create a new session for Create query")
	}

	defer session.Close()

	_, err = session.ReadTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
		result, err := transaction.Run(
			"MATCH (user:User {uuid: $currentUser})-[r:TWEET|REPOST]->(posts)-[ow:OWNS]->(us:User) RETURN posts, r, user, us",
			map[string]interface{}{
				"currentUser": userId.String(),
			},
		)

		if err != nil {
			return nil, err
		}

		list, err = dataMapper(result)
		if err != nil {
			return nil, err
		}

		return list, nil
	})

	if err != nil {
		return []post.UserPost{}, err
	}

	return list, nil
}

func dataMapper(result neo4j.Result) ([]post.UserPost, error) {
	var list []post.UserPost

	for result.Next() {
		postRecord := result.Record().GetByIndex(0).(neo4j.Node)
		relation := result.Record().GetByIndex(1).(neo4j.Relationship)
		userRecord := result.Record().GetByIndex(2).(neo4j.Node)

		postProps := postRecord.Props()
		parsedPost := post.Post{
			Id:        uuid.MustParse(postProps["uuid"].(string)),
			Content:   postProps["content"].(string),
			CreatedAt: postProps["created_at"].(time.Time),
		}

		userProps := userRecord.Props()
		postIsFrom := user.User{
			Id:       uuid.MustParse(userProps["uuid"].(string)),
			Username: userProps["username"].(string),
		}

		if len(relation.Props()) > 0 {
			props := relation.Props()
			quote := props["quote"]
			id := props["uuid"]
			createdAt := props["created_at"]

			if quote == nil {
				quote = ""
			}

			userRepost := result.Record().GetByIndex(3).(neo4j.Node)
			userRepostProps := userRepost.Props()
			repostWithUser := post.UserPost{
				User: user.User{
					Id:       uuid.MustParse(userRepostProps["uuid"].(string)),
					Username: userRepostProps["username"].(string),
				},
				Post: parsedPost,
			}

			repost := post.Repost{
				Id:        uuid.MustParse(id.(string)),
				Quote:     quote.(string),
				Parent:    repostWithUser,
				CreatedAt: createdAt.(time.Time),
			}

			userPost := post.UserPost{
				User: postIsFrom,
				Post: repost,
			}

			list = append(list, userPost)
		} else {
			list = append(list, post.UserPost{
				User: postIsFrom,
				Post: parsedPost,
			})
		}
	}

	return list, nil
}
