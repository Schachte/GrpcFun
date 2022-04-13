package main

import (
	"context"
	"fmt"
	"io"
	"log"

	"github.com/schachte/grpc/calculator/calculatorpb"
	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	c := calculatorpb.NewCalculatorClient(conn)
	req := &calculatorpb.SumRequest{
		Val_1: 250,
		Val_2: 150,
	}
	res, err := c.Sum(context.Background(), req)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("The value of adding %d and %d returned %d", req.Val_1, req.Val_2, res.Result)

	primeReq := &calculatorpb.PrimeRequest{
		PrimeTarget: uint32(120),
	}

	prime, err := c.PrimeDecomposition(context.Background(), primeReq)
	if err != nil {
		log.Fatal(err)
	}
	for {
		res, err := prime.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(res.GetValue())
	}
}
