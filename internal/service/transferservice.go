package service

import (
	"bufio"
	"bytes"
	"io"
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
	for {

		req, err := stream.Recv()
		if err == io.EOF {
			log.Print("no more data")
			break
		}
		if err != nil {
			return genErr.NewError(err, genErr.ErrRecievingChunkData)
		}
		chunk := req.GetImage()
		_, err = imageData.Write(chunk)
		if err != nil {
			return genErr.NewError(err, genErr.ErrWritingChunkData)
		}
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

	reader := bufio.NewReader(file)
	buffer := make([]byte, 1024)

	for {
		n, err := reader.Read(buffer)
		if err == io.EOF {
			break
		}
		if err != nil {
			return genErr.NewError(err, genErr.ErrReadingChunk)
		}

		res := &v1.DownloadResponse{
			Data: &v1.DownloadResponse_Image{
				Image: buffer[:n],
			},
		}
		err = stream.Send(res)
		if err != nil {
			return genErr.NewError(err, genErr.ErrSendingChunkData)
		}
	}

	log.Println("Download completed")

	return nil
}
