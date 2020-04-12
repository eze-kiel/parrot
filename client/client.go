package client

import (
	"bufio"
	"net"
	"os"
	"strings"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
	"github.com/gobuffalo/packr"
	"github.com/marcusolsson/tui-go"
	log "github.com/sirupsen/logrus"
)

type Client struct {
	Nick   string
	Server string
}

func (c *Client) Run(sound bool) error {
	// Connect to server
	conn, err := net.Dial("tcp", c.Server)
	if err != nil {
		log.Fatalf("Error: %s unreachable\n", c.Server)
		return err
	}

	w := bufio.NewWriter(conn)
	w.WriteString(c.Nick + "\n")
	w.Flush()

	c.startUI(conn, sound)
	return nil
}

func (c *Client) startUI(cnx net.Conn, sound bool) {
	topbar := tui.NewVBox(
		tui.NewLabel("To quit, press Esc"),
		tui.NewSpacer(),
		tui.NewLabel("Commands available: /date"),
		tui.NewSpacer(),
	)
	topbar.SetBorder(true)
	topbar.SetTitle(" Welcome on Parrot Chat ")

	history := tui.NewVBox()

	historyScroll := tui.NewScrollArea(history)
	historyScroll.SetAutoscrollToBottom(true)

	historyBox := tui.NewVBox(historyScroll)
	historyBox.SetBorder(true)
	historyBox.SetTitle(" Chat ")

	input := tui.NewEntry()
	input.SetFocused(true)
	input.SetSizePolicy(tui.Expanding, tui.Maximum)

	inputBox := tui.NewHBox(input)
	inputBox.SetBorder(true)
	inputBox.SetSizePolicy(tui.Expanding, tui.Maximum)
	inputBox.SetTitle(" Your message ")

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

	if sound {
		box := packr.NewBox("../assets")
		notifRaw, err := box.FindString("notification.mp3")
		if err != nil {
			log.Printf("Error opening notification.mp3: %v", err)
		}

		f, err := os.Create("notification.mp3")
		if err != nil {
			log.Printf("Error while creating notification.mp3: %v", err)
		}
		defer f.Close()

		_, err = f.WriteString(notifRaw)
		if err != nil {
			log.Printf("Error while writing in notification.mp3")
		}
	}

	go func() {

		for {
			r := bufio.NewReader(cnx)
			message, err := r.ReadString('\n')

			if err != nil {
				log.Print(err)
			}
			message = strings.TrimSpace(message)

			if sound {
				// Play sound at each new message
				playSound("notification.mp3")
			}

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

func playSound(track string) {
	f, err := os.Open(track)
	if err != nil {
		log.Fatalf("Error opening %s: %v", track, err)
	}

	streamer, format, err := mp3.Decode(f)
	if err != nil {
		log.Fatal(err)
	}
	defer streamer.Close()

	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
	done := make(chan bool)
	speaker.Play(beep.Seq(streamer, beep.Callback(func() {
		done <- true
	})))

	<-done
}
