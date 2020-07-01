package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/mrpineapples/go-grpc/calculator/calculatorpb"
	"google.golang.org/grpc"
)

type server struct{}

func (s *server) Sum(ctx context.Context, req *calculatorpb.SumRequest) (*calculatorpb.SumResponse, error) {
	fmt.Printf("Sum function was invoked with: %v\n", req)

	firstNum := req.GetFirstNumber()
	secondNum := req.GetSecondNumber()
	sum := firstNum + secondNum
	res := &calculatorpb.SumResponse{
		Sum: sum,
	}

	return res, nil
}

func (s *server) PrimeDecomposition(req *calculatorpb.PrimeDecompositionRequest, stream calculatorpb.CalculatorService_PrimeDecompositionServer) error {
	fmt.Printf("PrimeDecomposition function was invoked with: %v\n", req)

	divisor := int64(2)
	num := req.GetNumber()
	for num > 1 {
		if num%divisor == 0 {
			res := &calculatorpb.PrimeDecompositionResponse{
				PrimeNumber: divisor,
			}
			stream.Send(res)
			time.Sleep(1000 * time.Millisecond)
			num = num / divisor
		} else {
			divisor++
		}
	}

	return nil
}

func main() {
	fmt.Println("Calculator Server")

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	calculatorpb.RegisterCalculatorServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to server: %v", err)
	}
}
