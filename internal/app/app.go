package app

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"
)

type BeforeFunc func(context.Context) error
type ActionFunc func(context.Context) error
type AfterFunc func(context.Context) error

type App struct {
	Ctx    context.Context
	Before BeforeFunc
	Action ActionFunc
	After  AfterFunc
	stop   chan struct{}
	Name   string
}

func (a *App) Run() error {
	a.Ctx = context.Background()
	if a.Before != nil {
		if err := a.Before(a.Ctx); err != nil {
			return err
		}
	}
	a.listener()
	if a.Action != nil {
		return a.Action(a.Ctx)
	}
	return nil
}
func (a *App) listener() {
	a.stop = make(chan struct{})
	go func() {
		sig := make(chan os.Signal)
		signal.Notify(sig, os.Interrupt, syscall.SIGINT, syscall.SIGKILL, syscall.SIGTERM)
		<-sig
		ctx, cancel := context.WithTimeout(a.Ctx, 8*time.Second)
		defer cancel()
		if a.After != nil {
			if err := a.After(ctx); err != nil {
				_, _ = fmt.Fprintln(os.Stderr, err)
			}
		}
		close(a.stop)
	}()
}
func (a *App) Wait() {
	<-a.stop
}

func NewApp() *App {
	return &App{
		Name: filepath.Base(os.Args[0]),
	}
}
