package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/go-grpc-udemy-course/go-grpc-course/calculator/calculatorpb"

	"google.golang.org/grpc"
)

type server struct{}

// Attached method to the struct
func (*server) Sum(ctx context.Context, req *calculatorpb.SumRequest) (*calculatorpb.SumResponse, error) {
	fmt.Printf("Receive SUM RPC: %v", req)
	firstNumber := req.FirstNumber
	secondNumber := req.SecondNumber
	sum := firstNumber + secondNumber
	res := &calculatorpb.SumResponse{
		SumResult: sum,
	}
	return res, nil
}

func (*server) PrimeNumberDecomposition(req *calculatorpb.PrimerNumberDescompositionRequest, stream calculatorpb.CalculatorService_PrimeNumberDecompositionServer) error {
	fmt.Printf("Primer Number Decomposition was called with: %v\n", req)
	N := req.GetNumber()
	k := 2
	for N > 1 {
		if N%int32(k) == 0 {
			res := &calculatorpb.PrimerNumberDescompositionResponse{
				Result: int32(k),
			}
			stream.Send(res)
			time.Sleep(1000 * time.Millisecond)
			N = N / int32(k)
		} else {
			k = k + 1
		}
	}
	return nil
}

func main() {
	fmt.Println("Calculator Server")

	// 50051 defult port for gRPC
	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	// Register new service
	s := grpc.NewServer()
	calculatorpb.RegisterCalculatorServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
