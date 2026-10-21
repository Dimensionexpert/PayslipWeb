package whatsapp

import (
	"context"
	"fmt"
	"os"

	"github.com/mdp/qrterminal/v3"
)

func (s *Service) Pair(ctx context.Context) error {
	if s.client == nil {
		return fmt.Errorf("WhatsApp client is not initialized")
	}

	// Existing device: reconnect using the saved session.
	if s.client.Store.ID != nil {
		if err := s.client.Connect(); err != nil {
			return fmt.Errorf("connect to WhatsApp using saved session: %w", err)
		}
		return nil
	}

	// No saved device: prepare the QR pairing flow.
	qrChan, err := s.client.GetQRChannel(ctx)
	if err != nil {
		return fmt.Errorf("get WhatsApp QR channel: %w", err)
	}

	if err := s.client.Connect(); err != nil {
		return fmt.Errorf("connect to WhatsApp for pairing: %w", err)
	}

	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("WhatsApp pairing cancelled: %w", ctx.Err())

		case evt, ok := <-qrChan:
			if !ok {
				return fmt.Errorf("WhatsApp QR channel closed before pairing completed")
			}

			switch evt.Event {
			case "code":
				qrterminal.GenerateHalfBlock(
					evt.Code,
					qrterminal.L,
					os.Stdout,
				)
			case "success":
				fmt.Println("WhatsApp pairing successful")
				return nil

			case "timeout":
				return fmt.Errorf("WhatsApp QR code expired; try pairing again")

			default:
				fmt.Printf("WhatsApp login event: %s\n", evt.Event)
			}
		}
	}
}
