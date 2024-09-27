package client

import (
	"context"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/kTowkA/GophKeeper/client/config"
	"github.com/kTowkA/GophKeeper/client/models"
	"github.com/kTowkA/GophKeeper/client/models/options"
	pb "github.com/kTowkA/GophKeeper/grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	waitTime = time.Second * 10
)

type Gophkeeper struct {
	token   string
	gClient pb.GophKeeperClient
}

func (gc *Gophkeeper) IsLogged() bool {
	if gc.token != "" {
		return true
	}
	return false
}

func Run(ctx context.Context, cfg config.Config) error {
	if cfg.WithLog() {
		fl, err := tea.LogToFile(cfg.LogFile(), "debug")
		if err != nil {
			return err
		}
		defer fl.Close()
	}

	conn, err := grpc.NewClient(cfg.Address(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}
	defer conn.Close()

	gClient := pb.NewGophKeeperClient(conn)

	p := tea.NewProgram(models.NewController(ctx, options.CreateOptions(gClient, options.SimpleCrypter{})))
	if _, err := p.Run(); err != nil {
		return err
	}

	return nil
}
