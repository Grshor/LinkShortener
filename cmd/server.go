package main

import (
	"context"
	"github.com/jackc/pgx/v4"
	"linkShortener/pkg"
	"log"
	"net"
	//"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"google.golang.org/grpc"
	pb "linkShortener/pkg/proto"
)

const (
	port = pkg.Port
	//linkStart = pkg.Address
)

//port = pkg

type linkShortenerServer struct {
	//pgConn *pgx.Conn // для обычной (конкурентно небезопасной) pgx сессии
	pgConn *pgxpool.Pool // для пула конкурентно безопасной pgx сессии
	pb.UnimplementedLinkShortenerServer
}

// Get - метод, который принимает сокращённый URL (shortLink), ищет его в базе и возвращает оригинальный (longLink)
func (s *linkShortenerServer) Get(ctx context.Context, in *pb.ShortLink) (link *pb.LongLink, error error) {
	_ = ctx
	shortLink := in.GetLink()
	log.Printf("Принял get: %v", shortLink)

	var longLink string
	err := s.pgConn.QueryRow(context.Background(),
		"select longLink from links where shortLink = $1", shortLink).Scan(&longLink)
	if err != nil {
		return nil, err
	}
	return &pb.LongLink{Link: longLink}, nil
}

// Create - метод, который сохраняет оригинальный URL (longLink) в базе и возвращать сокращённый (shortLink)
// проводит 3 запроса в бд. Можно снизить до 2-х, если использовать обратную кодировку для метода Get, вместо запроса
// на получение longLink через соответствие shortLink.
// Тогда в таблице не нужно будет вообще хранить shortLink, или хранить longLink, это на выбор.
// НЕТ ПРОВЕРКИ longLink на то, реально ли это url.
func (s *linkShortenerServer) Create(ctx context.Context, in *pb.LongLink) (link *pb.ShortLink, error error) {
	_ = ctx
	var shortLink string
	longLink := in.GetLink()
	log.Printf("Принял create: %v", longLink)

	tx, err := s.pgConn.Begin(context.Background()) // открываем транзакцию
	if err != nil {
		//return nil, err // дальше идти невозможно
		log.Fatalf("Ошибка при conn.Begin: %v", err) // заменить на return
	}

	// проверяем longUrl на наличие в бд
	row := tx.QueryRow(context.Background(),
		"select shortLink from links where longLink = $1", longLink).Scan(&shortLink)
	if row != pgx.ErrNoRows { // то-есть действительно нашлась такая запись
		log.Printf("Найден дупликат")
		error = row
		link = &pb.ShortLink{Link: shortLink}
		return
	}

	// пишем новый longLink в бд
	var linkId int // из этого красивого числа мы и получим shortLink
	row = tx.QueryRow(context.Background(),
		"insert into links(longLink, shortLink) values ($1, $2) RETURNING id", longLink, shortLink).Scan(&linkId)
	if row == pgx.ErrNoRows { // хотя какая может быть ошибка на insert?
		log.Fatalf("Ошибка при добавлении записи в базу: %v", err) // заменить на return
	}
	// создаём новый shortLink
	shortLink = pkg.DehydrateAndUpgrade(linkId)
	// добавляем к новом longLink его shortLink
	_, err = tx.Exec(context.Background(),
		"update links SET shortLink = $1 where id = $2", shortLink, linkId)
	if err != nil {
		return nil, err // дальше идти невозможно
		//log.Fatalf("Ошибка при добавлении записи в базу: %v", err) // заменить на return
	}

	tx.Commit(context.Background())
	error = err
	link = &pb.ShortLink{Link: shortLink}
	return
}

func main() {
	log.Printf("Подключаемся к БД")
	pgConn, err := pgxpool.Connect(context.Background(), pkg.DatabaseUrl)
	if err != nil {
		log.Fatalf("Ошибка подключения к бд по %v : %v", pkg.DatabaseUrl, err) // скорее всего порт занят
	}
	createSql := "create table if not exists links(id SERIAL PRIMARY KEY, longLink text, shortLink text);"
	_, err = pgConn.Exec(context.Background(), createSql)
	if err != nil {
		log.Fatalf("Не удалось создать таблицу: %v", err)
	}

	//defer pgConn.Close(context.Background()) // это для обычной (конкурентно небезопасной) pgx сессии
	defer pgConn.Close() // для пула конкурентно безопасной pgx сессии
	s := &linkShortenerServer{pgConn: pgConn}
	if err := s.Run(); err != nil {
		log.Fatalf("Не удалось запустить сервер: %v", err)
	}
}

func (s *linkShortenerServer) Run() error {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Ошибка tcp-listen с портом %v : %v", port, err) // скорее всего порт занят
	}

	server := grpc.NewServer()
	pb.RegisterLinkShortenerServer(server, s)
	log.Printf("Слушаем адрес: %v", lis.Addr())

	return server.Serve(lis)
}
