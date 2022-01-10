package chat

import (
	"github.com/jroimartin/gocui"
)

type InputView struct {
	name string
	dim ViewDimensionHandler
	maxLength int
	g *gocui.Gui
	chatManager *ChatManager
}

func NewInputView(name string, dim ViewDimensionHandler, maxLength int) *InputView {
	return &InputView{name: name, dim: dim, maxLength: maxLength}
}

func (cv *InputView) Layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	x, y, w, h := cv.dim(maxX, maxY)
	v, err := g.SetView(cv.name, x, y, x + w - 1, y + h - 1)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Editor = cv
		v.Editable = true
		cv.g = g
	}
	return nil
}

func (cv *InputView) Edit(v *gocui.View, key gocui.Key, ch rune, mod gocui.Modifier) {
	cx, _ := v.Cursor()
	ox, _ := v.Origin()
	limit := ox + cx + 1 > cv.maxLength
	switch {
	case ch != 0 && mod == 0 && !limit:
		v.EditWrite(ch)
	case key == gocui.KeyEnter:
		msg := v.Buffer()
		if msg != "" {
			cv.chatManager.SendMessage(msg)
		}
		v.SetOrigin(0, 0)
		v.SetCursor(0, 0)
		v.Clear()
	case key == gocui.KeyPgup:
		cv.chatManager.ScrollUp()
	case key == gocui.KeyPgdn:
		cv.chatManager.ScrollDown()			
	case key == gocui.KeySpace && !limit:
		v.EditWrite(' ')
	case key == gocui.KeyBackspace || key == gocui.KeyBackspace2:
		v.EditDelete(true)
	}
}

func (cv *InputView) SetChatManager(chatManager *ChatManager) {
	cv.chatManager = chatManager
}