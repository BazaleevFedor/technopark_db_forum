package thread

import "github.com/BazaleevFedor/technopark_db_forum/internal/models"

type Repo interface {
	Create(forum *models.Thread) (*models.Thread, error)
	GetBySlugOrId(slug string, id int) (*models.Thread, error)
	Vote(vote *models.Vote) (*models.Thread, error)
	UpdateThread(thread *models.Thread) (*models.Thread, error)
}
