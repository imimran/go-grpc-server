package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"

	pb "github.com/imimran/go-grpc-server/fileservice"
	"google.golang.org/grpc"
)

// server implements the gRPC ImageServiceServer interface
type server struct {
	pb.UnimplementedImageServiceServer
}

// UploadImage handles image upload from client
func (s *server) UploadImage(ctx context.Context, req *pb.ImageRequest) (*pb.ImageResponse, error) {
	// Create uploads folder if it does not exist
	if _, err := os.Stat("uploads"); os.IsNotExist(err) {
		os.Mkdir("uploads", 0755)
	}

	// Save file
	filename := fmt.Sprintf("uploads/%s", req.Filename)
	err := ioutil.WriteFile(filename, req.Data, 0644)
	if err != nil {
		return &pb.ImageResponse{
			Success: false,
			Message: fmt.Sprintf("Failed to save file: %v", err),
		}, nil
	}

	log.Printf("✅ Received file: %s from user: %s meta: %s", req.Filename, req.UserId, req.Meta)

	return &pb.ImageResponse{
		Success: true,
		Message: "Image uploaded successfully",
	}, nil
}

func main() {
	// Listen on TCP port 50051
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("❌ failed to listen: %v", err)
	}

	// Create gRPC server
	s := grpc.NewServer()
	pb.RegisterImageServiceServer(s, &server{})

	log.Println("🚀 Go gRPC server running on port 50051")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("❌ failed to serve: %v", err)
	}
}
