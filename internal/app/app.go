package app

import (
	"context"
	"github.com/core-go/health"
	"github.com/core-go/log"
	mgo "github.com/core-go/mongo"
	mq "github.com/core-go/mongo/query"
	"github.com/core-go/search"
	sv "github.com/core-go/service"
	v "github.com/core-go/service/v10"
	"reflect"

	"go-service/internal/usecase/user"
)

type ApplicationContext struct {
	HealthHandler *health.Handler
	UserHandler   user.UserHandler
}

func NewApp(ctx context.Context, root Root) (*ApplicationContext, error) {
	db, err := mgo.Setup(ctx, root.Mongo)
	if err != nil {
		return nil, err
	}
	logError := log.ErrorMsg
	status := sv.InitializeStatus(root.Status)

	userType := reflect.TypeOf(user.User{})
	userQueryBuilder := mq.NewBuilder(userType)
	userSearchBuilder := mgo.NewSearchBuilder(db, "users", userQueryBuilder.BuildQuery, search.GetSort)
	userRepository := mgo.NewRepository(db, "users", userType)

	userService := user.NewUserService(userRepository)
	validator := v.NewValidator()
	userHandler := user.NewUserHandler(userSearchBuilder.Search, userService, status, validator.Validate, logError)

	mongoChecker := mgo.NewHealthChecker(db)
	healthHandler := health.NewHandler(mongoChecker)

	return &ApplicationContext{
		HealthHandler: healthHandler,
		UserHandler:   userHandler,
	}, nil
}
