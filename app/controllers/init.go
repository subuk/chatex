package controllers

import (
	"github.com/robfig/revel"
	"github.com/robfig/revel/modules/db/app"
)

var MEDIA_URL string
var MEDIA_PATH string
var MAX_FILE_SIZE int64
var UPLOAD_FILE_PERMISSION int

func init() {
	revel.OnAppStart(db.Init)
	revel.OnAppStart(getConfig)
}

func getConfig() {
	MEDIA_URL = revel.Config.StringDefault("app.media_url", "/media/")
	MEDIA_PATH = revel.Config.StringDefault("app.media_path", "/var/tmp")
	MAX_FILE_SIZE = int64(revel.Config.IntDefault("app.max_file_size", 10*1024*1024))
	UPLOAD_FILE_PERMISSION = revel.Config.IntDefault("app.upload_file_permission", 0644)
	revel.INFO.Printf("MEDIA_URL (%s), MEDIA_PATH(%s), MAX_FILE_SIZE(%d) fetched from config", MEDIA_URL, MEDIA_PATH, MAX_FILE_SIZE)
}
