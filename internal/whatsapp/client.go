package whatsapp

import (
	"context"
	"fmt"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store/sqlstore"
	waLog "go.mau.fi/whatsmeow/util/log"
)

type Service struct {
	client    *whatsmeow.Client
	sessionDB *sqlstore.Container
}

func (s *Service) Disconnect() {
	if s.client != nil {
		s.client.Disconnect()
	}
}

func NewService(ctx context.Context, dbPath string) (*Service, error) {
	// Creating the Session logger
	sessionLog := waLog.Stdout("Database", "INFO", true)
	// Creating the Session db.
	sessionDB, err := sqlstore.New(
		ctx,
		"sqlite",
		dbPath,
		sessionLog,
	)
	if err != nil {
		return nil, fmt.Errorf("initialize WhatsApp session database: %w", err)
	}
	// Getting the device(client).
	// this is meant for single user

	clientStore, err := sessionDB.GetFirstDevice(ctx)
	if err != nil {
		return nil, fmt.Errorf("getting the client from session database: %w", err)
	}
	// Creating the WhatsApp client logger.
	clientLog := waLog.Stdout("Client", "INFO", true)

	// Creating the WhatsApp client.
	client := whatsmeow.NewClient(clientStore, clientLog)

	return &Service{
		client:    client,
		sessionDB: sessionDB,
	}, nil
}
