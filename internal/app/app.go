package app

import (
	"context"
	core "github.com/core-go/core"
	v "github.com/core-go/core/v10"
	"github.com/core-go/health"
	"github.com/core-go/log"
	"github.com/core-go/search/convert"
	q "github.com/core-go/sql"
	"github.com/core-go/sql/template"
	"github.com/lib/pq"
	"reflect"

	"go-service/internal/user"
)

type ApplicationContext struct {
	Health *health.Handler
	User   user.UserHandler
}

func NewApp(ctx context.Context, conf Config) (*ApplicationContext, error) {
	db, err := q.OpenByConfig(conf.Sql)
	if err != nil {
		return nil, err
	}
	logError := log.LogError
	status := core.InitializeStatus(conf.Status)
	action := core.InitializeAction(conf.Action)
	validator := v.NewValidator()

	buildParam := q.GetBuild(db)
	templates, err := template.LoadTemplates(template.Trim, "configs/query.xml")
	if err != nil {
		return nil, err
	}

	userType := reflect.TypeOf(user.User{})
	queryUser, err := template.UseQueryWithArray(conf.Template, user.BuildQuery, "user", templates, &userType, convert.ToMap, buildParam, pq.Array)
	userSearchBuilder, err := q.NewSearchBuilderWithArray(db, userType, queryUser, pq.Array)
	if err != nil {
		return nil, err
	}
	userRepository, err := q.NewRepositoryWithArray(db, "users", userType, pq.Array)
	if err != nil {
		return nil, err
	}
	userService := user.NewUserService(userRepository)
	userHandler := user.NewUserHandler(userSearchBuilder.Search, userService, status, logError, validator.Validate, &action)

	sqlChecker := q.NewHealthChecker(db)
	healthHandler := health.NewHandler(sqlChecker)

	return &ApplicationContext{
		Health: healthHandler,
		User:   userHandler,
	}, nil
}
