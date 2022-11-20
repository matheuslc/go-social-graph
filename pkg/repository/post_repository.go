package repository

//go:generate mockgen -source=./post_repository.go -destination=../mock/repository/post_repository.go

import (
	"context"
	"gosocialgraph/pkg/entity"
	"time"

	"github.com/google/uuid"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

// PostWriter defines the contact about how someone can write to storage a post
type PostWriter interface {
	Create(ctx context.Context, userID, content string) (entity.Post, error)
}

// Reposter defines the contract for how we can repost a post
type Reposter interface {
	Repost(ctx context.Context, user, parentID, quote string) error
}

// PostRepository holds the repository dependencies
type PostRepository struct {
	Client neo4j.DriverWithContext
}

// Create creates a new user post. The content right now is a simple string.
func (repo *PostRepository) Create(ctx context.Context, userID, content string) (entity.Post, error) {
	session := repo.Client.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close(ctx)

	persistedPost, err := session.ExecuteWrite(ctx, func(transaction neo4j.ManagedTransaction) (interface{}, error) {
		result, err := transaction.Run(
			ctx,
			"CREATE (a:Post) SET a.uuid = $uuid, a.content = $content, a.created_at = datetime($createdAt) RETURN a.uuid, a.content",
			map[string]interface{}{
				"uuid":      uuid.New().String(),
				"content":   content,
				"createdAt": time.Now().UTC(),
			})

		if err != nil {
			return nil, err
		}

		records, err := result.Collect(ctx)
		if err != nil {
			return nil, err
		}

		values := records[0].Values
		return entity.Post{
			ID:      uuid.MustParse(values[0].(string)),
			Content: values[1].(string),
		}, nil
	})

	if err != nil {
		return entity.Post{}, err
	}

	_, err = session.ExecuteWrite(ctx, func(transaction neo4j.ManagedTransaction) (interface{}, error) {
		_, err := transaction.Run(
			ctx,
			"MATCH (a:User), (b:Post) WHERE a.uuid = $userId AND b.uuid = $postId CREATE(a)-[r:TWEET]->(b) CREATE(b)-[own:OWNS]->(a)",
			map[string]interface{}{
				"userId": userID,
				"postId": persistedPost.(entity.Post).ID.String(),
			})

		if err != nil {
			return nil, err
		}

		return "ok", nil
	})

	if err != nil {
		return entity.Post{}, err
	}

	return persistedPost.(entity.Post), nil
}

// Repost creates a new post within the original post and more content
func (repo *PostRepository) Repost(ctx context.Context, userID, parentID, quote string) error {
	session := repo.Client.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close(ctx)

	_, err := session.ExecuteWrite(ctx, func(transaction neo4j.ManagedTransaction) (interface{}, error) {
		result, err := transaction.Run(
			ctx,
			"MATCH (a:User), (b:Post) WHERE a.uuid = $userId AND b.uuid = $postId CREATE(a)-[r:REPOST { uuid: $uuid, created_at: datetime($createdAt), quote: $quote}]->(b)",
			map[string]interface{}{
				"userId":    userID,
				"postId":    parentID,
				"uuid":      uuid.New().String(),
				"createdAt": time.Now().UTC(),
				"quote":     quote,
			})

		if err != nil {
			return nil, err
		}

		return result, err
	})

	if err != nil {
		return err
	}

	return nil
}
