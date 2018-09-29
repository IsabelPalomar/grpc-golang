package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"

	"github.com/grpc-golang/calculator/calculatorpb"
	"google.golang.org/grpc"
)

type server struct{}

func (*server) Calculator(ctx context.Context, req *calculatorpb.CalculatorRequest) (*calculatorpb.CalculatorResponse, error) {
	fmt.Printf("Calculator function was invoked with %v \n", req)
	firstValue := req.GetCalculator().GetFirstValue()
	secondValue := req.GetCalculator().GetSecondValue()
	result := firstValue + secondValue
	res := &calculatorpb.CalculatorResponse{
		Result: result,
	}
	return res, nil
}

func (*server) PrimeNumberDecomposition(req *calculatorpb.PrimeNumberDecompositionRequest, stream calculatorpb.CalculatorService_PrimeNumberDecompositionServer) error {
	fmt.Printf("Prime number Decomposition was invoked with %v \n", req)

	number := req.GetNumber()
	divisor := int64(2)

	for number > 1 {
		if number%divisor == 0 {
			stream.Send(&calculatorpb.PrimeNumberDecompositionResponse{
				PrimeFactor: divisor,
			})
			number = number / divisor
		} else {
			divisor++
			fmt.Printf("Divisor has increased to %v \n", divisor)
		}
	}
	return nil
}

func (*server) ComputeAverage(stream calculatorpb.CalculatorService_ComputeAverageServer) error {
	fmt.Println("Compute Average function was invoked with an streaming request")

	count := 0
	sum := int32(0)

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			// we have finished reading the client stream
			average := float64(sum) / float64(count)
			return stream.SendAndClose(&calculatorpb.ComputeAverageResponse{
				Average: average,
			})
		}
		if err != nil {
			log.Fatalf("Error while reading client stream: %v", err)
		}

		sum += int32(req.GetNumber())
		count++
	}

}

func main() {
	fmt.Println("Hello Server")

	lis, err := net.Listen("tcp", "0.0.0.0:50052")

	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	calculatorpb.RegisterCalculatorServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
