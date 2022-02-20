package app

import (
	"context"
	"github.com/core-go/health"
	"github.com/core-go/log"
	mgo "github.com/core-go/mongo"
	"github.com/core-go/search"
	mq "github.com/core-go/search/mongo"
	sv "github.com/core-go/service"
	v "github.com/core-go/service/v10"
	"reflect"

	. "go-service/internal/usecase/user"
)

type ApplicationContext struct {
	HealthHandler *health.Handler
	UserHandler   UserHandler
}

func NewApp(ctx context.Context, root Root) (*ApplicationContext, error) {
	db, err := mgo.Setup(ctx, root.Mongo)
	if err != nil {
		return nil, err
	}
	logError := log.ErrorMsg
	status := sv.InitializeStatus(root.Status)
	action := sv.InitializeAction(root.Action)
	validator := v.NewValidator()

	userType := reflect.TypeOf(User{})
	userQueryBuilder := mq.NewBuilder(userType)
	userSearchBuilder := mgo.NewSearchBuilder(db, "users", userQueryBuilder.BuildQuery, search.GetSort)
	userRepository := mgo.NewRepository(db, "users", userType)

	userService := NewUserService(userRepository)
	userHandler := NewUserHandler(userSearchBuilder.Search, userService, status, logError, validator.Validate, &action)

	mongoChecker := mgo.NewHealthChecker(db)
	healthHandler := health.NewHandler(mongoChecker)

	return &ApplicationContext{
		HealthHandler: healthHandler,
		UserHandler:   userHandler,
	}, nil
}
