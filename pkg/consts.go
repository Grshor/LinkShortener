package pkg

const (
	Port    = ":5000"
	Address = "localhost:5000"

	// подключения "по обычному"
	//DatabaseUrl = "postgres://postgres:LX3ZF4M3bAAM4eeeRRwFqJfkwUMbDRHR@0.0.0.0:5432/postgres"

	// подключения, если сервер запущен внутри docker контейнера
	DatabaseUrl = "postgres://postgres:LX3ZF4M3bAAM4eeeRRwFqJfkwUMbDRHR@db:5432/postgres"
)
