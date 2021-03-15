package main

import (
	"calculator/calculatorpb"
	"context"
	"fmt"
	"io"
	"log"

	"google.golang.org/grpc"
)

func main() {
	fmt.Println("Hello I am Client")
	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect on %v", err)
	}
	defer cc.Close()
	c := calculatorpb.NewCalculatorServiceClient(cc)
	//fmt.Printf("Created Client: %f", c)
	//doUnary(c)
	//doServerStreaming(c)
	doClientStreaming(c)
}
func doUnary(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("Starting to do a unary RPC...")
	var first, second int32
	fmt.Println("Enter a first number:->")
	fmt.Scanln(&first)
	fmt.Println("Enter a second number:->")
	fmt.Scanln(&second)
	req := &calculatorpb.CalculatorRequest{
		FirstNumber: first,
		LastNumber:  second,
	}
	res, err := c.Sum(context.Background(), req)
	if err != nil {
		log.Fatalf("error while creating Sum: %v", err)
	}
	fmt.Println()
	log.Printf("\n Response from greet:%v", res.Result)
}
func doServerStreaming(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("Starting to do a Server Streaming RPC...")
	var number int32
	fmt.Println("Enter a number:")
	fmt.Scanln(&number)
	req := &calculatorpb.CalculatorManyTimeRequest{
		Number: number,
	}
	resStream, err := c.PrimeDeComposition(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling CalculatorManytimes %v", err)
	}
	log.Printf("Prime Decompostion of %v number", number)
	for {
		msg, err := resStream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("error while reading stream from server side %v ", err)
		}
		log.Printf(" %v ", msg.GetResult())
	}
}
func doClientStreaming(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("Starting to do a ComputeAverage RPC...")
	stream, err := c.ComputeAverage(context.Background())
	if err != nil {
		log.Fatalf("Error while opening stream %v", err)
	}
	numbers := []int32{3, 5, 9, 54, 23}
	for _, number := range numbers {
		fmt.Printf("Sending number: %v", number)
		stream.Send(&calculatorpb.ComputeAverageRequest{
			Number: number,
		})
	}
	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("error while receiving response %v", err)
	}
	fmt.Printf("The Average is :%v", res.GetResult())
}
