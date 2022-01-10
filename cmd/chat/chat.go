package chat

import (
	"fmt"

	"github.com/akamensky/argparse"
	"github.com/jroimartin/gocui"

	"github.com/marcinbor85/nes/cmd"
)

type ChatContext struct {
	To 			*string
}

var Cmd = &cmd.Command{
	Name: "chat",
	Help: "Interactive chat with recipient",
	Context: &ChatContext{},
	Init: Init,
	Logic: Logic,
}

const (
	MAX_INPUT_LENGTH = 1024
)

type ViewDimensionHandler func(maxX, maxY int) (x, y, w, h int)

func Init(c *cmd.Command) {
	ctx := c.Context.(*ChatContext)
	ctx.To = c.Cmd.String("t", "to", &argparse.Options{Required: true, Help: "Recipient username"})
}

func SetFocus(name string) func(g *gocui.Gui) error {
	return func(g *gocui.Gui) error {
		_, err := g.SetCurrentView(name)
		return err
	}
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func Logic(c *cmd.Command) {
	ctx := c.Context.(*ChatContext)

	recipient := *ctx.To
	
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		fmt.Println("Cannot open terminal:", err.Error())
		return
	}
	defer g.Close()

	g.Cursor = true

	chat := NewChatView("chat", func(maxX, maxY int) (x, y, w, h int) {
		return 0, 0, maxX, maxY - 3
	})
	input := NewInputView("input", func(maxX, maxY int) (x, y, w, h int) {
		return 0, maxY - 3, maxX, 3
	}, MAX_INPUT_LENGTH)
	
	manager, err := NewChatManager(chat, input, recipient)
	if err != nil {
		g.Close()
		fmt.Println("Cannot start chat:", err.Error())
		return
	}
	
	focus := gocui.ManagerFunc(SetFocus("input"))
	g.SetManager(chat, input, focus)

	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		g.Close()
		manager.Disconnect()
		fmt.Println("Cannot set key bindings:", err.Error())
		return
	}
	
	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		g.Close()
		manager.Disconnect()
		fmt.Println("Error:", err.Error())
		return
	}
}
