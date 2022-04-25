package db

import (
	"fmt"

	"github.com/neo4j/neo4j-go-driver/neo4j"
)

// Graph database
func New(uri, username, password string) (neo4j.Driver, error) {
	driver, err := neo4j.NewDriver(uri, neo4j.BasicAuth(username, password, ""), func(c *neo4j.Config) {
		c.Encrypted = false
	})

	if err != nil {
		fmt.Printf("Cant connect to Neo4j. Reason: %s", err)
	}

	return driver, err
}

type Graph struct {
	Client neo4j.Driver
}

