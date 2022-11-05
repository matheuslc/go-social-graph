package repository

//go:generate mockgen -source=./user_repository.go -destination=../mock/repository/user_repository.go

import (
	"fmt"
	"gosocialgraph/pkg/entity"
	"time"

	"github.com/google/uuid"
	"github.com/neo4j/neo4j-go-driver/neo4j"
)

// UserRepository defines the Graph implementation of UserRepository
type UserRepository struct {
	Client neo4j.Driver
}

// UserReader defines a unit for finding an user
type UserReader interface {
	Find(userID uuid.UUID) (entity.User, error)
	FindByUsername(username string) (bool, error)
}

// Follower defines functions a follower must implement
type Follower interface {
	Follow(from, to string) (bool, error)
}

// Unfollower defines the small unit that needs to unfollow someone
type Unfollower interface {
	Unfollow(to, from string) (bool, error)
}

// Creater defines how to create a new user
type Creater interface {
	Create(username string) (entity.User, error)
}

// UserWriter compouns a write-only functions interface
type UserWriter interface {
	Creater
	Follower
	Unfollower
}

// UserReaderWriter groups reads and writes function
type UserReaderWriter interface {
	UserWriter
	UserReader
}

// Stats defines the methods needed to get user related stats
type Stats interface {
	CountFollowing(userID uuid.UUID) (int64, error)
	CountFollowers(userID uuid.UUID) (int64, error)
	CountPosts(userID uuid.UUID) (int64, error)
}

// Find finds a user by userID
func (repo *UserRepository) Find(userID uuid.UUID) (entity.User, error) {
	session, err := repo.Client.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	if err != nil {
		return entity.User{}, err
	}

	defer session.Close()

	persistedUser, err := session.ReadTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
		result, err := transaction.Run(
			"MATCH (a:User { uuid: $userId }) RETURN a LIMIT 1",
			map[string]interface{}{
				"userId": userID,
			},
		)

		if err != nil {
			return nil, err
		}

		if result.Next() {
			userProps := result.Record().GetByIndex(0).(neo4j.Node).Props()
			return entity.User{
				ID:        uuid.MustParse(userProps["uuid"].(string)),
				Username:  userProps["username"].(string),
				CreatedAt: userProps["created_at"].(time.Time),
			}, nil
		}

		return nil, result.Err()
	})

	if err != nil || persistedUser == nil {
		return entity.User{}, err
	}

	return persistedUser.(entity.User), nil
}

// FindByUsername finds a user by its username
func (repo *UserRepository) FindByUsername(username string) (bool, error) {
	session, err := repo.Client.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	if err != nil {
		return false, err
	}

	defer session.Close()

	persistedUser, err := session.ReadTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
		result, err := transaction.Run(
			"MATCH (a:User { username: $username }) RETURN a",
			map[string]interface{}{
				"username": username,
			},
		)

		if err != nil {
			return false, err
		}

		if result.Next() {
			return true, nil
		}

		return false, result.Err()
	})

	if err != nil {
		return false, err
	}

	return persistedUser.(bool), nil
}

// Follow match users and creates a relantionship between them
func (repo *UserRepository) Follow(to, from string) (bool, error) {
	session, err := repo.Client.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	if err != nil {
		return false, err
	}

	defer session.Close()

	_, err = session.WriteTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
		result, err := transaction.Run(
			"MATCH (a:User), (b:User) where a.uuid = $to AND b.uuid = $from CREATE (a)-[r:FOLLOW]->(b)",
			map[string]interface{}{
				"to":   to,
				"from": from,
			},
		)

		if err != nil {
			return nil, err
		}

		if result.Next() {
			return result.Record(), nil
		}

		return nil, result.Err()
	})

	if err != nil {
		return false, err
	}

	return true, nil
}

// Unfollow
func (repo *UserRepository) Unfollow(to, from string) (bool, error) {
	session, err := repo.Client.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	if err != nil {
		return false, err
	}

	defer session.Close()

	_, err = session.WriteTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
		result, err := transaction.Run(
			"MATCH (u:User { uuid: $to })-[r:FOLLOW]->(following:User { uuid: $from }) DELETE r ",
			map[string]interface{}{
				"to":   to,
				"from": from,
			},
		)

		if err != nil {
			return nil, err
		}

		if result.Next() {
			return result.Record(), nil
		}

		return nil, result.Err()
	})

	if err != nil {
		return false, err
	}

	return true, nil
}

// Create
func (repo *UserRepository) Create(username string) (entity.User, error) {
	session, err := repo.Client.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	if err != nil {
		return entity.User{}, err
	}
	defer session.Close()

	persistedUser, err := session.WriteTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
		result, err := transaction.Run(
			"CREATE (u:User) SET u.uuid = $uuid, u.username = $username, u.created_at = datetime($createdAt) RETURN u.uuid, u.username, u.created_at",
			map[string]interface{}{
				"uuid":      uuid.New().String(),
				"username":  username,
				"createdAt": time.Now().UTC(),
			},
		)

		if err != nil {
			return nil, err
		}

		if result.Next() {
			fmt.Println(result.Record().Keys())

			return entity.User{
				ID:        uuid.MustParse(result.Record().GetByIndex(0).(string)),
				Username:  result.Record().GetByIndex(1).(string),
				CreatedAt: result.Record().GetByIndex(2).(time.Time),
			}, nil
		}

		return nil, result.Err()
	})

	if err != nil {
		return entity.User{}, err
	}

	return persistedUser.(entity.User), nil
}

func (repo *UserRepository) CountFollowing(userID uuid.UUID) (int64, error) {
	session, err := repo.Client.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	if err != nil {
		return 0, err
	}

	defer session.Close()

	count, err := session.ReadTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
		result, err := transaction.Run(
			"MATCH (user:User { uuid: $userId })-[:FOLLOW]->(following) RETURN count(following)",
			map[string]interface{}{
				"userId": userID,
			},
		)

		if err != nil {
			return nil, err
		}

		if result.Next() {
			return result.Record().GetByIndex(0), nil
		}

		return nil, result.Err()
	})

	if err != nil {
		return 0, err
	}

	return count.(int64), nil
}

func (repo *UserRepository) CountFollowers(userID uuid.UUID) (int64, error) {
	session, err := repo.Client.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	if err != nil {
		return 0, err
	}

	defer session.Close()

	count, err := session.ReadTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
		result, err := transaction.Run(
			"MATCH (user:User { uuid: $userId })<-[:FOLLOW]-(followers) RETURN count(followers)",
			map[string]interface{}{
				"userId": userID,
			},
		)

		if err != nil {
			return nil, err
		}

		if result.Next() {
			return result.Record().GetByIndex(0), nil
		}

		return nil, result.Err()
	})

	if err != nil {
		return 0, err
	}

	return count.(int64), nil
}

func (repo *UserRepository) CountPosts(userID uuid.UUID) (int64, error) {
	session, err := repo.Client.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	if err != nil {
		return 0, err
	}
	defer session.Close()

	count, err := session.ReadTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
		result, err := transaction.Run(
			"MATCH (user:User { uuid: $userId })-[:TWEET|:REPOST]->(posts) RETURN count(posts)",
			map[string]interface{}{
				"userId": userID,
			},
		)

		if err != nil {
			return nil, err
		}

		if result.Next() {
			return result.Record().GetByIndex(0), nil
		}

		return nil, result.Err()
	})

	if err != nil {
		return 0, err
	}

	return count.(int64), nil
}
