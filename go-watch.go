// Copyright 2015 Mario Krapp. All rights reserved.

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"time"

	"github.com/jroimartin/gocui"
)

const short_form = "Mon Jan _2 2006 3:04pm"

// struct definitions
type Project struct {
	Time     time.Time
	Duration *time.Duration
	Total    *time.Duration
}

type Projects map[string]Project

// global projects variable
var projects Projects
var projFile = "test.json"

//var projects = map[string]Project{"p1": {time.Now(),0,0}, "p2": {time.Now(),0,0}}
var current_proj Project
var current_name string

// JSON file handling
func read(fnm string) Projects {
	var p Projects
	data, _ := ioutil.ReadFile(fnm)
	json.Unmarshal(data, &p)
	return p
}

func (p Projects) save(fnm string) {
	b, _ := json.Marshal(p)
	ioutil.WriteFile(fnm, b, 0644)
}

// starting/stopping methods
func (p *Project) start() {
	p.Time = time.Now()
	//fmt.Println(p.Time.Format(short_form))
}

func (p *Project) stop() {
	//fmt.Println("Stopping", name)
	//fmt.Println(time.Now().Format(short_form))
	*p.Duration = time.Now().Sub(p.Time)
	*p.Total += *p.Duration
	//p.report()
}

func updateProjInfo(g *gocui.Gui, v *gocui.View) error {
	delView(g, "proj_info")
	maxX, maxY := g.Size()
	if v, err := g.SetView("proj_info", maxX-51, 11, maxX-3, maxY-5); err != nil {
		if err != gocui.ErrorUnkView {
			return err
		}
		fmt.Fprintln(v, current_name)
		fmt.Fprintln(v, "")
		fmt.Fprintln(v, " started on", projects[current_name].Time.Format(short_form))
		fmt.Fprintln(v, " total time spent:", projects[current_name].Total)
		fmt.Fprintln(v, "  last time spent:", projects[current_name].Duration)
	}
	return nil
}

func cursorDown(g *gocui.Gui, v *gocui.View) error {
	var err error
	if v != nil {
		cx, cy := v.Cursor()
		if current_name, err = v.Line(cy + 1); err != nil {
			current_name = ""
		}
		if l, err := v.Line(cy + 2); err != nil && l == "" {
		} else {
			if err := v.SetCursor(cx, cy+1); err != nil {
				ox, oy := v.Origin()
				if err := v.SetOrigin(ox, oy+1); err != nil {
					return err
				}
			}
			updateProjInfo(g, v)
		}
	}
	return nil
}

func cursorUp(g *gocui.Gui, v *gocui.View) error {
	var err error
	if v != nil {
		//_, oy := v.Origin()
		cx, cy := v.Cursor()
		if current_name, err = v.Line(cy - 1); err != nil {
			current_name = ""
		}
		if l, err := v.Line(cy - 1); err != nil && l == "" {
		} else {
			if err := v.SetCursor(cx, cy-1); err != nil {
				ox, oy := v.Origin()
				if err := v.SetOrigin(ox, oy-1); err != nil {
					return err
				}
			}
			updateProjInfo(g, v)
		}
	}
	return nil
}

func getLine(g *gocui.Gui, v *gocui.View) error {
	var l string
	var err error

	g.ShowCursor = false

	_, cy := v.Cursor()
	if l, err = v.Line(cy); err != nil {
		l = ""
	}

	maxX, maxY := g.Size()
	length := 10 + len(current_name)
	if l != "" {
		current_name = l
		if v, err := g.SetView("msg", maxX/2-length/2, maxY/2-3, maxX/2+length/2, maxY/2+3); err != nil {
			v.BgColor = gocui.ColorGreen
			v.FgColor = gocui.ColorBlack
			if err != gocui.ErrorUnkView {
				return err
			}
			current_proj = projects[current_name]
			current_proj.start()
			fmt.Fprintln(v, "")
			fmt.Fprintln(v, "")
			fmt.Fprintln(v, strings.Repeat(" ", (length-15)/2),"Active Project",)
			fmt.Fprintln(v, "")
			fmt.Fprintln(v, strings.Repeat(" ", 5), current_name)
			fmt.Fprintln(v, "")
			setView(g, "msg")
		}
	}
	return nil
}

func delMsg(g *gocui.Gui, v *gocui.View) error {
	current_proj.stop()
	delView(g, "msg")
	setView(g, "list")
	updateProjInfo(g, v)
	return nil
}

func setView(g *gocui.Gui, s string) error {
	if err := g.SetCurrentView(s); err != nil {
		return err
	}
	return nil

}

func delView(g *gocui.Gui, s string) error {
	if err := g.DeleteView(s); err != nil {
		return err
	}
	return nil

}

func saveProj(g *gocui.Gui, v *gocui.View) error {
	if l := strings.TrimSpace(v.Buffer()); l != "" {
		init_t, _ := time.ParseDuration("0s")
		init_d, _ := time.ParseDuration("0s")
		projects[l] = Project{time.Now(), &init_d, &init_t}
	}
	g.ShowCursor = false
	delView(g, "save_proj")
	delView(g, "list")
	g.Flush()
	setView(g, "list")
	updateProjInfo(g, v)
	return nil
}

func delProj(g *gocui.Gui, v *gocui.View) error {
	if l := current_name; l != "" {
		delete(projects, l)
	}
	g.ShowCursor = false
	delView(g, "del_proj")
	delView(g, "list")
	g.Flush()
	setView(g, "list")

	return nil
}

func abortDelProj(g *gocui.Gui, v *gocui.View) error {
	delView(g, "del_proj")
	setView(g, "list")
	return nil
}

func removeProject(g *gocui.Gui, v *gocui.View) error {
	var l string
	var err error

	g.ShowCursor = false

	_, cy := v.Cursor()
	if l, err = v.Line(cy); err != nil {
		l = ""
	}
	current_name = l

	maxX, maxY := g.Size()
	length := 20
	if l != "" {
		if v, err := g.SetView("del_proj", maxX/2-length/2, maxY/2, maxX/2+length/2, maxY/2+2); err != nil {
			v.BgColor = gocui.ColorRed
			if err != gocui.ErrorUnkView {
				return err
			}
			fmt.Fprintln(v, "Press 'd' to delete")
			setView(g, "del_proj")
		}
	}
	return nil
}

func addProject(g *gocui.Gui, v *gocui.View) error {
	g.ShowCursor = true
	maxX, maxY := g.Size()
	if v, err := g.SetView("save_proj", maxX/2-30, maxY/2, maxX/2+30, maxY/2+2); err != nil {
		if err != gocui.ErrorUnkView {
			return err
		}
		setView(g, "save_proj")
		v.Editable = true
	}
	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	projects.save(projFile)
	return gocui.Quit
}

func keybindings(g *gocui.Gui) error {
	if err := g.SetKeybinding("list", 'a', gocui.ModNone, addProject); err != nil {
		return err
	}
	if err := g.SetKeybinding("list", 'd', gocui.ModNone, removeProject); err != nil {
		return err
	}
	if err := g.SetKeybinding("list", gocui.KeyArrowDown, gocui.ModNone, cursorDown); err != nil {
		return err
	}
	if err := g.SetKeybinding("list", gocui.KeyArrowUp, gocui.ModNone, cursorUp); err != nil {
		return err
	}
	if err := g.SetKeybinding("list", gocui.KeyEnter, gocui.ModNone, getLine); err != nil {
		return err
	}
	if err := g.SetKeybinding("msg", gocui.KeyEnter, gocui.ModNone, delMsg); err != nil {
		return err
	}
	if err := g.SetKeybinding("save_proj", gocui.KeyEnter, gocui.ModNone, saveProj); err != nil {
		return err
	}
	if err := g.SetKeybinding("del_proj", 'd', gocui.ModNone, delProj); err != nil {
		return err
	}
	if err := g.SetKeybinding("del_proj", gocui.KeyEnter, gocui.ModNone, abortDelProj); err != nil {
		return err
	}
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		return err
	}

	return nil
}

func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()

	if v, err := g.SetView("info", maxX-53, 9, maxX-1, maxY-4); err != nil {
		if err != gocui.ErrorUnkView {
			return err
		}
		fmt.Fprintln(v, " PROJECT INFORMATION")
	}

	if v, err := g.SetView("list", 0, 1, maxX-54, maxY-1); err != nil {
		if err != gocui.ErrorUnkView {
			return err
		}
		v.Highlight = true
		for k, _ := range projects {
			fmt.Fprintln(v, k)
		}
		setView(g, "list")
		if l := len(projects); l != 0 {
			_, cy := v.Cursor()
			if current_name, err = v.Line(cy); err != nil {
				current_name = ""
			}
			updateProjInfo(g, v)
		}
	}

	if v, err := g.SetView("legend", maxX-53, 1, maxX-1, 8); err != nil {
		if err != gocui.ErrorUnkView {
			return err
		}
		fmt.Fprintln(v, " KEYBINDINGS")
		fmt.Fprintln(v, " A: Add new project")
		fmt.Fprintln(v, " D: Delete a project")
		fmt.Fprintln(v, " Enter: Activate a project")
		fmt.Fprintln(v, " ^C: Exit")
	}

	if v, err := g.SetView("label", maxX-53, maxY-3, maxX-1, maxY-1); err != nil {
		if err != gocui.ErrorUnkView {
			return err
		}
		fmt.Fprintln(v, " This is go-watch -- a time tracker")
	}
	if v, err := g.SetView("listlabel", 0, -1, 20, 1); err != nil {
		if err != gocui.ErrorUnkView {
			return err
		}
		fmt.Fprintln(v, "LIST OF PROJECTS")
		v.Frame = false
	}

	return nil
}

func main() {
	var err error

	projects = read(projFile)

	g := gocui.NewGui()
	if err := g.Init(); err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	g.SetLayout(layout)
	if err := keybindings(g); err != nil {
		log.Panicln(err)
	}
	g.SelBgColor = gocui.ColorGreen
	g.SelFgColor = gocui.ColorBlack
	//g.ShowCursor = true

	err = g.MainLoop()
	if err != nil && err != gocui.Quit {
		log.Panicln(err)
	}
}
