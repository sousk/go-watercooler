package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	exitOK = iota
	exitError
)

var (
	version string
	versionPrinted = flag.Bool("version", false, "version")
)

func main() {
	os.Exit(RealMain(os.Args))
}

func RealMain(args []string) int {
	err := Execute(args)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "[ERROR] Failed: %s\n", err)
		return exitError
	}

	return exitOK
}

func Execute(args []string) error {
	flag.Parse()

	if *versionPrinted {
		fmt.Println(version)
		return nil
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	middlewares(ctx)

	sigChan := make(chan os.Signal, 1)

	signal.Notify(sigChan,
		syscall.SIGHUP,
		 syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	cmderr := make(chan error)
	go func() {
		cmderr <- command(ctx, args[1:])
	}()

	select {
	case err := <- cmderr:
		return err
	case sig := <-sigChan:
		fmt.Printf("signal--> %v\n", sig)
	}

	fmt.Println("existing")

	return nil
}

func command(ctx context.Context, args []string) error {
	fmt.Println("I'm a command")
	time.Sleep(5 * time.Second)
	// return fmt.Errorf("error happens")
	return nil
}

func middlewares(ctx context.Context) {
	fmt.Println("middlewares are working")
	go func() {
		select {
		case <- ctx.Done():
			fmt.Println("middlewares shutting down gracefully")
		}
	}()
}

