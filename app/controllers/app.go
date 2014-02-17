package controllers

import (
	"chatex/app/models"
	"github.com/robfig/revel"
)

type App struct {
	*revel.Controller
}

func (self App) Index() revel.Result {
	var r = make(map[string][]models.Room)
	r["rooms"] = models.GetRoomList()
	return self.Render(r)
}
