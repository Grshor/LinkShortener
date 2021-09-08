package main

import (
	"context"
	"linkShortener/pkg"
	"log"
	"os"
	"time"

	"google.golang.org/grpc"
	pb "linkShortener/pkg/proto"
)

const (
	address = pkg.Address
)

func main() {
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("Не удалось подключится: %v", err)
	}

	defer func(conn *grpc.ClientConn) {
		err := conn.Close()
		if err != nil {
			os.Exit(1)
		}
	}(conn)

	c := pb.NewLinkShortenerClient(conn)
	method := "Create"
	link := "https://www.google.com/"
	if len(os.Args) > 2 {
		method = os.Args[1]
		link = os.Args[2]
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	switch method {
	case "Create":
		r, e := c.Create(ctx, &pb.LongLink{Link: link})
		if e != nil {
			err = e
			break
		}
		log.Printf("Ответ: %s", r.GetLink())

	case "Get":
		r, e := c.Get(ctx, &pb.ShortLink{Link: link})
		if e != nil {
			err = e
			break
		}
		log.Printf("Ответ: %s", r.GetLink())

	default:
		log.Fatalf("Метод не поддерживается")
	}

	if err != nil {
		log.Fatalf("Ошибка сервера: %v", err)
	}
}
