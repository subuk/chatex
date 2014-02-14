package controllers

import (
	"github.com/robfig/revel"
	"github.com/robfig/revel/modules/db/app"
)

func init() {
	revel.OnAppStart(db.Init)
}
