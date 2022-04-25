package user

import (
	"time"

	"github.com/google/uuid"
	"github.com/neo4j/neo4j-go-driver/neo4j"
)

type UserRepository struct {
	Client neo4j.Driver
}

type Reader interface {
	Find(userId string) (User, error)
}

type Writer interface {
	Create(username string) (User, error)
	Follow(to, from string) (bool, error)
	Unfollow(to, from string) (bool, error)
}

type Stats interface {
	CountFollowing(userId string) (int64, error)
	CountFollowers(userId string) (int64, error)
	CountPosts(userId string) (int64, error)
}

func (repo UserRepository) Find(userId string) (User, error) {
	session, err := repo.Client.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	if err != nil {
		return User{}, err
	}

	defer session.Close()

	persistedUser, err := session.ReadTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
		result, err := transaction.Run(
			"MATCH (a:User { uuid: $userId }) RETURN a LIMIT 1",
			map[string]interface{}{
				"userId": userId,
			},
		)

		if err != nil {
			return nil, err
		}

		if result.Next() {
			userProps := result.Record().GetByIndex(0).(neo4j.Node).Props()
			return User{
				Id:        uuid.MustParse(userProps["uuid"].(string)),
				Username:  userProps["username"].(string),
				CreatedAt: userProps["created_at"].(time.Time),
			}, nil
		}

		return nil, result.Err()
	})

	if err != nil {
		return User{}, err
	}

	return persistedUser.(User), nil
}

func (repo UserRepository) Follow(to, from string) (bool, error) {
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

func (repo UserRepository) Unfollow(to, from string) (bool, error) {
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

func (repo UserRepository) Create(username string) (User, error) {
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
				Id:       uuid.MustParse(result.Record().GetByIndex(0).(string)),
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

func (repo UserRepository) CountFollowing(userId string) (int64, error) {
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

func (repo UserRepository) CountFollowers(userId string) (int64, error) {
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

func (repo UserRepository) CountPosts(userId string) (int64, error) {
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
