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

	"github.com/devignesh/grpc/clientstream/clientpb"
)

type server struct{}

var ctx = context.Background()

func (*server) ClientManytime(stream clientpb.GreetService_ClientManytimeServer) error {
	fmt.Printf("Client stream rpc service\n")

	for i := 0; i < 3; i++ {
		if ctx.Err() == context.DeadlineExceeded {

			fmt.Println("The client canceled the request!")
			return status.Error(codes.Canceled, "the client canceled the request")
		}
		time.Sleep(1 * time.Second)
	}
	result := ""
	for {
		req, err := stream.Recv()
		if err == io.EOF {

			return stream.SendAndClose(&clientpb.ClientManytimeResponse{
				Result: result,
			})
		}
		if err != nil {
			log.Fatalf("Error while reading client stream: %v", err)
		}

		firstName := req.GetGreeting().GetFirstName()
		lastName := req.GetGreeting().GetLastName()
		result += "Hello" + firstName + " " + lastName + "! "
	}
}

func main() {

	fmt.Println("gRPC server tests")

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("failed to listen %v", err)
	}

	s := grpc.NewServer()
	clientpb.RegisterGreetServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}
