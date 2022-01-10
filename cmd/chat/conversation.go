package chat

import (
	"fmt"
	"time"

	"github.com/jroimartin/gocui"
)

type ChatView struct {
	name string
	dim ViewDimensionHandler
	g *gocui.Gui
	chatManager *ChatManager
}

func NewChatView(name string, dim ViewDimensionHandler) *ChatView {
	return &ChatView{name: name, dim: dim}
}

func (cv *ChatView) Layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	x, y, w, h := cv.dim(maxX, maxY)
	v, err := g.SetView(cv.name, x, y, x + w - 1, y + h - 1)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Autoscroll = true
		v.Wrap = true
		v.Title = cv.chatManager.Recipient()
		cv.g = g

		err = cv.chatManager.Start()
		if err != nil {
			return err
		}
	}
	return nil
}

func (cv *ChatView) AddMessage(msg string) {
	cv.g.Update(func(g *gocui.Gui) error {
		view, err := g.View(cv.name)
		if err != nil {
			return err
		}
		t := time.Now()
		tm := t.Format("15:04:05")
		fmt.Fprintf(view, "[%s] %s\n", tm, msg)
		return nil
	})
}

func (cv *ChatView) ScrollUp() {
	view, _ := cv.g.View(cv.name)
	ox, oy := view.Origin()
	if oy > 0 {
		view.Autoscroll = false
		view.SetOrigin(ox, oy - 1)
	}
}

func (cv *ChatView) ScrollDown() {
	view, _ := cv.g.View(cv.name)
	_, h := view.Size()
	ox, oy := view.Origin()
	lines := view.BufferLines()
	linesNum := len(lines)
	if linesNum >= oy + h {
		nOy := oy + 1
		if linesNum > nOy + h {
			view.Autoscroll = false
			view.SetOrigin(ox, nOy)
		} else {
			view.Autoscroll = true
		}
	}
}

func (cv *ChatView) SetChatManager(chatManager *ChatManager) {
	cv.chatManager = chatManager
}



