// Copyright 2014 The gocui Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"log"

	"github.com/evilgroot/gocui"
)

type manager struct {
	quit bool
}

func (m *manager) Layout(g *gocui.Gui) error {
	if v, err := g.SetView("main", 0, 0, 20, 3); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		if _, err = g.SetCurrentView("main"); err != nil {
			return err
		}
		fmt.Fprint(v, "Press \"q\" to exit")
	}

	if m.quit {
		// show the dialog box
		maxX, maxY := g.Size()
		x0 := (maxX - 20) / 2
		y0 := (maxY - 3) / 2
		if v, err := g.SetView("dialog", x0, y0, x0+20, y0+3); err != nil {
			if err != gocui.ErrUnknownView {
				return err
			}
			if _, err = g.SetCurrentView("dialog"); err != nil {
				return err
			}
			fmt.Fprint(v, "Are you sure? [y/n]")
		}
	}

	return nil
}

func (m *manager) keyBindYes(_ *gocui.Gui, _ *gocui.View) error {
	return gocui.ErrQuit
}

func (m *manager) keyBindNo(g *gocui.Gui, _ *gocui.View) error {
	// remove the dialog box
	m.quit = false
	// set the main view as the current view
	if _, err := g.SetCurrentView("main"); err != nil {
		return err
	}
	return g.DeleteView("dialog")
}

func (m *manager) keyBindQuit(_ *gocui.Gui, _ *gocui.View) error {
	// enable dialog box
	m.quit = true
	return nil
}

func main() {
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Fatal(err)
	}
	defer g.Close()

	man := &manager{}

	g.SetManager(man)

	if err = g.SetKeybinding("main", 'q', gocui.ModNone, man.keyBindQuit); err != nil {
		log.Fatal(err)
	}
	if err = g.SetKeybinding("dialog", 'y', gocui.ModNone, man.keyBindYes); err != nil {
		log.Fatal(err)
	}
	if err = g.SetKeybinding("dialog", 'n', gocui.ModNone, man.keyBindNo); err != nil {
		log.Fatal(err)
	}

	if err = g.MainLoop(); err != gocui.ErrQuit {
		log.Fatal(err)
	}

}
