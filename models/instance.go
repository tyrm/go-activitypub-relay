package models

import (
	"database/sql"
	"time"
)

type Instance struct {
	Hostname string

	JoinedAt   time.Time
	ApprovedAt sql.NullTime

	// metadata
	id int
}

const sqlApproveInstance = `
UPDATE "public"."instances"
SET approved_at = current_timestamp
WHERE id = $1
RETURNING approved_at;`
func (i *Instance) Approve() error {

	var approvedAt time.Time

	err := db.QueryRow(sqlApproveInstance, i.id).Scan(&approvedAt)
	if err != nil {
		logger.Tracef("(%v) Approve() (%s)", &i, err)
		return err
	}

	logger.Tracef("(%v) Approve() (nil)", &i)
	return nil
}


const sqlCreateInstance = `
INSERT INTO "public"."instances" (hostname)
VALUES ($1)
RETURNING id, joined_at;`

func CreateInstance(h string) (*Instance, error) {
	var id int


	var joinedAt time.Time

	err := db.QueryRow(sqlCreateInstance, h).Scan(&id, &joinedAt)
	if err != nil {
		logger.Tracef("CreateInstance(%s) (nil, %s)", h, err)
		return nil, err
	}

	instance := &Instance{
		id:         id,
		Hostname:   h,
		JoinedAt:   joinedAt,
	}

	logger.Tracef("CreateInstance(%s) (%v, nil)", h, &instance)
	return instance, nil
}


const sqlGetApprovedInstances = `
SELECT id, hostname, joined_at, approved_at
FROM instances
WHERE approved_at IS NOT NULL;`

func GetApprovedInstances() (*[]Instance, error) {
	rows, err := db.Query(sqlGetApprovedInstances)
	if err != nil {
		logger.Tracef("GetApprovedInstances() (nil, %s)", err)
		return nil, err
	}
	defer rows.Close()

	var instanceList []Instance
	for rows.Next() {
		var id int
		var hostname string
		var joinedAt time.Time
		var approvedAt sql.NullTime

		if err := rows.Scan(&id, &hostname, &joinedAt, &approvedAt); err != nil {
			logger.Tracef("GetApprovedInstances() (nil, %s)", err)
			return nil, err
		}

		instance := Instance{
			id:         id,
			Hostname:   hostname,
			JoinedAt:   joinedAt,
			ApprovedAt: approvedAt,
		}
		instanceList = append(instanceList, instance)
	}

	logger.Tracef("GetApprovedInstances() (%d, nil)", len(instanceList))
	return &instanceList, nil
}

const sqlReadInstance = `
SELECT id, hostname, joined_at, approved_at
FROM instances
WHERE hostname = $1;`

func ReadInstance(h string) (*Instance, error) {
	var id int

	var hostname string

	var joinedAt time.Time
	var approvedAt sql.NullTime

	err := db.QueryRow(sqlReadInstance, h).Scan(&id, &hostname, &joinedAt, &approvedAt)
	if err != nil {
		logger.Tracef("ReadInstance(%s) (nil, %s)", h, err)
		return nil, err
	}

	instance := &Instance{
		id:         id,
		Hostname:   hostname,
		JoinedAt:   joinedAt,
		ApprovedAt: approvedAt,
	}

	logger.Tracef("ReadInstance(%s) (%v, nil)", h, &instance)
	return instance, nil
}
