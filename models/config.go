package models

import (
	"time"
)

type Config struct {
	Key   string
	Value string

	// metadata
	id        int
	createdAt time.Time
	updatedAt time.Time
}

const sqlReadConfig = `
SELECT id, key, value, created_at, updated_at
FROM config
WHERE key = $1;`

func ReadConfig(k string) (*Config, error) {

	var id int

	var key string
	var value string

	var createdAt time.Time
	var updatedAt time.Time

	err := db.QueryRow(sqlReadConfig, k).Scan(&id, &key, &value, &createdAt, &updatedAt)
	if err != nil {
		logger.Tracef("ReadConfig(%s) (nil, %s)", k, err)
		return nil, err
	}

	newConfig := &Config{
		id:        id,
		Key:       key,
		Value:     value,
		createdAt: createdAt,
		updatedAt: updatedAt,
	}

	logger.Tracef("ReadConfig(%s) (%v, nil)", k, &newConfig)
	return newConfig, nil
}

const sqlCreateConfig = `
INSERT INTO "public"."config" (key, value)
VALUES ($1, $2)
RETURNING id, created_at, updated_at;`

func CreateConfig(k string, v string) (*Config, error) {
	var newId int
	var newCreatedAt time.Time
	var newUpdatedAt time.Time

	err := db.QueryRow(sqlCreateConfig, k, v).Scan(&newId, &newCreatedAt, &newUpdatedAt)
	if err != nil {
		logger.Tracef("CreateConfig(%s, %s) (nil, %s)", k, v, err)
		return nil, err
	}

	newConfig := &Config{
		id:        newId,
		Key:       k,
		Value:     v,
		createdAt: newCreatedAt,
		updatedAt: newUpdatedAt,
	}

	logger.Tracef("CreateConfig(%s, %s) (%s, nil)", k, v, &newConfig)
	return newConfig, nil
}
