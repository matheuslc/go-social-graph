package post

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/neo4j/neo4j-go-driver/neo4j"
)

type Writer interface {
	Create(userId, content string) (Post, error)
}

type Reposter interface {
	Repost(user, parentId, quote string) (bool, error)
}

type Repository struct {
	Client neo4j.Driver
}

func (repo Repository) Create(userId, content string) (Post, error) {
	session, err := repo.Client.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	if err != nil {
		return Post{}, fmt.Errorf("could not create a new session for Create query")
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
			return Post{
				Id:      uuid.MustParse(result.Record().GetByIndex(0).(string)),
				Content: result.Record().GetByIndex(1).(string),
			}, nil
		}

		return nil, result.Err()
	})

	if err != nil {
		return Post{}, err
	}

	_, err = session.WriteTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
		result, err := transaction.Run(
			"MATCH (a:User), (b:Post) WHERE a.uuid = $userId AND b.uuid = $postId CREATE(a)-[r:TWEET]->(b) CREATE(b)-[own:OWNS]->(a)",
			map[string]interface{}{
				"userId": userId,
				"postId": persistedPost.(Post).Id.String(),
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
		return Post{}, err
	}

	return persistedPost.(Post), nil
}

func (repo Repository) Repost(userId, parentId, quote string) (bool, error) {
	session, err := repo.Client.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	if err != nil {
		return false, fmt.Errorf("could not create a new session for Create query")
	}
	defer session.Close()

	_, err = session.WriteTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
		result, err := transaction.Run(
			"MATCH (a:User), (b:Post) WHERE a.uuid = $userId AND b.uuid = $postId CREATE(a)-[r:REPOST { uuid: $uuid, created_at: datetime($createdAt), quote: $quote}]->(b)",
			map[string]interface{}{
				"userId":    userId,
				"postId":    parentId,
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
		return false, err
	}

	return true, nil
}
