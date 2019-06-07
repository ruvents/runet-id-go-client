package models

import (
	"fmt"
	"os"
)

const cfgWebRoot = "/backend/www/files/"
const cfgWebHost = "https://static.runet-id.com/"
const cfgEmptyPhoto = "files/photo/nophoto_200.png"

type UserPhoto struct {
	Original   string
	ModifiedAt int64
}

func (photo *UserPhoto) SetUserID(id uint32) {
	relativePhotoPath := fmt.Sprintf("photo/%d/%d.jpg", (id - id % 10000) / 10000, id)
	if file, err := os.Stat(cfgWebRoot + relativePhotoPath); err == nil {
		photo.Original = cfgWebHost + relativePhotoPath
		photo.ModifiedAt = file.ModTime().Unix()
	} else {
		photo.Original = cfgWebHost + cfgEmptyPhoto
	}
}
