package serving

import (
	"fmt"
	// "log"
	"net/http"
	"time"
	
	"github.com/nsf/termbox-go"
)

const defaultRatio float64 = 7.0 / 3.0 // The terminal's default cursor width/height ratio

var (
	fg termbox.Attribute = termbox.ColorRed
	bg termbox.Attribute = termbox.ColorDefault
)

var box [][]rune

// var X = 0

// var Y int = 0

func tbprint(x, y int, fg, bg termbox.Attribute, msg string) {
	box = append(box, []rune{})
	for _, c := range msg {
		box[len(box)-1] = append(box[len(box)-1], c)
	}
}

func print(msg string) {
	if msg != "" {
		tbprint(0, len(box), fg, bg, msg)
	} else {
		box = append(box, []rune{})
	}
}

func Serve(path string) error {
	// fmt.Println("Welcome to serve")

	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()

	termbox.SetOutputMode(termbox.Output256)

	termbox.SetOutputMode(termbox.Output256)

	// print("Hey there!")
	fmt.Printf("Serving at http://localhost:8080/\n")
	print(fmt.Sprintf("Serving at http://localhost:8080/\n"))

	go func() {
		http.Handle("/", http.FileServer(http.Dir("./dist")))

		// http.Handle("/view", http.FileServer(http.Dir("./static")))
		// print(fmt.Sprintf("Serving at http://localhost:8080/%s\n", path))

		if err := http.ListenAndServe(":8080", nil); err != nil {
			print(err.Error())
			time.Sleep(4 * time.Second)
			return
		}
	}()

	draw()
	tbprint(12, 0, termbox.ColorWhite, termbox.ColorDefault, "helpMsg")
	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			if ev.Key == termbox.KeyEsc || ev.Ch == 'q' || ev.Key == termbox.KeyCtrlC {
				return nil
			}
			if ev.Key == termbox.KeyEnter {
				return nil
			}
		case termbox.EventResize:
			draw()
		case termbox.EventMouse:

		default:
			time.Sleep(10 * time.Millisecond)
		}
	}
}

func draw() {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	defer termbox.Flush()

	for y, v := range box {
		for x, m := range v {
			termbox.SetCell(x, y, m, fg, bg)
		}
	}
}