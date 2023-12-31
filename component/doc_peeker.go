package component

import (
	"bytes"
	"context"
	"encoding/json"

	"github.com/kopecmaciej/mongui/config"
	"github.com/kopecmaciej/mongui/manager"
	"github.com/kopecmaciej/mongui/mongo"
	"github.com/kopecmaciej/mongui/primitives"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/rs/zerolog/log"
)

const (
	TextPeekerComponent manager.Component = "TextPeeker"
)

type peekerState struct {
	mongo.CollectionState
	rawDocument string
	id          primitive.ObjectID
}

type DocPeeker struct {
	*primitives.ModalView

	app         *App
	style       *config.DocPeeker
	eventChan   chan interface{}
	docModifier *DocModifier
	dao         *mongo.Dao
	state       peekerState
	manager     *manager.ComponentManager
}

func NewDocPeeker(dao *mongo.Dao) *DocPeeker {
	return &DocPeeker{
		ModalView:   primitives.NewModalView(),
		docModifier: NewDocModifier(dao),
		dao:         dao,
	}
}

func (dc *DocPeeker) Init(ctx context.Context) error {
	app, err := GetApp(ctx)
	if err != nil {
		return err
	}
	dc.app = app

	dc.setStyle()
	dc.setShortcuts(ctx)

	dc.manager = dc.app.ComponentManager

	if err := dc.docModifier.Init(ctx); err != nil {
		return err
	}

	return nil
}

func (dc *DocPeeker) setStyle() {
	dc.style = &dc.app.Styles.DocPeeker
	dc.SetBorder(true)
	dc.SetTitle("Document Details")
	dc.SetTitleAlign(tview.AlignLeft)
	dc.SetTitleColor(dc.style.TitleColor.Color())

	dc.ModalView.AddButtons([]string{"Edit", "Close"})
}

func (dc *DocPeeker) setShortcuts(ctx context.Context) {
	dc.ModalView.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyCtrlR {
			if err := dc.render(ctx); err != nil {
				log.Error().Err(err).Msg("Error refreshing document")
			}
			return nil
		}
		return event
	})
}

func (dc *DocPeeker) Peek(ctx context.Context, db, coll string, jsonString string) error {
	dc.state = peekerState{
		CollectionState: mongo.CollectionState{
			Db:   db,
			Coll: coll,
		},
		rawDocument: jsonString,
	}
	var prettyJson bytes.Buffer
	err := json.Indent(&prettyJson, []byte(jsonString), "", "  ")
	if err != nil {
		log.Printf("Error marshaling JSON: %v", err)
		return nil
	}
	text := string(prettyJson.Bytes())

	dc.ModalView.SetText(primitives.Text{
		Content: text,
		Color:   dc.style.ValueColor.Color(),
		Align:   tview.AlignLeft,
	})

	root := dc.app.Root
	root.AddPage(TextPeekerComponent, dc.ModalView, true, true)
	dc.ModalView.SetDoneFunc(func(buttonIndex int, buttonLabel string) {
		if buttonLabel == "Edit" {
			updatedDoc, err := dc.docModifier.Edit(ctx, db, coll, jsonString)
			if err != nil {
				log.Error().Err(err)
				return
			}
			dc.state.rawDocument = updatedDoc
			dc.render(ctx)
		} else if buttonLabel == "Close" || buttonLabel == "" {
			root.RemovePage(TextPeekerComponent)
		}
	})
	return nil
}

func (dc *DocPeeker) render(ctx context.Context) error {
	return dc.Peek(ctx, dc.state.Db, dc.state.Coll, dc.state.rawDocument)
}
