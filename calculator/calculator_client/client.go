package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/grpc-golang/calculator/calculatorpb"
	"google.golang.org/grpc"
)

type server struct{}

func main() {
	fmt.Println("Hello, I'm the calculator client")
	cc, err := grpc.Dial("localhost:50052", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	defer cc.Close()

	c := calculatorpb.NewCalculatorServiceClient(cc)
	//doUnary(c)
	//doServerStreaming(c)
	doClientStreaming(c)

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

func doClientStreaming(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("Starting to do a Client Streaming RPC...")

	request := []*calculatorpb.ComputeAverageRequest{
		&calculatorpb.ComputeAverageRequest{
			Number: 1,
		},
		&calculatorpb.ComputeAverageRequest{
			Number: 2,
		},
		&calculatorpb.ComputeAverageRequest{
			Number: 3,
		},
		&calculatorpb.ComputeAverageRequest{
			Number: 4,
		},
	}

	stream, err := c.ComputeAverage(context.Background())
	if err != nil {
		log.Fatalf("Error while calling Compute Average: %v", err)
	}

	// we iterate over our slice and send each number individually
	for _, req := range request {
		fmt.Printf("Sending req: %v \n", req)
		stream.Send(req)
		time.Sleep(1000 * time.Millisecond)
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("error while receiving response from ComputeAverage: %v", err)
	}

	fmt.Printf("Compute Average response %v \n", res)

}
