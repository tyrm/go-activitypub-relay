package models

import (
	"time"
)

type BlacklistItem struct {
	Hostname string

	CreatedAt   time.Time

	// metadata
	id int
}

const sqlReadBlacklistItem = `
SELECT id, hostname, created_at
FROM blacklist
WHERE hostname = $1;`

func ReadBlacklistItem(h string) (*BlacklistItem, error) {
	var id int

	var hostname string

	var createdAt time.Time

	err := db.QueryRow(sqlReadBlacklistItem, h).Scan(&id, &hostname, &createdAt)
	if err != nil {
		logger.Tracef("ReadBlacklistItem(%s) (nil, %s)", h, err)
		return nil, err
	}

	instance := &BlacklistItem{
		id:         id,
		Hostname:   hostname,
		CreatedAt:   createdAt,
	}

	logger.Tracef("ReadBlacklistItem(%s) (%v, nil)", h, &instance)
	return instance, nil
}
