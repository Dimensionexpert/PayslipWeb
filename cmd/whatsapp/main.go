package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Dimensionexpert/payslip/internal/whatsapp"
	_ "modernc.org/sqlite"
)

func run() error {
	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	// Initialize the WhatsApp service.
	service, err := whatsapp.NewService(
		ctx,
		"file:whatsapp-session.db?_foreign_keys=on",
	)
	if err != nil {
		return fmt.Errorf("initialize WhatsApp: %w", err)
	}

	// Disconnect when run returns.
	defer service.Disconnect()

	fmt.Println("Starting WhatsApp pairing / connection...")

	if err := service.Pair(ctx); err != nil {
		if errors.Is(err, context.Canceled) {
			return nil
		}
		return fmt.Errorf("WhatsApp pairing/connection failed: %w", err)
	}

	fmt.Println("WhatsApp connected or pairing completed.")

	// Keep the process alive until Ctrl+C.
	<-ctx.Done()

	fmt.Println("\nShutting down WhatsApp...")
	return nil
}

func main() {
	if err := run(); err != nil {
		log.Printf("Error: %v", err)
		os.Exit(1)
	}
}
