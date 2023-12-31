package forum

import (
	"github.com/BazaleevFedor/technopark_db_forum/internal/models"
)

type Repo interface {
	Create(threadSlug string, threadId int, posts []models.Post) ([]models.Post, error)
	GetThreadPosts(threadSlug string, threadId int, desc bool, limit int, since int, sort string) ([]models.Post, error)
	CheckThreadBySlugOrId(slug string, id int) (bool, error)
	GetPostByIdRelated(id int, related []string) (*models.Post, *models.User, *models.Forum, *models.Thread, error)
	UpdatePost(post *models.Post) (*models.Post, error)
}
