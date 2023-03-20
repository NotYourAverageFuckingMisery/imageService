package main

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	v1 "github.com/NotYourAverageFuckingMisery/imageService/proto/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"
)

func getInfo(client v1.ImageInfoServiceClient) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	stream, err := client.GetImageList(ctx, &emptypb.Empty{})
	if err != nil {
		log.Fatal("cannot get image list from a server: ", err)
	}
	data, err := stream.Recv()
	if err != nil {
		log.Fatal("cannot get image list from a server: ", err)
	}

	type info struct {
		name    string
		changed time.Time
		created time.Time
	}
	infoList := make([]info, 0)
	for _, i := range data.ImageList {
		infoList = append(infoList, info{
			name:    i.ImageName,
			changed: i.LastModified.AsTime(),
			created: i.CreatedAt.AsTime(),
		})
	}

	err = stream.CloseSend()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(infoList)

}
func upload(client v1.TransferImageServiceClient, imagePath string) {
	file, err := os.Open(imagePath)
	if err != nil {
		log.Fatal("Error opening image: ", err)
	}
	defer file.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	stream, err := client.Upload(ctx)
	if err != nil {
		log.Fatal("cannot send image info to server: ", err, stream.RecvMsg(nil))
	}

	name := strings.TrimPrefix(file.Name(), "client/img/")

	req := &v1.UploadRequest{
		Data: &v1.UploadRequest_ImageInfo{
			ImageInfo: &v1.Info{
				ImageType: filepath.Ext(imagePath),
				ImageName: name,
			},
		},
	}
	err = stream.Send(req)
	if err != nil {
		log.Fatal("cannot send image info to server: ", err, stream.RecvMsg(nil))
	}

	bytes, err := os.ReadFile(imagePath)
	if err != nil {
		log.Fatal("can not read file")
	}
	req = &v1.UploadRequest{
		Data: &v1.UploadRequest_Image{
			Image: bytes,
		},
	}
	err = stream.Send(req)
	if err != nil {
		log.Fatal("cannot send chunk to server: ", err)
	}

	_, err = stream.CloseAndRecv()
	if err != nil {
		log.Fatal("cannot close stream: ", err)
	}

	log.Println("upload completed")
}

func download(client v1.TransferImageServiceClient, imagePath string) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	stream, err := client.Download(ctx, &v1.DownloadRequest{Filename: "jojo.jpeg"})
	if err != nil {
		log.Fatal("cannot get image from a server: ", err)
	}

	req, err := stream.Recv()
	if err != nil {
		log.Println(err)
	}

	imageName := req.GetInfo().ImageName
	imageData := bytes.Buffer{}

	log.Println("Geting image data...")

	req, err = stream.Recv()
	if err != nil {
		log.Fatal("Can not recive data")
	}
	chunk := req.GetImage()
	_, err = imageData.Write(chunk)
	if err != nil {
		log.Fatal("Can not write data")
	}

	err = stream.CloseSend()
	if err != nil {
		log.Fatal(err)
	}

	file, err := os.Create(imagePath + imageName)
	if err != nil {
		log.Fatal(err)
	}
	_, err = imageData.WriteTo(file)
	if err != nil {
		log.Fatal(err)
	}

}

func main() {
	conn, err := grpc.Dial("0.0.0.0:5051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	c := v1.NewTransferImageServiceClient(conn)
	conn2, err := grpc.Dial("0.0.0.0:5069", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	c2 := v1.NewImageInfoServiceClient(conn2)
	wg := sync.WaitGroup{}
	for {
		time.Sleep(10 * time.Millisecond)
		wg.Add(3)
		//t := time.Now()
		go func() {
			upload(c, "client/img/20-facts-might-know-bad-santa.jpg")
			wg.Done()
		}()
		go func() {
			download(c, "client/img/")
			wg.Done()
		}()
		go func() {
			getInfo(c2)
			wg.Done()
		}()
	}
	//wg.Wait()
	//fmt.Println(time.Since(t))
}
