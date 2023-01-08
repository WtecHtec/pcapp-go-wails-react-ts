package main

import (
	"context"
	"fmt"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

func (a *App) shutdown(ctx context.Context) {
	WxBot.Logout()
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	fmt.Println("Greet===")
	WxLogin()
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

func (a *App) LoginCode() {
	WxLogin()
}

// 获取 二维码路径
func (a *App) GetQRpath() string {
	return QR_PATH
}

func (a *App) Logger(msg interface{}) {
	fmt.Println("Logger =====")
	fmt.Println(msg)
}
