package main

import (
  _ "embed"
  "github.com/wailsapp/wails"
)

func basic() string {
  return "Hello World!"
}

//go:embed frontend/dist/app.js
var js string

//go:embed frontend/dist/app.css
var css string

func main() {

  app := wails.CreateApp(&wails.AppConfig{
    Width:  1280,
    Height: 800,
    Title:  "git-inside",
    JS:     js,
    CSS:    css,
    Colour: "#ffffff",
  })
  app.Bind(basic)

  git, _ := NewGit()
  app.Bind(git)

  app.Run()
}
