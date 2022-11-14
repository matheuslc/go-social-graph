package repository

//go:generate mockgen -source=./user_repository.go -destination=../mock/repository/user_repository.go

import (
	"fmt"
	"gosocialgraph/pkg/entity"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/neo4j/neo4j-go-driver/neo4j"
	"golang.org/x/crypto/bcrypt"
)

var Salt = os.Getenv("JWT_SALT")

// UserRepository defines the Graph implementation of UserRepository
type UserRepository struct {
	Client neo4j.Driver
}

// UserReader defines a unit for finding an user
type UserReader interface {
	Find(userID uuid.UUID) (entity.User, error)
	FindByUsername(username string) (entity.User, error)
}

// Follower defines functions a follower must implement
type Follower interface {
	Follow(from, to string) error
}

// Unfollower defines the small unit that needs to unfollow someone
type Unfollower interface {
	Unfollow(to, from string) error
}

// Creater defines how to create a new user
type Creater interface {
	Create(username, password string) (entity.User, error)
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
				"userId": userID.String(),
			},
		)

		if err != nil {
			return nil, err
		}

		if result.Next() {
			userProps := result.Record().GetByIndex(0).(neo4j.Node).Props()
			uuidParsed := uuid.MustParse(userProps["uuid"].(string))

			return entity.User{
				ID:        uuidParsed,
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
func (repo *UserRepository) FindByUsername(username string) (user entity.User, err error) {
	session, err := repo.Client.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	if err != nil {
		return user, err
	}

	defer session.Close()

	userFromDatabase, err := session.ReadTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
		result, err := transaction.Run(
			"MATCH (a:User { username: $username }) RETURN a LIMIT 1",
			map[string]interface{}{
				"username": username,
			},
		)

		if err != nil {
			return entity.User{}, err
		}

		if result.Next() {
			userProps := result.Record().GetByIndex(0).(neo4j.Node).Props()
			uuidParsed := uuid.MustParse(userProps["uuid"].(string))

			return entity.User{
				ID:        uuidParsed,
				Username:  userProps["username"].(string),
				Password:  userProps["password"].([]byte),
				CreatedAt: userProps["created_at"].(time.Time),
			}, nil
		}

		return entity.User{}, result.Err()
	})

	if err != nil {
		return user, err
	}

	return userFromDatabase.(entity.User), nil
}

// Follow match users and creates a relantionship between them
func (repo *UserRepository) Follow(to, from string) error {
	session, err := repo.Client.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	if err != nil {
		return err
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
		return err
	}

	return nil
}

// Unfollow
func (repo *UserRepository) Unfollow(to, from string) error {
	session, err := repo.Client.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	if err != nil {
		return err
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
		return err
	}

	return nil
}

// Create
func (repo *UserRepository) Create(username, password string) (entity.User, error) {
	session, err := repo.Client.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	if err != nil {
		return entity.User{}, err
	}
	defer session.Close()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(fmt.Sprintf("%s%s", Salt, password)), 8)
	if err != nil {
		return entity.User{}, err
	}

	persistedUser, err := session.WriteTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
		result, err := transaction.Run(
			"CREATE (u:User) SET u.uuid = $uuid, u.username = $username, u.password = $password, u.created_at = datetime($createdAt) RETURN u.uuid, u.username, u.created_at",
			map[string]interface{}{
				"uuid":      uuid.New().String(),
				"username":  username,
				"password":  hashedPassword,
				"createdAt": time.Now().UTC(),
			},
		)

		if err != nil {
			return nil, err
		}

		if result.Next() {
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
				"userId": userID.String(),
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
				"userId": userID.String(),
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
				"userId": userID.String(),
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
