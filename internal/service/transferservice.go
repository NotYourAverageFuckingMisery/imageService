package service

import (
	"bytes"
	"log"
	"os"
	"path/filepath"

	"github.com/NotYourAverageFuckingMisery/imageService/internal/genErr"
	"github.com/NotYourAverageFuckingMisery/imageService/internal/store"
	v1 "github.com/NotYourAverageFuckingMisery/imageService/proto/v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

type TransferImageServer struct {
	v1.UnimplementedTransferImageServiceServer
	*store.DiskImageStore
}

func (s *TransferImageServer) Upload(stream v1.TransferImageService_UploadServer) error {
	req, err := stream.Recv()
	if err != nil {
		log.Println(err)
		return err
	}
	imageName := req.GetImageInfo().ImageName
	imageData := bytes.Buffer{}

	log.Println("Geting image data...")

	req, err = stream.Recv()
	if err != nil {
		return genErr.NewError(err, genErr.ErrRecievingData)
	}
	bytes := req.GetImage()
	_, err = imageData.Write(bytes)
	if err != nil {
		return genErr.NewError(err, genErr.ErrWritingData)
	}

	err = s.Save(imageName, imageData)
	if err != nil {
		return genErr.NewError(err, genErr.ErrFailedToSave)
	}

	res := &emptypb.Empty{}

	err = stream.SendAndClose(res)
	if err != nil {
		return genErr.NewError(err, genErr.ErrClosingStream)
	}

	log.Println("saved image to disk")

	return nil
}

func (s *TransferImageServer) Download(req *v1.DownloadRequest, stream v1.TransferImageService_DownloadServer) error {
	file, err := os.Open(s.ImageFolder + "/" + req.Filename)
	if err != nil {
		return genErr.NewError(err, genErr.ErrOpeningFile)
	}
	defer file.Close()

	res := &v1.DownloadResponse{
		Data: &v1.DownloadResponse_Info{
			Info: &v1.Info{
				ImageType: filepath.Ext(s.ImageFolder + "/" + req.Filename),
				ImageName: req.Filename,
			},
		},
	}
	err = stream.Send(res)
	if err != nil {
		return genErr.NewError(err, genErr.ErrSendingImage)
	}

	bytes, err := os.ReadFile(s.ImageFolder + "/" + req.Filename)
	if err != nil {
		return genErr.NewError(err, genErr.ErrReadingFile)
	}

	res = &v1.DownloadResponse{
		Data: &v1.DownloadResponse_Image{
			Image: bytes,
		},
	}
	err = stream.Send(res)
	if err != nil {
		return genErr.NewError(err, genErr.ErrSendingData)
	}

	log.Println("Download completed")

	return nil
}
