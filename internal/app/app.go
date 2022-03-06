package app

import (
	"context"
	"github.com/core-go/health"
	"github.com/core-go/log"
	"github.com/core-go/search/convert"
	sv "github.com/core-go/service"
	v "github.com/core-go/service/v10"
	q "github.com/core-go/sql"
	"github.com/core-go/sql/template"
	"github.com/lib/pq"
	"reflect"

	. "go-service/internal/usecase/user"
)

type ApplicationContext struct {
	HealthHandler *health.Handler
	UserHandler   UserHandler
}

func NewApp(ctx context.Context, conf Config) (*ApplicationContext, error) {
	db, err := q.OpenByConfig(conf.Sql)
	if err != nil {
		return nil, err
	}
	logError := log.ErrorMsg
	status := sv.InitializeStatus(conf.Status)
	action := sv.InitializeAction(conf.Action)
	validator := v.NewValidator()

	buildParam := q.GetBuild(db)
	templates, err := template.LoadTemplates(template.Trim, "configs/query.xml")
	if err != nil {
		return nil, err
	}

	userType := reflect.TypeOf(User{})
	//userQueryBuilder := query.NewBuilder(db, "users", userType)
	queryUser, err := template.UseQueryWithArray(conf.Template, BuildQuery, "user", templates, &userType, convert.ToMap, buildParam, pq.Array)
	userSearchBuilder, err := q.NewSearchBuilderWithArray(db, userType, queryUser, pq.Array)
	if err != nil {
		return nil, err
	}
	userRepository, err := q.NewRepositoryWithArray(db, "users", userType, pq.Array)
	if err != nil {
		return nil, err
	}
	userService := NewUserService(userRepository)
	userHandler := NewUserHandler(userSearchBuilder.Search, userService, status, logError, validator.Validate, &action)

	sqlChecker := q.NewHealthChecker(db)
	healthHandler := health.NewHandler(sqlChecker)

	return &ApplicationContext{
		HealthHandler: healthHandler,
		UserHandler:   userHandler,
	}, nil
}
