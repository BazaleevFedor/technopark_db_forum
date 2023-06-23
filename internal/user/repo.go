package user

import (
	"github.com/BazaleevFedor/technopark_db_forum/internal/models"
)

type Repo interface {
	Create(user *models.User) (*models.User, error)
	GetByEmailOrNick(user *models.User) ([]models.User, error)
	GetByNick(nick string) (*models.User, error)
	GetByEmail(email string) (string, error)
	Update(user *models.User) (*models.User, error)
}
