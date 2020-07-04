package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/mrpineapples/go-grpc/calculator/calculatorpb"
	"google.golang.org/grpc"
)

func main() {
	fmt.Println("Calculator Client")
	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect: %v", err)
	}
	defer cc.Close()

	c := calculatorpb.NewCalculatorServiceClient(cc)
	// doUnary(c)

	// doServerStreaming(c)
	doClientStreaming(c)
}

func doUnary(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("Starting a Unary RPC...")
	req := &calculatorpb.SumRequest{
		FirstNumber:  10,
		SecondNumber: 14,
	}
	res, err := c.Sum(context.Background(), req)
	if err != nil {
		log.Fatalf("Error while calling Sum RPC: %v", err)
	}
	log.Printf("Response from Sum: %v", res.Sum)
}

func doServerStreaming(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("Starting a Server Stream RPC...")

	req := &calculatorpb.PrimeDecompositionRequest{
		Number: 120,
	}
	resStream, err := c.PrimeDecomposition(context.Background(), req)
	if err != nil {
		log.Fatalf("Error while calling PrimeDecomposition RPC: %v", err)
	}

	for {
		msg, err := resStream.Recv()
		if err == io.EOF {
			// reached end of stream
			break
		} else if err != nil {
			log.Fatalf("Error while reading stream: %v", err)
		}

		log.Printf("Response from PrimeDecomposition: %v", msg.GetPrimeNumber())
	}
}

func doClientStreaming(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("Starting a Client Stream RPC...")

	stream, err := c.ComputeAverage(context.Background())
	if err != nil {
		log.Fatalf("Error while reading client stream: %v", err)
	}

	numbers := []int32{1, 2, 3, 4}
	for _, num := range numbers {
		fmt.Printf("Sending req: %v\n", num)
		stream.Send(&calculatorpb.ComputeAverageRequest{
			Number: num,
		})
		time.Sleep(1000 * time.Millisecond)
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Error while receiving response from ComputeAverage: %v", err)
	}

	fmt.Printf("ComputeAverage Response: %v\n", res.GetAverage())
}
