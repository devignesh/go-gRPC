package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/devignesh/grpc/bidistream/bidipb"
)

type server struct {
}

func (*server) BidiStream(stream bidipb.GreetService_BidiStreamServer) error {
	fmt.Printf("bi direction streaming api rpc\n")

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			log.Fatalf("Error while reading client stream: %v", err)
			return err
		}

		ctx := context.Background()
		for i := 0; i < 3; i++ {
			if ctx.Err() == context.DeadlineExceeded {

				fmt.Println("The client canceled the request!")
				return status.Error(codes.Canceled, "the client canceled the request")
			}
			time.Sleep(1 * time.Second)
		}

		firstName := req.GetGreeting().GetFirstName()
		lastName := req.GetGreeting().GetLastName()
		result := "Hello " + firstName + " " + lastName + "! "

		sendErr := stream.Send(&bidipb.BidiResponse{
			Result: result,
		})
		if sendErr != nil {
			log.Fatalf("Error while sending data to client: %v", sendErr)
			return sendErr
		}
	}

}

func main() {

	fmt.Println("gRPC server tests")

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("failed to listen %v", err)
	}

	s := grpc.NewServer()
	bidipb.RegisterGreetServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}
