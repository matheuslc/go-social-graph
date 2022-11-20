package db

import (
	"fmt"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

// New database
func New(uri, username, password string) (neo4j.DriverWithContext, error) {
	auth := neo4j.BasicAuth(username, password, "")

	driver, err := neo4j.NewDriverWithContext(uri, auth)
	if err != nil {
		panic(fmt.Errorf("Cant connect to Neo4j. Reason: %s", err))
	}

	return driver, err
}

// Graph
type Graph struct {
	Client neo4j.DriverWithContext
}
