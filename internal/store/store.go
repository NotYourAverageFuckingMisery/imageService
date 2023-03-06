package store

import (
	"bytes"
	"log"
	"os"
	"syscall"
	"time"

	"github.com/NotYourAverageFuckingMisery/imageService/internal/models"
)

// This type is responsible for storing and getting images from the disk.
type DiskImageStore struct {
	ImageFolder string
}

func NewImageStore(path string) *DiskImageStore {
	return &DiskImageStore{
		ImageFolder: path,
	}
}

func (d *DiskImageStore) Save(imageName string, imageData bytes.Buffer) error {
	file, err := os.Create(d.ImageFolder + "/" + imageName)
	if err != nil {
		log.Println(err)
		return err
	}
	_, err = imageData.WriteTo(file)
	if err != nil {
		os.Remove(d.ImageFolder + file.Name())
		log.Println(err)
		return err
	}

	return nil
}

func (d *DiskImageStore) GetInfo() ([]models.FileInfo, error) {
	files, err := os.ReadDir(d.ImageFolder)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	infoList := make([]models.FileInfo, 0, len(files))
	for _, file := range files {

		fi, err := os.Stat(d.ImageFolder + "/" + file.Name())
		if err != nil {
			return nil, err
		}
		stat := fi.Sys().(*syscall.Stat_t)
		ctime := stat.Birthtimespec.Sec

		info, err := file.Info()
		if err != nil {
			return nil, err
		}
		infoList = append(infoList, models.FileInfo{
			Name:         info.Name(),
			CreatedAt:    time.Unix(ctime, 0),
			LastModified: info.ModTime(),
		})
	}

	return infoList, nil
}
