package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/devignesh/grpc/bidistream/bidipb"
)

func main() {
	fmt.Println("gRPC client test")
	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatal("couldn't connect: %v", err)
	}

	defer cc.Close()

	c := bidipb.NewGreetServiceClient(cc)

	bidistreaming(c, 50*time.Second)
	bidistreaming(c, 1*time.Second)

}

func bidistreaming(c bidipb.GreetServiceClient, timeout time.Duration) {
	fmt.Println("Starting to do a BiDi Streaming RPC...")

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	stream, err := c.BidiStream(ctx)
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

		// return
		log.Fatalf("Error while creating stream: %v", err)
		return
	}

	requests := []*bidipb.BidiRequest{
		&bidipb.BidiRequest{
			Greeting: &bidipb.Greeting{
				FirstName: "Vignesh",
				LastName:  "vicky",
			},
		},
		&bidipb.BidiRequest{
			Greeting: &bidipb.Greeting{
				FirstName: "Mani",
				LastName:  "ma",
			},
		},
		&bidipb.BidiRequest{
			Greeting: &bidipb.Greeting{
				FirstName: "subash",
				LastName:  "mental",
			},
		},
		&bidipb.BidiRequest{
			Greeting: &bidipb.Greeting{
				FirstName: "shofia",
				LastName:  "paithiyam",
			},
		},
		&bidipb.BidiRequest{
			Greeting: &bidipb.Greeting{
				FirstName: "sangee",
				LastName:  "loosu",
			},
		},
	}

	waitc := make(chan struct{})

	go func() {

		for _, req := range requests {
			fmt.Printf("Sending message: %v\n", req)
			stream.Send(req)
			time.Sleep(1000 * time.Millisecond)
		}
		stream.CloseSend()
	}()
	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {

				log.Fatalf("Error while receiving: %v", err)
				break
			}
			fmt.Printf("Received: %v\n", res.GetResult())
		}
		close(waitc)
	}()

	<-waitc
}
