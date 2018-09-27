package main

import (
	"context"
	"fmt"
	"io"
	"log"

	"github.com/grpc-golang/calculator/calculatorpb"
	"google.golang.org/grpc"
)

type server struct{}

func main() {
	fmt.Println("Hello, I'm the calcualtor client")
	cc, err := grpc.Dial("localhost:50052", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	defer cc.Close()

	c := calculatorpb.NewCalculatorServiceClient(cc)
	//doUnary(c)

	doServerStreaming(c)

}

func doUnary(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("Starting to do a Unary RPC ...")
	req := &calculatorpb.CalculatorRequest{
		Calculator: &calculatorpb.Calculator{
			FirstValue:  300,
			SecondValue: 100,
		},
	}

	res, err := c.Calculator(context.Background(), req)

	if err != nil {
		log.Fatalf("Error while calling Calculator RPC: %v", err)
	}

	log.Printf("Response from Calculator: %v", res)
}

func doServerStreaming(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("Starting to do a PrimeDecomposition Server Streaming RPC ...")
	req := &calculatorpb.PrimeNumberDecompositionRequest{
		Number: 1988,
	}

	stream, err := c.PrimeNumberDecomposition(context.Background(), req)

	if err != nil {
		log.Fatalf("Error while calling PrimeDecomposition RPC: %v", err)
	}

	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Something happened: %v", err)
		}
		fmt.Println(res.GetPrimeFactor())
	}
}
