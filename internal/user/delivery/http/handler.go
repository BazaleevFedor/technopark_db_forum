package handler

import (
	"net/http"

	"github.com/BazaleevFedor/technopark_db_forum/internal/models"
	"github.com/BazaleevFedor/technopark_db_forum/internal/tools/errors"
	"github.com/BazaleevFedor/technopark_db_forum/internal/user"
	"github.com/jackc/pgx"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	Repo user.Repo
}

const (
	NickCtxKey = "username"
)

func NewHandler(repo user.Repo) *Handler {
	return &Handler{Repo: repo}
}
func (h *Handler) CreateUser(ctx echo.Context) error {
	var newUserReq models.User
	if err := ctx.Bind(&newUserReq); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, errors.BAD_BODY)
	}
	newUserReq.Nick = ctx.Param(NickCtxKey)
	newUserResp, err := h.Repo.Create(&newUserReq)
	if err != nil {
		conflictUsers, err := h.Repo.GetByEmailOrNick(&newUserReq)
		if err != nil || len(conflictUsers) == 0 {
			return echo.NewHTTPError(http.StatusInternalServerError, errors.INTERNAL_SERVER_ERROR)
		}
		return ctx.JSON(http.StatusConflict, conflictUsers)
	}
	return ctx.JSON(http.StatusCreated, newUserResp)
}

func (h *Handler) GetUser(ctx echo.Context) error {
	nick := ctx.Param(NickCtxKey)
	userResp, err := h.Repo.GetByNick(nick)
	if err != nil {
		if err == pgx.ErrNoRows {
			return echo.NewHTTPError(http.StatusNotFound, errors.NOT_FOUND_USER_BY_NICK+nick)
		}
		return echo.NewHTTPError(http.StatusInternalServerError, errors.INTERNAL_SERVER_ERROR)
	}
	return ctx.JSON(http.StatusOK, userResp)
}

func (h *Handler) UpdateUser(ctx echo.Context) error {
	var updateUserReq models.User
	if err := ctx.Bind(&updateUserReq); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, errors.BAD_BODY)
	}
	updateUserReq.Nick = ctx.Param(NickCtxKey)
	newUserResp, err := h.Repo.Update(&updateUserReq)
	if err != nil {
		if err == pgx.ErrNoRows {
			return echo.NewHTTPError(http.StatusNotFound, errors.NOT_FOUND_USER_BY_NICK+updateUserReq.Nick)
		}
		conflictUser, err := h.Repo.GetByEmail(updateUserReq.Email)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, errors.INTERNAL_SERVER_ERROR)
		}
		return echo.NewHTTPError(http.StatusConflict, errors.EMAIL_ALREADY_IN_USE+conflictUser)
	}
	return ctx.JSON(http.StatusOK, newUserResp)
}
