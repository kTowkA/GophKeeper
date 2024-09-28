package client

import (
	"context"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/kTowkA/GophKeeper/client/config"
	"github.com/kTowkA/GophKeeper/client/models"
	pb "github.com/kTowkA/GophKeeper/grpc"
	"github.com/kTowkA/GophKeeper/internal/crypto/simplecrypter"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Run запуск клиента
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

	p := tea.NewProgram(models.NewController(ctx, models.NewServices(gClient, simplecrypter.SimpleCrypter)))
	if _, err := p.Run(); err != nil {
		return err
	}

	return nil
}
