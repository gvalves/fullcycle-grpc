package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/gvalves/fullcycle-grpc/pb"
	"google.golang.org/grpc"
)

func main() {

	connection, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not create connection: %v", err)
	}
	defer connection.Close()

	client := pb.NewUserServiceClient(connection)
	// AddUser(client)
	// AddUserVerbose(client)
	// AddUsers(client)
	AddUserStreamBoth(client)

}

func AddUser(client pb.UserServiceClient) {

	req := &pb.User{
		Id:    "0",
		Name:  "João",
		Email: "j@j.com",
	}

	res, err := client.AddUser(context.Background(), req)
	if err != nil {
		log.Fatalf("Could not make request: %v", err)
	}

	fmt.Println(res)

}

func AddUserVerbose(client pb.UserServiceClient) {

	req := &pb.User{
		Id:    "0",
		Name:  "João",
		Email: "j@j.com",
	}

	responseStream, err := client.AddUserVerbose(context.Background(), req)
	if err != nil {
		log.Fatalf("Could not make request: %v", err)
	}

	for {
		stream, err := responseStream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Could not receive msg: %v", err)
		}
		fmt.Println("Status:", stream.Status, "-", stream.GetUser())
	}

}

func AddUsers(client pb.UserServiceClient) {

	reqs := []*pb.User{
		{
			Id:    "1",
			Name:  "Fulano 1",
			Email: "fulano1@gmail.com",
		},
		{
			Id:    "2",
			Name:  "Fulano 2",
			Email: "fulano2@gmail.com",
		},
		{
			Id:    "3",
			Name:  "Fulano 3",
			Email: "fulano3@gmail.com",
		},
		{
			Id:    "4",
			Name:  "Fulano 4",
			Email: "fulano4@gmail.com",
		},
		{
			Id:    "5",
			Name:  "Fulano 5",
			Email: "fulano5@gmail.com",
		},
	}

	stream, err := client.AddUsers(context.Background())
	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}

	for _, req := range reqs {
		stream.Send(req)
		time.Sleep(time.Second * 3)
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Error receiving response: %v", err)
	}

	fmt.Println(res)

}

func AddUserStreamBoth(client pb.UserServiceClient) {

	stream, err := client.AddUserStreamBoth(context.Background())
	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}

	reqs := []*pb.User{
		{
			Id:    "1",
			Name:  "Fulano 1",
			Email: "fulano1@gmail.com",
		},
		{
			Id:    "2",
			Name:  "Fulano 2",
			Email: "fulano2@gmail.com",
		},
		{
			Id:    "3",
			Name:  "Fulano 3",
			Email: "fulano3@gmail.com",
		},
		{
			Id:    "4",
			Name:  "Fulano 4",
			Email: "fulano4@gmail.com",
		},
		{
			Id:    "5",
			Name:  "Fulano 5",
			Email: "fulano5@gmail.com",
		},
	}

	wait := make(chan int)

	go func() {
		for _, req := range reqs {
			fmt.Println("Sending user:", req.Name)
			stream.Send(req)
			time.Sleep(time.Second * 2)
		}
		stream.CloseSend()
	}()

	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("Error receiving data: %v", err)
				break
			}
			fmt.Printf("Recebendo user %v com status: %v\n", res.GetUser(), res.GetStatus())
		}
		close(wait)
	}()

	<-wait

}
