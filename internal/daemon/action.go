package daemon

import (
	"context"
	"time"

	"github.com/kardianos/service"
	"github.com/kitabisa/mubeng/internal/server"
)

func (p *program) Start(s service.Service) error {
	go server.Run(p.opt)
	return nil
}

func (p *program) Stop(s service.Service) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	server.Stop(ctx)
	return nil
}
