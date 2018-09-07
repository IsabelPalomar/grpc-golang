package main

import (
	"context"
	"fmt"
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
