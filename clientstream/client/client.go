package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/devignesh/grpc/clientstream/clientpb"
)

func main() {
	fmt.Println("gRPC client test")
	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatal("couldn't connect: %v", err)
	}

	defer cc.Close()

	c := clientpb.NewGreetServiceClient(cc)

	clientstreaming(c, 10*time.Second)
	clientstreaming(c, 1*time.Second)

}

func clientstreaming(c clientpb.GreetServiceClient, timeout time.Duration) {
	fmt.Println("Client streaming gRPC...")

	requests := []*clientpb.ClientManytimeRequest{

		&clientpb.ClientManytimeRequest{
			Greeting: &clientpb.Greeting{
				FirstName: "VIgnesh",
				LastName:  "vicky",
			},
		},

		&clientpb.ClientManytimeRequest{
			Greeting: &clientpb.Greeting{
				FirstName: "Mani",
				LastName:  "d",
			},
		},

		&clientpb.ClientManytimeRequest{
			Greeting: &clientpb.Greeting{
				FirstName: "Suji",
				LastName:  "Paithiyma",
			},
		},

		&clientpb.ClientManytimeRequest{
			Greeting: &clientpb.Greeting{
				FirstName: "subash",
				LastName:  "aunty",
			},
		},

		&clientpb.ClientManytimeRequest{
			Greeting: &clientpb.Greeting{
				FirstName: "dogy",
				LastName:  "vicky",
			},
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	stream, err := c.ClientManytime(ctx)

	if err != nil {
		log.Fatalf("error calling clientstreaming: %v", err)
	}

	for _, req := range requests {
		fmt.Printf("Sending req: %v\n", req)
		stream.Send(req)
		time.Sleep(1000 * time.Millisecond)
	}
	res, err := stream.CloseAndRecv()

	if err != nil {
		statusErr, ok := status.FromError(err)
		if ok {

			if statusErr.Code() == codes.DeadlineExceeded {
				fmt.Println("Timeout hited, deadline exceeded")
			} else {
				fmt.Printf("unexpected err: %v", statusErr)
			}

		} else {
			log.Fatalf("error while calling clientstream: %v", err)
		}

		return
	}
	fmt.Printf("Client stream Response: %v\n", res)

}
