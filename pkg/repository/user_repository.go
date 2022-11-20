package repository

//go:generate mockgen -source=./user_repository.go -destination=../mock/repository/user_repository.go

import (
	"context"
	"fmt"
	"gosocialgraph/pkg/entity"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"golang.org/x/crypto/bcrypt"
)

var Salt = os.Getenv("JWT_SALT")

// UserRepository defines the Graph implementation of UserRepository
type UserRepository struct {
	Client neo4j.DriverWithContext
}

// UserReader defines a unit for finding an user
type UserReader interface {
	Find(ctx context.Context, userID uuid.UUID) (entity.User, error)
	FindByUsername(ctx context.Context, username string) (entity.User, error)
}

// Follower defines functions a follower must implement
type Follower interface {
	Follow(ctx context.Context, from, to string) error
}

// Unfollower defines the small unit that needs to unfollow someone
type Unfollower interface {
	Unfollow(ctx context.Context, to, from string) error
}

// Creater defines how to create a new user
type Creater interface {
	Create(ctx context.Context, username, email, password string) (entity.User, error)
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
	CountFollowing(ctx context.Context, userID uuid.UUID) (int64, error)
	CountFollowers(ctx context.Context, userID uuid.UUID) (int64, error)
	CountPosts(ctx context.Context, userID uuid.UUID) (int64, error)
}

// Find finds a user by userID
func (repo *UserRepository) Find(ctx context.Context, userID uuid.UUID) (entity.User, error) {
	session := repo.Client.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer session.Close(ctx)

	persistedUser, err := session.ExecuteRead(ctx, func(transaction neo4j.ManagedTransaction) (interface{}, error) {
		result, err := transaction.Run(
			ctx,
			"MATCH (a:User { uuid: $userId }) RETURN a LIMIT 1",
			map[string]interface{}{
				"userId": userID.String(),
			},
		)

		if err != nil {
			return nil, err
		}

		if result.Next(ctx) {
			values := result.Record().Values
			userProps := values[0].(neo4j.Node).Props
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
func (repo *UserRepository) FindByUsername(ctx context.Context, username string) (user entity.User, err error) {
	session := repo.Client.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer session.Close(ctx)

	userFromDatabase, err := session.ExecuteRead(ctx, func(transaction neo4j.ManagedTransaction) (interface{}, error) {
		result, err := transaction.Run(
			ctx,
			"MATCH (a:User { username: $username }) RETURN a LIMIT 1",
			map[string]interface{}{
				"username": username,
			},
		)

		if err != nil {
			return entity.User{}, err
		}

		if result.Next(ctx) {
			values := result.Record().Values
			userProps := values[0].(neo4j.Node).Props
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
func (repo *UserRepository) Follow(ctx context.Context, to, from string) error {
	session := repo.Client.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close(ctx)

	_, err := session.ExecuteWrite(ctx, func(transaction neo4j.ManagedTransaction) (interface{}, error) {
		result, err := transaction.Run(
			ctx,
			"MATCH (a:User), (b:User) where a.uuid = $to AND b.uuid = $from CREATE (a)-[r:FOLLOW]->(b)",
			map[string]interface{}{
				"to":   to,
				"from": from,
			},
		)

		if err != nil {
			return nil, err
		}

		if result.Next(ctx) {
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
func (repo *UserRepository) Unfollow(ctx context.Context, to, from string) error {
	session := repo.Client.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close(ctx)

	_, err := session.ExecuteWrite(ctx, func(transaction neo4j.ManagedTransaction) (interface{}, error) {
		result, err := transaction.Run(
			ctx,
			"MATCH (u:User { uuid: $to })-[r:FOLLOW]->(following:User { uuid: $from }) DELETE r ",
			map[string]interface{}{
				"to":   to,
				"from": from,
			},
		)

		if err != nil {
			return nil, err
		}

		if result.Next(ctx) {
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
func (repo *UserRepository) Create(ctx context.Context, username, email, password string) (entity.User, error) {
	session := repo.Client.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close(ctx)

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(fmt.Sprintf("%s%s", Salt, password)), 8)
	if err != nil {
		return entity.User{}, err
	}

	persistedUser, err := session.ExecuteWrite(ctx, func(transaction neo4j.ManagedTransaction) (interface{}, error) {
		result, err := transaction.Run(
			ctx,
			"CREATE (u:User) SET u.uuid = $uuid, u.username = $username, u.email = $email, u.password = $password, u.created_at = datetime($createdAt) RETURN u.uuid, u.username, u.email, u.created_at",
			map[string]interface{}{
				"uuid":      uuid.New().String(),
				"username":  username,
				"email":     email,
				"password":  hashedPassword,
				"createdAt": time.Now().UTC(),
			},
		)

		if err != nil {
			return nil, err
		}

		if result.Next(ctx) {
			values := result.Record().Values
			return entity.User{
				ID:        uuid.MustParse(values[0].(string)),
				Username:  values[1].(string),
				Email:     values[2].(string),
				CreatedAt: values[3].(time.Time),
			}, nil
		}

		return nil, result.Err()
	})

	if err != nil {
		return entity.User{}, err
	}

	return persistedUser.(entity.User), nil
}

func (repo *UserRepository) CountFollowing(ctx context.Context, userID uuid.UUID) (int64, error) {
	session := repo.Client.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close(ctx)

	count, err := session.ExecuteRead(ctx, func(transaction neo4j.ManagedTransaction) (interface{}, error) {
		result, err := transaction.Run(
			ctx,
			"MATCH (user:User { uuid: $userId })-[:FOLLOW]->(following) RETURN count(following)",
			map[string]interface{}{
				"userId": userID.String(),
			},
		)

		if err != nil {
			return nil, err
		}

		if result.Next(ctx) {
			values := result.Record().Values
			return values[0], nil
		}

		return nil, result.Err()
	})

	if err != nil {
		return 0, err
	}

	return count.(int64), nil
}

func (repo *UserRepository) CountFollowers(ctx context.Context, userID uuid.UUID) (int64, error) {
	session := repo.Client.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close(ctx)

	count, err := session.ExecuteRead(ctx, func(transaction neo4j.ManagedTransaction) (interface{}, error) {
		result, err := transaction.Run(
			ctx,
			"MATCH (user:User { uuid: $userId })<-[:FOLLOW]-(followers) RETURN count(followers)",
			map[string]interface{}{
				"userId": userID.String(),
			},
		)

		if err != nil {
			return nil, err
		}

		if result.Next(ctx) {
			values := result.Record().Values
			return values[0], nil
		}

		return nil, result.Err()
	})

	if err != nil {
		return 0, err
	}

	return count.(int64), nil
}

func (repo *UserRepository) CountPosts(ctx context.Context, userID uuid.UUID) (int64, error) {
	session := repo.Client.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close(ctx)

	count, err := session.ExecuteRead(ctx, func(transaction neo4j.ManagedTransaction) (interface{}, error) {
		result, err := transaction.Run(
			ctx,
			"MATCH (user:User { uuid: $userId })-[:TWEET|:REPOST]->(posts) RETURN count(posts)",
			map[string]interface{}{
				"userId": userID.String(),
			},
		)

		if err != nil {
			return nil, err
		}

		if result.Next(ctx) {
			values := result.Record().Values
			return values[0], nil
		}

		return nil, result.Err()
	})

	if err != nil {
		return 0, err
	}

	return count.(int64), nil
}
