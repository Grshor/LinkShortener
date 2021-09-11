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

// не хороший клиент
func main() {
	method := "Create"
	link := "cute link 64"

	log.Printf("Начали")
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("Не удалось подключится: %v", err)
	}
	log.Printf("Подключились")
	defer conn.Close()
	c := pb.NewLinkShortenerClient(conn)

	// берём аргументы запроса
	argsLen := len(os.Args)
	if argsLen > 2 {
		method = os.Args[argsLen-2]
		link = os.Args[argsLen-1]
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
