package component

import (
	"context"
	"mongo-ui/mongo"

	"github.com/rivo/tview"
)

const (
	appCtxKey = "app"
)

type App struct {
	*tview.Application

	Root *Root
}

func NewApp() App {
	client := mongo.NewClient()
	client.Connect()
	mongoDao := mongo.NewDao(client.Client, client.Config)

	app := App{
		Application: tview.NewApplication(),
		Root:        NewRoot(mongoDao),
	}

	return app
}

func (a *App) Init() error {
	ctx := LoadApp(context.Background(), a)
	err := a.Root.Init(ctx)
	if err != nil {
		return err
	}
	focus := a.GetFocus()
	a.SetRoot(a.Root.Pages, true).EnableMouse(true)
	a.SetFocus(focus)
	return a.Run()
}

func GetApp(ctx context.Context) *App {
	app, ok := ctx.Value(appCtxKey).(*App)
	if !ok {
		panic("App not found in context")
	}
	return app
}

func LoadApp(ctx context.Context, app *App) context.Context {
	return context.WithValue(ctx, appCtxKey, app)
}
