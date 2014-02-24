package controllers

import (
	"github.com/robfig/revel"
	"io"
	"os"
	"path/filepath"
)

type File struct {
	*revel.Controller
}

func (self *File) Upload(image io.ReadSeeker) revel.Result {

	imageMetaAll := self.Params.Files["image"]
	if len(imageMetaAll) == 0 {
		revel.WARN.Println("File upload without multipart header, possible canceled by user.")
		self.Response.Status = 400
		return self.RenderJson(map[string]string{"error": "Bad upload strategy"})
	}

	imageMeta := imageMetaAll[0]

	if GetUploadFileSize(image) > MAX_FILE_SIZE {
		revel.INFO.Println("File is too large. Please check app.max_file_size option")
		self.Response.Status = 400
		return self.RenderJson(map[string]string{"error": "File too large"})
	}

	name := GetSha1Hex(image) + filepath.Ext(imageMeta.Filename)
	path := filepath.Join(MEDIA_PATH, name)
	url := MEDIA_URL + name

	revel.INFO.Printf("Saving file to %s. URL: %s", path, url)

	dest, err := os.Create(path)

	if err != nil {
		panic(err)
	}
	io.Copy(dest, image)

	err = dest.Close()
	if err != nil {
		panic(err)
	}

	return self.RenderJson(map[string]string{"url": url})
}
