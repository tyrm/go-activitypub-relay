package activitypub

import (
	"database/sql"

	"github.com/tyrm/go-activitypub-relay/models"
)

func InstanceExists(hostname string) (bool, error) {
	_, err := models.ReadInstance(hostname)
	if err == sql.ErrNoRows {
		// Instance not on list
		logger.Tracef("InstanceExists(%s) (false, nil)", hostname)
		return false, nil
	} else if err != nil {
		// Return Error
		logger.Tracef("InstanceExists(%s) (false, %s)", hostname, err)
		return false, err
	}

	// Instance Exists and is approved
	logger.Tracef("InstanceExists(%s) (true, nil)", hostname)
	return true, nil
}

func ApprovedInstanceExists(hostname string) (bool, error) {
	instance, err := models.ReadInstance(hostname)
	if err == sql.ErrNoRows {
		// Instance not on list
		logger.Tracef("ApprovedInstanceExists(%s) (false, nil)", hostname)
		return false, nil
	} else if err != nil {
		// Return Error
		logger.Tracef("ApprovedInstanceExists(%s) (false, %s)", hostname, err)
		return false, err
	}

	nullTime := sql.NullTime{Valid: false}
	if instance.ApprovedAt == nullTime {
		// Instance exists but isn't approved
		logger.Tracef("ApprovedInstanceExists(%s) (false, nil)", hostname)
		return false, nil
	}

	// Instance Exists and is approved
	logger.Tracef("ApprovedInstanceExists(%s) (true, nil)", hostname)
	return true, nil
}
