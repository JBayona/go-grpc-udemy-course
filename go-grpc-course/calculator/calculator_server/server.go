package main

import (
	"context"
	"fmt"
	"io"
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
	k := int32(2)
	for N > 1 {
		if N%k == 0 {
			res := &calculatorpb.PrimerNumberDescompositionResponse{
				Result: k,
			}
			stream.Send(res)
			time.Sleep(1000 * time.Millisecond)
			N = N / k
		} else {
			k = k + 1
		}
	}
	return nil
}

func (*server) ComputeAverage(stream calculatorpb.CalculatorService_ComputeAverageServer) error {
	fmt.Printf("ComputeAverage function was invoked with a streaming request\n")
	result := int32(0)
	count := 0
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(
				&calculatorpb.ComputerAverageResponse{
					Result: float64(result) / float64(count),
				},
			)
		}
		if err != nil {
			log.Fatalf("error while reading client stream: %v", err)
		}
		result += req.GetNumber()
		count++
	}
}

func (*server) FindMaximum(stream calculatorpb.CalculatorService_FindMaximumServer) error {
	fmt.Printf("FindMaximum was invoked with a streaming request\n")
	max := int32(0)
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			log.Fatalf("Error while reading client stream: %v", err)
			return err
		}
		number := req.GetNumber()
		if number > max {
			max = number
			sendErr := stream.Send(
				&calculatorpb.FindMaximumResponse{
					Reult: max,
				},
			)
			if sendErr != nil {
				log.Fatalf("Error while streaming to client: %v", sendErr)
			}
		}
	}
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
