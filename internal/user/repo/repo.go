package userRepo

import (
	"github.com/BazaleevFedor/technopark_db_forum/internal/models"
	"github.com/jackc/pgx"
)

type Repo struct {
	Conn *pgx.ConnPool
}

func NewRepo(conn *pgx.ConnPool) *Repo {
	conn.Prepare("create_user", "INSERT into users(name, nick, email, about) VALUES ($1,$2,$3,$4)")
	conn.Prepare("update_user", "UPDATE users SET name=COALESCE(NULLIF($1, ''), name), email=COALESCE(NULLIF($2, ''), email), about=COALESCE(NULLIF($3, ''), about) WHERE nick = $4 RETURNING name,nick,email,about")
	conn.Prepare("get_user_by_email_or_nick", "SELECT name,nick,email,about FROM users WHERE nick=$1 OR email=$2")
	conn.Prepare("get_user_by_nick", "SELECT name, nick, email, about FROM users WHERE nick=$1")
	conn.Prepare("get_user_by_email", "SELECT nick FROM users WHERE email=$1")

	return &Repo{Conn: conn}
}
func (r *Repo) Create(user *models.User) (*models.User, error) {
	_, err := r.Conn.Exec(`EXECUTE create_user($1,$2,$3,$4)`, user.Name, user.Nick, user.Email, user.About)
	if err != nil {
		return nil, err
	}
	return user, nil
}
func (r *Repo) Update(user *models.User) (*models.User, error) {
	err := r.Conn.QueryRow("EXECUTE update_user($1,$2,$3,$4)", user.Name, user.Email, user.About, user.Nick).Scan(&user.Name, &user.Nick, &user.Email, &user.About)
	if err != nil {
		return nil, err
	}
	return user, nil
}
func (r *Repo) GetByEmailOrNick(user *models.User) ([]models.User, error) {
	userResp := make([]models.User, 0, 2)
	userRows, err := r.Conn.Query(`EXECUTE get_user_by_email_or_nick($1,$2)`, user.Nick, user.Email)
	if err != nil {
		return nil, err
	}
	defer userRows.Close()
	for userRows.Next() {
		user := models.User{}
		err = userRows.Scan(&user.Name, &user.Nick, &user.Email, &user.About)
		if err != nil {
			return nil, err
		}
		userResp = append(userResp, user)
	}
	return userResp, nil
}
func (r *Repo) GetByNick(nick string) (*models.User, error) {
	user := &models.User{}
	err := r.Conn.QueryRow("EXECUTE get_user_by_nick($1)", nick).Scan(&user.Name, &user.Nick, &user.Email, &user.About)
	if err != nil {
		return nil, err
	}
	return user, nil
}
func (r *Repo) GetByEmail(email string) (string, error) {
	var userNick string
	err := r.Conn.QueryRow("EXECUTE get_user_by_email($1)", email).Scan(&userNick)
	if err != nil {
		return "", err
	}
	return userNick, nil
}
