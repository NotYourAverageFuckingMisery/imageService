package service

import (
	"log"

	"github.com/NotYourAverageFuckingMisery/imageService/internal/genErr"
	"github.com/NotYourAverageFuckingMisery/imageService/internal/store"
	v1 "github.com/NotYourAverageFuckingMisery/imageService/proto/v1"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// ImageInfoServer is responsible for getting information about the images
type ImageInfoServer struct {
	v1.UnimplementedImageInfoServiceServer
	*store.DiskImageStore
}

func (s *ImageInfoServer) GetImageList(req *emptypb.Empty, stream v1.ImageInfoService_GetImageListServer) error {
	log.Println("InfoService called")
	info, err := s.GetInfo()
	if err != nil {
		return genErr.NewError(err, genErr.ErrGetingInfo)
	}
	imageList := make([]*v1.ImageInfo, 0, len(info))
	for _, v := range info {
		imageList = append(imageList, &v1.ImageInfo{
			ImageName:    v.Name,
			CreatedAt:    timestamppb.New(v.CreatedAt),
			LastModified: timestamppb.New(v.LastModified),
		})
	}
	resp := &v1.GetImageListResponse{
		ImageList: imageList,
	}

	err = stream.Send(resp)
	if err != nil {
		return genErr.NewError(err, genErr.ErrSendingInfo)
	}

	return nil
}
