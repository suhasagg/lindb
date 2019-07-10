package service

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/eleme/lindb/models"
	"github.com/eleme/lindb/pkg/state"
)

const databaseConfigNode = "/lindb/database/config"

// DatabaseService defines database service interface
type DatabaseService interface {
	// Save saves database config
	Save(database models.Database) error
	// Get gets database config by name
	Get(name string) (models.Database, error)
}

// databaseService implements DatabaseService interface
type databaseService struct {
	repo state.Repository
}

// NewDatabaseService creates database service
func NewDatabaseService(repo state.Repository) DatabaseService {
	return &databaseService{
		repo: repo,
	}
}

// Save saves database config into state's repo
func (db *databaseService) Save(database models.Database) error {
	if len(database.Name) == 0 {
		return fmt.Errorf("name cannot be empty")
	}
	if database.NumOfShard <= 0 {
		return fmt.Errorf("num. of shard must be > 0")
	}
	if database.ReplicaFactor <= 0 {
		return fmt.Errorf("replica factor must be > 0")
	}
	data, err := json.Marshal(database)
	if err != nil {
		return fmt.Errorf("marshal database config error:%s", err)
	}
	return db.repo.Put(context.TODO(), getDatabasePath(database.Name), data)
}

// Get returns the database config in the state's repo
func (db *databaseService) Get(name string) (models.Database, error) {
	database := models.Database{}
	if name == "" {
		return database, fmt.Errorf("database name must not be null")
	}
	configBytes, err := db.repo.Get(context.TODO(), getDatabasePath(name))
	if err != nil {
		return database, err
	}
	err = json.Unmarshal(configBytes, &database)
	if err != nil {
		return database, err
	}
	return database, nil
}

// getDatabasePath gets the path where the database config is stored
func getDatabasePath(databaseName string) string {
	return databaseConfigNode + "/" + databaseName
}
