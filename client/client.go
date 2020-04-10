package client

import (
	"bufio"
	"log"
	"net"
	"strings"
	"time"

	"github.com/marcusolsson/tui-go"
)

type Client struct {
	Nick   string
	Server string
}

func (c *Client) Run() error {
	// Connect to server
	conn, err := net.Dial("tcp", c.Server)
	if err != nil {
		return err
	}

	c.startUI(conn)
	return nil
}

func (c *Client) startUI(cnx net.Conn) {
	topbar := tui.NewVBox(
		tui.NewLabel("WELCOME on parrot ! Your first message will be your nickname."),
		tui.NewSpacer(),
		tui.NewLabel("To quit, press Esc"),
		tui.NewSpacer(),
		tui.NewLabel("Commands available: /date"),
		tui.NewSpacer(),
	)
	topbar.SetBorder(true)

	history := tui.NewVBox()

	historyScroll := tui.NewScrollArea(history)
	historyScroll.SetAutoscrollToBottom(true)

	historyBox := tui.NewVBox(historyScroll)
	historyBox.SetBorder(true)

	input := tui.NewEntry()
	input.SetFocused(true)
	input.SetSizePolicy(tui.Expanding, tui.Maximum)

	inputBox := tui.NewHBox(input)
	inputBox.SetBorder(true)
	inputBox.SetSizePolicy(tui.Expanding, tui.Maximum)

	chat := tui.NewVBox(historyBox, inputBox)
	chat.SetSizePolicy(tui.Expanding, tui.Expanding)

	input.OnSubmit(func(e *tui.Entry) {
		writer := bufio.NewWriter(cnx)
		writer.WriteString(e.Text() + "\n")
		writer.Flush()

		input.SetText("")
	})

	root := tui.NewVBox(topbar, chat)

	ui, err := tui.New(root)
	if err != nil {
		log.Fatal(err)
	}

	ui.SetKeybinding("Esc", func() { ui.Quit() })
	ui.SetKeybinding("Ctrl+c", func() { ui.Quit() })

	go func() {
		for {
			r := bufio.NewReader(cnx)
			message, err := r.ReadString('\n')

			if err != nil {
				log.Print(err)
			}
			message = strings.TrimSpace(message)

			ui.Update(func() {
				history.Append(tui.NewHBox(
					tui.NewLabel(time.Now().Format("15:04")),
					tui.NewPadder(1, 0, tui.NewLabel(message)),
					tui.NewSpacer(),
				))
			})
		}
	}()

	if err := ui.Run(); err != nil {
		log.Fatal(err)
	}
}
