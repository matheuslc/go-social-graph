package repository

//go:generate mockgen -source=./post_repository.go -destination=../mock/repository/post_repository.go

import (
	"fmt"
	"gosocialgraph/pkg/entity"
	"time"

	"github.com/google/uuid"
	"github.com/neo4j/neo4j-go-driver/neo4j"
)

// PostWriter defines the contact about how someone can write to storage a post
type PostWriter interface {
	Create(userID, content string) (entity.Post, error)
}

// Reposter defines the contract for how we can repost a post
type Reposter interface {
	Repost(user, parentID, quote string) error
}

// PostRepository holds the repository dependencies
type PostRepository struct {
	Client neo4j.Driver
}

// Create creates a new user post. The content right now is a simple string.
func (repo *PostRepository) Create(userID, content string) (entity.Post, error) {
	session, err := repo.Client.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	if err != nil {
		return entity.Post{}, fmt.Errorf("could not create a new session for Create query")
	}
	defer session.Close()

	persistedPost, err := session.WriteTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
		result, err := transaction.Run(
			"CREATE (a:Post) SET a.uuid = $uuid, a.content = $content, a.created_at = datetime($createdAt) RETURN a.uuid, a.content",
			map[string]interface{}{
				"uuid":      uuid.New().String(),
				"content":   content,
				"createdAt": time.Now().UTC(),
			})

		if err != nil {
			return nil, err
		}

		if result.Next() {
			return entity.Post{
				ID:      uuid.MustParse(result.Record().GetByIndex(0).(string)),
				Content: result.Record().GetByIndex(1).(string),
			}, nil
		}

		return nil, result.Err()
	})

	if err != nil {
		return entity.Post{}, err
	}

	_, err = session.WriteTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
		result, err := transaction.Run(
			"MATCH (a:User), (b:Post) WHERE a.uuid = $userId AND b.uuid = $postId CREATE(a)-[r:TWEET]->(b) CREATE(b)-[own:OWNS]->(a)",
			map[string]interface{}{
				"userId": userID,
				"postId": persistedPost.(entity.Post).ID.String(),
			})

		if err != nil {
			return nil, err
		}

		if result.Next() {
			return result.Record(), nil
		}

		return nil, result.Err()
	})

	if err != nil {
		return entity.Post{}, err
	}

	return persistedPost.(entity.Post), nil
}

// Repost creates a new post within the original post and more content
func (repo *PostRepository) Repost(userID, parentID, quote string) error {
	session, err := repo.Client.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	if err != nil {
		return fmt.Errorf("could not create a new session for Create query")
	}
	defer session.Close()

	_, err = session.WriteTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
		result, err := transaction.Run(
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
