package repository

//go:generate mockgen -source=./timeline_repository.go -destination=../mock/repository/timeline_repository.go

import (
	"context"
	"gosocialgraph/pkg/entity"
	"time"

	"github.com/google/uuid"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

// TimelineRepository defines what the timeline repository must holds
type TimelineRepository struct {
	Client neo4j.DriverWithContext
}

// TimelineReader defines the read operations for timeline
type TimelineReader interface {
	All(ctx context.Context) ([]entity.UserPost, error)
	TimelineFor(ctx context.Context, userID uuid.UUID) ([]entity.UserPost, error)
	UserPosts(ctx context.Context, userID uuid.UUID) ([]entity.UserPost, error)
}

// All return all posts and its users from the system
// You need to take care when using this method, and is perform a whole search
// into the database. If we start to hold too much data, we should paginate or remove
// this method.
func (repo *TimelineRepository) All(ctx context.Context) ([]entity.UserPost, error) {
	var list []entity.UserPost

	session := repo.Client.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close(ctx)

	_, err := session.ExecuteRead(ctx, func(transaction neo4j.ManagedTransaction) (interface{}, error) {
		result, err := transaction.Run(
			ctx,
			"MATCH (users:User)-[r:TWEET|REPOST]->(posts)-[ow:OWNS]->(us:User) RETURN posts, r, users, us ORDER BY posts.created_at",
			map[string]interface{}{},
		)

		if err != nil {
			return nil, err
		}

		list, err = dataMapper(ctx, result)
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

// TimelineFor generates a timeline for a specific user
// The timeline contains posts from following users
func (repo *TimelineRepository) TimelineFor(ctx context.Context, userID uuid.UUID) ([]entity.UserPost, error) {
	var list []entity.UserPost
	session := repo.Client.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close(ctx)

	_, err := session.ExecuteRead(ctx, func(transaction neo4j.ManagedTransaction) (interface{}, error) {
		result, err := transaction.Run(
			ctx,
			"MATCH (user:User {uuid: $currentUser})-[:FOLLOW]->(followers)-[r:TWEET|REPOST]->(posts:Post)-[ow:OWNS]->(us:User) RETURN posts, r, followers, us ORDER BY posts.created_at",
			map[string]interface{}{
				"currentUser": userID.String(),
			},
		)

		if err != nil {
			return nil, err
		}

		list, err = dataMapper(ctx, result)
		if err != nil {
			return nil, err
		}

		return list, nil
	})

	if err != nil {
		return []entity.UserPost{}, err
	}

	return list, nil
}

// UserPosts retuns the posts from an user
func (repo *TimelineRepository) UserPosts(ctx context.Context, userID uuid.UUID) ([]entity.UserPost, error) {
	var list []entity.UserPost
	session := repo.Client.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close(ctx)

	_, err := session.ExecuteRead(ctx, func(transaction neo4j.ManagedTransaction) (interface{}, error) {
		result, err := transaction.Run(
			ctx,
			"MATCH (user:User {uuid: $currentUser})-[r:TWEET|REPOST]->(posts)-[ow:OWNS]->(us:User) RETURN posts, r, user, us",
			map[string]interface{}{
				"currentUser": userID.String(),
			},
		)

		if err != nil {
			return nil, err
		}

		list, err = dataMapper(ctx, result)
		if err != nil {
			return nil, err
		}

		return list, nil
	})

	if err != nil {
		return []entity.UserPost{}, err
	}

	return list, nil
}

func dataMapper(ctx context.Context, result neo4j.ResultWithContext) ([]entity.UserPost, error) {
	var list []entity.UserPost

	for result.Next(ctx) {
		values := result.Record().Values
		postRecord := values[0].(neo4j.Node)
		relation := values[1].(neo4j.Relationship)
		userRecord := values[2].(neo4j.Node)

		postProps := postRecord.Props
		parsedPost := entity.Post{
			ID:        uuid.MustParse(postProps["uuid"].(string)),
			Content:   postProps["content"].(string),
			CreatedAt: postProps["created_at"].(time.Time),
		}

		userProps := userRecord.Props
		postIsFrom := entity.User{
			ID:       uuid.MustParse(userProps["uuid"].(string)),
			Username: userProps["username"].(string),
		}

		if len(relation.Props) > 0 {
			props := relation.Props
			quote := props["quote"]
			id := props["uuid"]
			createdAt := props["created_at"]

			if quote == nil {
				quote = ""
			}

			userRepost := values[3].(neo4j.Node)
			userRepostProps := userRepost.Props
			repostWithUser := entity.UserPost{
				User: entity.User{
					ID:       uuid.MustParse(userRepostProps["uuid"].(string)),
					Username: userRepostProps["username"].(string),
				},
				Post: parsedPost,
			}

			repost := entity.Post{
				ID:        uuid.MustParse(id.(string)),
				Parent:    &repostWithUser,
				Content:   quote.(string),
				CreatedAt: createdAt.(time.Time),
			}

			userPost := entity.UserPost{User: postIsFrom, Post: repost}

			list = append(list, userPost)
		} else {
			list = append(list, entity.UserPost{
				User: postIsFrom,
				Post: parsedPost,
			})
		}
	}

	return list, nil
}
