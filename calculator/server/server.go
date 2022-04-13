package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/schachte/grpc/calculator/calculatorpb"
	"google.golang.org/grpc"
)

type server struct{}

func main() {
	fmt.Println("Server")

	listen, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Error listening to server on port 50051: %v", err)
	}

	defer listen.Close()
	grpcServer := grpc.NewServer()

	calculatorpb.RegisterCalculatorServer(grpcServer, &server{})

	if err := grpcServer.Serve(listen); err != nil {
		log.Fatalf("Failed to serve the listener: %v", err)
	}
}

func (*server) Sum(c context.Context, req *calculatorpb.SumRequest) (*calculatorpb.SumResponse, error) {
	val1, val2 := req.GetVal_1(), req.GetVal_2()
	sum := val1 + val2
	resp := &calculatorpb.SumResponse{
		Result: sum,
	}
	return resp, nil
}

func (*server) PrimeDecomposition(req *calculatorpb.PrimeRequest, stream calculatorpb.Calculator_PrimeDecompositionServer) error {
	interestedVal := req.GetPrimeTarget()
	values := retrievePrimeDecompositionValues(int(interestedVal))
	for _, val := range values {
		stream.Send(&calculatorpb.PrimeResponse{
			Value: uint32(val),
		})
		time.Sleep(1 * time.Second)
	}
	return nil
}

func retrievePrimeDecompositionValues(target int) []int {
	var data []int
	k := 2
	N := target
	for N > 1 {
		if N%k == 0 { // if k evenly divides into N
			data = append(data, k)
			N = N / k // divide N by k so that we have the rest of the number left.
			continue
		}
		k = k + 1
	}
	return data
}
