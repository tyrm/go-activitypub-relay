package activitypub

import (
	"database/sql"

	"github.com/tyrm/go-activitypub-relay/models"
)

func OnBlacklist(hostname string) (bool, error) {
	_, err := models.ReadBlacklistItem(hostname)
	if err == sql.ErrNoRows {
		// Instance not on list
		logger.Tracef("OnBlacklist(%s) (false, nil)", hostname)
		return false, nil
	} else if err != nil {
		// Return Error
		logger.Tracef("OnBlacklist(%s) (false, %s)", hostname, err)
		return false, err
	}

	// Instance Exists and is approved
	logger.Tracef("OnBlacklist(%s) (true, nil)", hostname)
	return true, nil
}