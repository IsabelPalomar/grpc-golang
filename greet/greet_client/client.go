package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"
	"github.com/grpc-golang/greet/greetpb"

	"google.golang.org/grpc"
)

type server struct{}

func main() {
	fmt.Println("Hello, I'm a client")
	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	defer cc.Close()

	c := greetpb.NewGreetServiceClient(cc)
	
	// doUnary(c)
	//doServerStreaming(c)
	doClientStreaming(c)

}

func doUnary(c greetpb.GreetServiceClient) {
	fmt.Println("Starting to do a Unary RPC...")
	req := &greetpb.GreetRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Isabel",
			LastName:  "Palomar",
		},
	}
	res, err := c.Greet(context.Background(), req)

	if err != nil {
		log.Fatalf("Error while calling Greet RPC: %v", err)
	}

	log.Printf("Response from Greet: %v", res.Result)
}

func doServerStreaming(c greetpb.GreetServiceClient) {
	fmt.Println("Starting to do a Serving Streaming RPC...")

	req := &greetpb.GreetManyTimesRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Isabel",
			LastName:  "Palomar",
		},
	}
	resStream, err := c.GreetManyTimes(context.Background(), req)

	if err != nil {
		log.Fatalf("Error while calling GreetManyTimes RPC: %v", err)
	}

	for {
		msg, err := resStream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Error while reading stream: %v", err)
		}
		log.Printf("Response from GreetManyTimes: %v", msg.GetResult())
	}

}

func doClientStreaming(c greetpb.GreetServiceClient){
	fmt.Println("Starting to do a Client Streaming RPC...")

	request := []*greetpb.LongGreetRequest{
		&greetpb.LongGreetRequest {
			Greeting: &greetpb.Greeting{
				FirstName: "Isabel",
			},
		},
		&greetpb.LongGreetRequest {
			Greeting: &greetpb.Greeting{
				FirstName: "Paty",
			},
		},
		&greetpb.LongGreetRequest {
			Greeting: &greetpb.Greeting{
				FirstName: "Jorge",
			},
		},
		&greetpb.LongGreetRequest {
			Greeting: &greetpb.Greeting{
				FirstName: "Bau",
			},
		},
		&greetpb.LongGreetRequest {
			Greeting: &greetpb.Greeting{
				FirstName: "Jorge PG",
			},
		},
	}

	stream, err := c.LongGreet(context.Background())
	if err != nil {
		log.Fatalf("Error while calling LongGreet: %v", err)
	}

	// we iterate over our slice and send each message individually
	for _, req := range request {
		fmt.Printf("Sending req: %v \n", req)
		stream.Send(req)
		time.Sleep(1000 * time.Millisecond)
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("error whiele receiving response from LongGreet: %v", err)
	}

	fmt.Printf("Long greet response %v \n", res)

}
