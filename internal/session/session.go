package session

import (
	"context"

	"github.com/jinzhu/copier"
	"github.com/lixvyang/betxin.one/config"
	"github.com/lixvyang/betxin.one/pkg/logger"
	"github.com/rs/zerolog"
)

type Session struct {
	Version string
	Logger  *zerolog.Logger
	Conf    *config.AppConfig
}

var contextKey struct{}

func With(ctx context.Context, s *Session) context.Context {
	return context.WithValue(ctx, contextKey, s)
}

func From(ctx context.Context) *Session {
	return ctx.Value(contextKey).(*Session)
}

func (s *Session) WithConf(conf *config.AppConfig) *Session {
	s.Conf = conf
	return s
}

func (s *Session) WithLogger(conf *config.LogConfig) *Session {
	logConf := new(logger.LogConfig)
	copier.Copy(logConf, conf)
	s.Logger = logger.New(logConf)
	return s
}
