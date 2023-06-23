package service

import "github.com/BazaleevFedor/technopark_db_forum/internal/models"

type Repo interface {
	Status() (*models.Status, error)
	TruncateDB() error
}
