package models

import (
	"time"
)

type FileInfo struct {
	Name         string
	CreatedAt    time.Time
	LastModified time.Time
}
