package main

import (
	"context"
	"linkShortener/pkg"
	"log"
	"net"

	"google.golang.org/grpc"
	pb "linkShortener/pkg/proto"
)

const (
	port = pkg.Port
)

//port = pkg

type linkShortenerServer struct {
	pb.UnimplementedLinkShortenerServer
}

// CreateShortLink
// хэширует строку (longLink), добавляет хэш в базу (если его нет), возвращает хэш (shortLink)
func CreateShortLink(longLink string) (shortLink string, error error) {
	error = nil
	hashed := pkg.HashFnv(longLink)
	log.Printf("Сгенерировал: %v", hashed)
	return hashed, error
}

// Create - метод, который сохраняет оригинальный URL (longLink) в базе и возвращать сокращённый (shortLink)
func (s *linkShortenerServer) Create(ctx context.Context, in *pb.LongLink) (link *pb.ShortLink, error error) {
	_ = ctx

	longLink := in.GetLink()
	log.Printf("Принял: %v", longLink)
	shortLink, err := CreateShortLink(longLink)

	error = err
	link = &pb.ShortLink{Link: shortLink}
	return
	//return &pb.ShortLink{Link: shortLink}, nil
}

func main() {
	lis, err := net.Listen("tcp", port)

	if err != nil {
		log.Fatalf("Ошибка tcp-listen с портом %v : %v", port, err) // скорее всего порт занят
	}

	server := grpc.NewServer()
	pb.RegisterLinkShortenerServer(server, &linkShortenerServer{})
	log.Printf("Слушаем адрес: %v", lis.Addr())

	if err := server.Serve(lis); err != nil {
		log.Fatalf("%v", err)
	}
}
