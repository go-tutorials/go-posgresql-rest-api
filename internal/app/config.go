package app

import (
	"github.com/core-go/core"
	"github.com/core-go/log"
	mid "github.com/core-go/log/middleware"
	"github.com/core-go/sql"
)

type Config struct {
	Server     core.ServerConf    `mapstructure:"server"`
	Template   bool               `mapstructure:"template"`
	Sql        sql.Config         `mapstructure:"sql"`
	Log        log.Config         `mapstructure:"log"`
	MiddleWare mid.LogConfig      `mapstructure:"middleware"`
	Status     *core.StatusConfig `mapstructure:"status"`
	Action     *core.ActionConfig `mapstructure:"action"`
}
