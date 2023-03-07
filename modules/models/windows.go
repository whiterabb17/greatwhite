package models

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/asticode/go-astikit"
	"github.com/asticode/go-astilectron"
	bootstrap "github.com/asticode/go-astilectron-bootstrap"
)

func AbsoluteResourcesPath(a *astilectron.Astilectron, relativeResourcesPath string) string {
	return filepath.Join(a.Paths().DataDirectory(), relativeResourcesPath)
}

func handleMessages(w *astilectron.Window, messageHandler bootstrap.MessageHandler, l astikit.SeverityLogger) astilectron.ListenerMessage {
	return func(m *astilectron.EventMessage) (v interface{}) {
		// Unmarshal message
		var i bootstrap.MessageIn
		var err error
		if err = m.Unmarshal(&i); err != nil {
			l.Error(fmt.Errorf("unmarshaling message %+v failed: %w", *m, err))
			return
		}

		// Handle message
		var p interface{}
		if p, err = messageHandler(w, i); err != nil {
			l.Error(fmt.Errorf("handling message %+v failed: %w", i, err))
		}

		// Return message
		if p != nil {
			o := &bootstrap.MessageOut{Name: i.Name + ".callback", Payload: p}
			if err != nil {
				o.Name = "error"
			}
			v = o
		}
		return
	}
}

func CreateWindow(window []*bootstrap.Window) (win []*astilectron.Window, err error) {
	func() {
		l := astikit.AdaptStdLogger(Logger)
		var w = make([]*astilectron.Window, len(WindowOpts))
		for i, wo := range WindowOpts {
			var url = wo.Homepage
			if !strings.Contains(url, "://") && !strings.HasPrefix(url, string(filepath.Separator)) {
				url = filepath.Join(AbsoluteResourcesPath(AppPtr, "resources"), "app", url)
			}
			if w[i], err = App.NewWindow(url, wo.Options); err != nil {
				err = fmt.Errorf("new window failed: %w", err)
			}

			// Handle messages
			if wo.MessageHandler != nil {
				w[i].OnMessage(handleMessages(w[i], wo.MessageHandler, l))
			}

			// Adapt window
			if wo.Adapter != nil {
				wo.Adapter(w[i])
			}

			// Create window
			if err = w[i].Create(); err != nil {
				err = fmt.Errorf("creating window failed: %w", err)
			}
			win = append(win, w[i])
		}
		// Create menu options
		mo := AppOpts.MenuOptions
		if AppOpts.MenuOptionsFunc != nil {
			mo = AppOpts.MenuOptionsFunc(AppPtr)
		}

		// Debug

		if AppOpts.Debug {
			// Create menu item
			var debug bool
			mi := &astilectron.MenuItemOptions{
				Accelerator: astilectron.NewAccelerator("Control", "d"),
				Label:       astikit.StrPtr("Debug"),
				OnClick: func(e astilectron.Event) (deleteListener bool) {
					for i, window := range w {
						width := *AppOpts.Windows[i].Options.Width
						if debug {
							if err := window.CloseDevTools(); err != nil {
								l.Error(fmt.Errorf("closing dev tools failed: %w", err))
							}
							if err := window.Resize(width, *AppOpts.Windows[i].Options.Height); err != nil {
								l.Error(fmt.Errorf("resizing window failed: %w", err))
							}
						} else {
							if err := window.OpenDevTools(); err != nil {
								l.Error(fmt.Errorf("opening dev tools failed: %w", err))
							}
							if err := window.Resize(width+700, *AppOpts.Windows[i].Options.Height); err != nil {
								l.Error(fmt.Errorf("resizing window failed: %w", err))
							}
						}
					}

					debug = !debug
					return
				},
				Type: astilectron.MenuItemTypeCheckbox,
			}

			// Add menu item
			if len(mo) == 0 {
				mo = []*astilectron.MenuItemOptions{{SubMenu: []*astilectron.MenuItemOptions{mi}}}
			} else {
				if len(mo[0].SubMenu) > 0 {
					mo[0].SubMenu = append(mo[0].SubMenu, &astilectron.MenuItemOptions{Type: astilectron.MenuItemTypeSeparator})
				}
				mo[0].SubMenu = append(mo[0].SubMenu, mi)
			}
		}

		// Menu
		var m *astilectron.Menu
		if len(mo) > 0 {
			// Init menu
			m = App.NewMenu(mo)

			// Create menu
			if err = m.Create(); err != nil {
				err = fmt.Errorf("creating menu failed: %w", err)
			}
		}

	}()
	return
}
