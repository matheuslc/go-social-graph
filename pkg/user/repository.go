package user

import (
	"time"

	"github.com/google/uuid"
	"github.com/neo4j/neo4j-go-driver/neo4j"
)

// Repository defines the Graph implementation of UserRepository
type Repository struct {
	Client neo4j.Driver
}

// Reader defines a unit for finding an user
type Reader interface {
	Find(userID string) (User, error)
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
	Create(username string) (User, error)
}

// Writer compouns a write-only functions interface
type Writer interface {
	Creater
	Follower
	Unfollower
}

// ReaderWriter groups reads and writes function
type ReaderWriter interface {
	Creater
	Reader
}

// Stats defines the methods needed to get user related stats
type Stats interface {
	CountFollowing(userID string) (int64, error)
	CountFollowers(userID string) (int64, error)
	CountPosts(userID string) (int64, error)
}

func (repo *Repository) Find(userID string) (User, error) {
	session, err := repo.Client.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	if err != nil {
		return User{}, err
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
			return User{
				ID:        uuid.MustParse(userProps["uuid"].(string)),
				Username:  userProps["username"].(string),
				CreatedAt: userProps["created_at"].(time.Time),
			}, nil
		}

		return nil, result.Err()
	})

	if err != nil || persistedUser == nil {
		return User{}, err
	}

	return persistedUser.(User), nil
}

func (repo *Repository) FindByUsername(username string) (bool, error) {
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
			return nil, err
		}

		if result.Next() {
			return true, nil
		}

		return nil, result.Err()
	})

	if err != nil {
		return false, err
	}

	return persistedUser.(bool), nil
}

func (repo *Repository) Follow(to, from string) (bool, error) {
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

func (repo *Repository) Unfollow(to, from string) (bool, error) {
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

func (repo *Repository) Create(username string) (User, error) {
	session, err := repo.Client.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	if err != nil {
		return User{}, err
	}
	defer session.Close()

	persistedUser, err := session.WriteTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
		result, err := transaction.Run(
			"CREATE (u:User) SET u.uuid = $uuid, u.username = $username, u.created_at = datetime($createdAt) RETURN u.uuid, u.username",
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
			return User{
				ID:       uuid.MustParse(result.Record().GetByIndex(0).(string)),
				Username: result.Record().GetByIndex(1).(string),
			}, nil
		}

		return nil, result.Err()
	})

	if err != nil {
		return User{}, err
	}

	return persistedUser.(User), nil
}

func (repo *Repository) CountFollowing(userId string) (int64, error) {
	session, err := repo.Client.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	if err != nil {
		return 0, err
	}

	defer session.Close()

	count, err := session.ReadTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
		result, err := transaction.Run(
			"MATCH (user:User { uuid: $userId })-[:FOLLOW]->(following) RETURN count(following)",
			map[string]interface{}{
				"userId": userId,
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

func (repo *Repository) CountFollowers(userId string) (int64, error) {
	session, err := repo.Client.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	if err != nil {
		return 0, err
	}

	defer session.Close()

	count, err := session.ReadTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
		result, err := transaction.Run(
			"MATCH (user:User { uuid: $userId })<-[:FOLLOW]-(followers) RETURN count(followers)",
			map[string]interface{}{
				"userId": userId,
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

func (repo *Repository) CountPosts(userId string) (int64, error) {
	session, err := repo.Client.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	if err != nil {
		return 0, err
	}
	defer session.Close()

	count, err := session.ReadTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
		result, err := transaction.Run(
			"MATCH (user:User { uuid: $userId })-[:TWEET|:REPOST]->(posts) RETURN count(posts)",
			map[string]interface{}{
				"userId": userId,
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
