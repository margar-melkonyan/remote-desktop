package storage

import "database/sql"

type StorageConnector interface {
	NewConnection() (*sql.DB, error)
}
