package main

import (
	"calculator/calculatorpb"
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"time"

	"google.golang.org/grpc"
)

type server struct {
	calculatorpb.UnimplementedCalculatorServiceServer
}

func (*server) Sum(ctx context.Context, req *calculatorpb.CalculatorRequest) (res *calculatorpb.CalculatorResponse, err error) {
	fmt.Printf("Sum function was invoked %v", req)
	firstNumber := req.GetFirstNumber()
	lastNumber := req.GetLastNumber()
	sum := firstNumber + lastNumber
	res = &calculatorpb.CalculatorResponse{
		Result: sum,
	}
	return res, nil
}
func (*server) PrimeDeComposition(req *calculatorpb.CalculatorManyTimeRequest, stream calculatorpb.CalculatorService_PrimeDeCompositionServer) error {
	fmt.Printf("PrimeDecompostion function was invoked %v", req)
	n := req.GetNumber()
	var k int32
	k = 2
	for n > 1 {
		if n%k == 0 {
			res := &calculatorpb.CalculatorManyTimeResponse{
				Result: k,
			}
			stream.Send(res)
			time.Sleep(1000 * time.Millisecond)
			n = n / k
		} else {
			k = k + 1
		}
	}
	return nil
}
func (*server) ComputeAverage(stream calculatorpb.CalculatorService_ComputeAverageServer) error {
	fmt.Printf("Received ComputeAverage RPC \n")
	sum := int32(0)
	count := 0
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			average := float64(sum) / float64(count)
			return stream.SendAndClose(&calculatorpb.ComputeAverageResponse{
				Result: average,
			})
		}
		if err != nil {
			log.Fatalf("error while reading client stream %v", err)
		}
		sum += req.GetNumber()
		count++
	}
}
func main() {
	fmt.Println("Hello I am Server:")
	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to Listem on %v :", err)
	}
	s := grpc.NewServer()
	calculatorpb.RegisterCalculatorServiceServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed To serve %v :", err)
	}
}
