package main

import (
	"fmt"
	"log"

	"github.com/BazaleevFedor/technopark_db_forum/config"
	"github.com/BazaleevFedor/technopark_db_forum/configRouting"
	forumDelivery "github.com/BazaleevFedor/technopark_db_forum/internal/forum/delivery/http"
	forumRepository "github.com/BazaleevFedor/technopark_db_forum/internal/forum/repo"
	postDelivery "github.com/BazaleevFedor/technopark_db_forum/internal/post/delivery/http"
	postRepository "github.com/BazaleevFedor/technopark_db_forum/internal/post/repo"
	serviceDelivery "github.com/BazaleevFedor/technopark_db_forum/internal/service/delivery/http"
	serviceRepository "github.com/BazaleevFedor/technopark_db_forum/internal/service/repo"
	threadDelivery "github.com/BazaleevFedor/technopark_db_forum/internal/thread/delivery/http"
	threadRepository "github.com/BazaleevFedor/technopark_db_forum/internal/thread/repo"
	userDelivery "github.com/BazaleevFedor/technopark_db_forum/internal/user/delivery/http"
	userRepository "github.com/BazaleevFedor/technopark_db_forum/internal/user/repo"
	"github.com/jackc/pgx"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/labstack/echo-contrib/pprof"
)

func main() {
	connStr := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable port=%s",
		config.DbConfig.User, config.DbConfig.Password, config.DbConfig.DBName, config.DbConfig.Port)
	pgxConn, err := pgx.ParseConnectionString(connStr)
	if err != nil {
		log.Fatal(err.Error())
	}
	pgxConn.PreferSimpleProtocol = true
	config := pgx.ConnPoolConfig{
		ConnConfig:     pgxConn,
		MaxConnections: config.DbConfig.MaxConnections,
		AfterConnect:   nil,
		AcquireTimeout: 0,
	}
	connPool, err := pgx.NewConnPool(config)

	if err != nil {
		log.Fatal(err.Error())
	}
	e := echo.New()
	pprof.Register(e)
	e.Use(middleware.Recover())
	// e.Use(middleware.Logger())
	userRepo := userRepository.NewRepo(connPool)
	userHandler := userDelivery.NewHandler(userRepo)
	forumRepo := forumRepository.NewRepo(connPool)
	forumHandler := forumDelivery.NewHandler(forumRepo)
	threadRepo := threadRepository.NewRepo(connPool)
	threadHandler := threadDelivery.NewHandler(threadRepo)
	postRepo := postRepository.NewRepo(connPool)
	postHandler := postDelivery.NewHandler(postRepo)
	servRepo := serviceRepository.NewRepo(connPool)
	servHandler := serviceDelivery.NewHandler(servRepo)

	handlers := configRouting.Handlers{
		UserHandler:    userHandler,
		ForumHandler:   forumHandler,
		ThreadHandler:  threadHandler,
		PostHandler:    postHandler,
		ServiceHandler: servHandler,
	}
	handlers.ConfigureRouting(e)

	e.Logger.Fatal(e.Start(":5000"))
}
