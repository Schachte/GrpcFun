package main

import (
	"context"
	"fmt"
	"io"
	"log"

	"github.com/schachte/grpc/greet/greetpb"
	"google.golang.org/grpc"
)

func main() {
	// SSL is enabled by default, so we are forcing WithInsecure until TLS/SSL is enabled
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect: %v", err)
	}
	defer conn.Close()
	c := greetpb.NewGreetServiceClient(conn)

	res, err := doUnary(c)
	if err != nil {
		log.Fatalf("Error on unary: %v", err)
	}
	fmt.Printf("The result from the unary call is: %s\n", res.Result)

	// doServerStreaming(c)
	// doClientStreaming(c)
	doBidiStreaming(c)
}

func doBidiStreaming(c greetpb.GreetServiceClient) {
	waitCh := make(chan struct{})
	requests := []*greetpb.GreetEveryoneRequest{
		{
			Greeting: &greetpb.Greeting{
				FirstName: "p1",
			},
		},
		{
			Greeting: &greetpb.Greeting{
				FirstName: "p2",
			},
		},
		{
			Greeting: &greetpb.Greeting{
				FirstName: "p3",
			},
		},
	}
	stream, err := c.GreetEveryone(context.Background())

	if err != nil {
		log.Fatal(err)
	}

	go func() {
		for _, val := range requests {
			err := stream.Send(val)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println("Sent!")
		}
		stream.CloseSend()
	}()

	go func() {
		for {
			data, err := stream.Recv()
			if err == io.EOF {
				close(waitCh)
				break
			}
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(data.Result)
		}
	}()

	// Block until everything is done
	<-waitCh
}

func doClientStreaming(c greetpb.GreetServiceClient) {

	stream, err := c.LongGreet(context.Background())
	if err != nil {
		log.Fatalf("Error w/ long greet: %v", err)
	}
	greetings := []*greetpb.LongGreetRequest{
		{
			Greeting: &greetpb.Greeting{
				FirstName: "firstName_1",
			},
		},
		{
			Greeting: &greetpb.Greeting{
				FirstName: "firstName_2",
			},
		},
		{
			Greeting: &greetpb.Greeting{
				FirstName: "firstName_3",
			},
		},
	}

	for _, val := range greetings {
		err := stream.Send(val)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Sent!")
	}
	stream.CloseAndRecv()
}

// doUnary shows how single req/resp is implemented in GRPC
func doUnary(c greetpb.GreetServiceClient) (*greetpb.GreetResponse, error) {
	req := greetpb.GreetRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Ryan",
			LastName:  "Schachte",
		},
	}

	resp, err := c.Greet(context.Background(), &req)
	if err != nil {
		log.Printf("There was a problem: %v", err)
		return &greetpb.GreetResponse{}, err
	}

	return resp, nil
}

func doServerStreaming(c greetpb.GreetServiceClient) {
	req := greetpb.GreetManyTimesRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Ryan",
			LastName:  "Schachte",
		},
	}

	resStream, err := c.GreetManyTimes(context.Background(), &req)
	if err != nil {
		log.Printf("There was a problem: %v", err)
		return
	}

	for {
		msg, err := resStream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Error receiving from stream: %v", err)
		}
		log.Printf("The result received is %s", msg.GetResult())
	}
}
