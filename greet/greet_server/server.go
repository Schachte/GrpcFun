package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"time"

	"github.com/schachte/grpc/greet/greetpb"
	"google.golang.org/grpc"
)

const PORT = 50051

type server struct{}

func (*server) Greet(ctx context.Context, req *greetpb.GreetRequest) (*greetpb.GreetResponse, error) {
	firstName := req.GetGreeting().GetFirstName()
	result := fmt.Sprintf("Hello %s", firstName)
	res := greetpb.GreetResponse{
		Result: result,
	}
	return &res, nil
}

func (*server) LongGreet(stream greetpb.GreetService_LongGreetServer) error {
	var result string
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			stream.SendMsg(&greetpb.LongGreetResponse{
				Result: result,
			})
			break
		}
		if err != nil {
			log.Fatalf("Error while reading client stream: %v", err)
		}
		result += fmt.Sprintf("%s, ", req.GetGreeting().GetFirstName())
		fmt.Println(result)
	}
	return stream.SendAndClose(nil)
}

func (*server) GreetEveryone(stream greetpb.GreetService_GreetEveryoneServer) error {
	for {
		res, err := stream.Recv()
		fName := res.GetGreeting().GetFirstName()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			log.Fatal(err)
		}
		if err := stream.Send(&greetpb.GreetEveryoneResponse{
			Result: fName + ", welcome!",
		}); err != nil {
			log.Fatal(err)
		}
	}
}

func (*server) GreetManyTimes(req *greetpb.GreetManyTimesRequest, stream greetpb.GreetService_GreetManyTimesServer) error {
	firstName := req.GetGreeting().GetFirstName()
	for i := 0; i < 10; i++ {
		res := &greetpb.GreetManyTimesResponse{
			Result: fmt.Sprintf("Hello %s - iteration %d", firstName, i),
		}
		stream.Send(res)
		time.Sleep(1 * time.Second)
	}
	return nil
}

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", PORT))
	fmt.Printf("Running server on port %d", PORT)

	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	// Create a server and register different services to it
	s := grpc.NewServer()
	greetpb.RegisterGreetServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to server: %v", err)
	}
}
