package main

import (
	"context"
	"fmt"
	"log"

	"github.com/go-grpc-udemy-course/go-grpc-course/greet/greetpb"
	"google.golang.org/grpc"
)

func main() {
	fmt.Println("Hello I'm a client")
	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())

	if err != nil {
		log.Fatalf("Could not connect: %v", err)
	}

	defer cc.Close() // Call at the very end

	c := greetpb.NewGreetServiceClient(cc)
	// fmt.Printf("Created client: %f", c)

	doUnary(c)
	doUnaryAddition(c)
}

func doUnary(c greetpb.GreetServiceClient) {
	fmt.Println("Starting to do a Unary RPC...")
	req := &greetpb.GreetRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Jorge",
			LastName:  "Bayona",
		},
	}
	fmt.Printf("Request: %v", req)
	res, err := c.Greet(context.Background(), req)
	if err != nil {
		log.Fatalf("Error while calling Greet RPC: %v", err)
	}
	log.Printf("Response from Greet: %v", res.Result)
}

func doUnaryAddition(c greetpb.GreetServiceClient) {
	req := &greetpb.SumRequest{
		FirstNumber:  3,
		SecondNumber: 10,
	}
	res, err := c.AdditionOperation(context.Background(), req)
	if err != nil {
		log.Fatalf("Error while calling Addition RPC: %v", err)
	}
	log.Printf("Response from Addition RPC: %v", res.Result)
}
