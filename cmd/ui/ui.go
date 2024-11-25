package main

import (
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"go_stream_events/pkg"
	"image/color"
	"log"
	"os"

	"gioui.org/app"
	"gioui.org/op"
	"gioui.org/text"
	"gioui.org/widget/material"
)

func main() {
	go func() {
		window := new(app.Window)
		window.Option(app.Title("Event UI"))
		window.Option(app.Size(unit.Dp(400), unit.Dp(600)))
		err := run(window)
		if err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()
	app.Main()
}

func run(window *app.Window) error {
	theme := material.NewTheme()
	var ops op.Ops
	var ed = widget.Editor{
		Alignment:  text.Middle,
		SingleLine: false,
		ReadOnly:   false,
		Submit:     true,
	}
	for {
		switch e := window.Event().(type) {
		case app.DestroyEvent:
			return e.Err
		case app.FrameEvent:
			// This graphics context is used for managing the rendering state.
			gtx := app.NewContext(&ops, e)

			// Define an large label with an appropriate text:
			title := material.H2(theme, "Send Test Message")

			// Change the color of the label.
			maroon := color.NRGBA{R: 127, G: 0, B: 0, A: 255}
			title.Color = maroon

			// Change the position of the label.
			title.Alignment = text.Middle

			layout.Flex{
				// Vertical alignment, from top to bottom
				Axis: layout.Vertical,
				// Empty space is left at the start, i.e. at the top
				Spacing: layout.SpaceStart,
			}.Layout(
				gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					margins := layout.Inset{
						Top:    unit.Dp(0),
						Right:  unit.Dp(10),
						Bottom: unit.Dp(20),
						Left:   unit.Dp(10),
					}

					// ... and borders ...
					border := widget.Border{
						Color:        color.NRGBA{R: 204, G: 204, B: 204, A: 255},
						CornerRadius: unit.Dp(3),
						Width:        unit.Dp(2),
					}
					editor := material.Editor(theme, &ed, "Add message")

					e, _ := ed.Update(gtx)
					if e, ok := e.(widget.SubmitEvent); ok {
						consumer := pkg.NewConsumerClient(e.Text)
						consumer.Connect()
						poll, _ := consumer.Poll()

						for i := uint64(0); i < uint64(len(poll)); i++ {
							log.Println(poll[i])
						}
					}

					return margins.Layout(gtx,
						func(gtx layout.Context) layout.Dimensions {
							return border.Layout(gtx, editor.Layout)
						},
					)
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return title.Layout(gtx)
				}),
			)
			e.Frame(gtx.Ops)
		}
	}
}
