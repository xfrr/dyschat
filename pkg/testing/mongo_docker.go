package testing

import (
	"context"
	"fmt"
	"log"

	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/xfrr/dyschat/pkg/mongo"

	mongodb "go.mongodb.org/mongo-driver/mongo"
)

type MongoDockerContainer struct {
	database *mongodb.Database
	pool     *dockertest.Pool
	resource *dockertest.Resource
}

func (dc MongoDockerContainer) URI() string {
	return fmt.Sprintf("mongodb://localhost:%s", dc.resource.GetPort("27017/tcp"))
}

func (dc MongoDockerContainer) DatabaseName() string {
	return "mongo_pkg_test"
}

func NewMongoContainer() *MongoDockerContainer {
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	dc := &MongoDockerContainer{
		pool: pool,
	}

	return dc
}

func (dc *MongoDockerContainer) Start(ctx context.Context, database string) (*mongodb.Database, error) {
	var err error
	// pull docker image for version
	dc.resource, err = dc.pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "mongo",
		Tag:        "5.0",
		Env: []string{
			// username and password for mongodb superuser
			"MONGO_INITDB_ROOT_USERNAME=test",
			"MONGO_INITDB_ROOT_PASSWORD=test",
		},
	}, func(config *docker.HostConfig) {
		// set AutoRemove to true so that stopped container goes away by itself
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{
			Name: "no",
		}
	})
	if err != nil {
		return nil, err
	}

	dc.resource.Expire(300)

	// exponential backoff-retry, because the application in the container might not be ready to accept connections yet
	err = dc.pool.Retry(func() error {
		var err error

		// creates new nats connection
		dc.database, err = mongo.Connect(ctx, mongo.Config{
			URI: fmt.Sprintf("mongodb://%s:%s@localhost:%s",
				"test", "test", dc.resource.GetPort("27017/tcp")),
			Database: database,
		})
		if err != nil {
			return err
		}

		return err
	})
	if err != nil {
		return nil, fmt.Errorf("could not connect to docker: %w", err)
	}

	return dc.database, nil
}

func (dc MongoDockerContainer) Purge() {
	if err := dc.pool.Purge(dc.resource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}
}
