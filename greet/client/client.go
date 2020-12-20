package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/devignesh/grpc/greet/greetpb"
)

func main() {
	fmt.Println("gRPC client test")
	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatal("couldn't connect: %v", err)
	}

	defer cc.Close()

	c := greetpb.NewGreetServiceClient(cc)

	doUnary(c, 5*time.Second)
	doUnary(c, 1*time.Second)

}

func doUnary(c greetpb.GreetServiceClient, timeout time.Duration) {
	fmt.Println("Unary api starts...")
	req := &greetpb.GreetRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Vignesh",
			LastName:  "vicky",
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	res, err := c.Greet(ctx, req)
	if err != nil {

		statusErr, ok := status.FromError(err)
		if ok {

			if statusErr.Code() == codes.DeadlineExceeded {
				fmt.Println("Timeout hited, deadline exceeded")
			} else {
				fmt.Printf("unexpected err: %v", statusErr)
			}

		} else {
			log.Fatalf("error while calling dounary: %v", err)
		}

		return
	}
	log.Printf("Response from Greet: %v", res.Result)
}
