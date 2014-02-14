package controllers

import (
	"github.com/robfig/revel"
)

type App struct {
	*revel.Controller
}

func (self App) Index() revel.Result {
	return self.Render()
}
