FROM golang:latest

ENV GOPATH=/
COPY ./ ./

# устанавливаем psql
RUN apt-get update
RUN apt-get -y install postgresql-client

# make wait-for-postgres.sh executable
RUN chmod +x wait-for-postgres.sh

# скачиваем зависимости
RUN go mod download
# компилируем исполняемый файл
RUN go build -o server ./cmd/server.go
CMD ["./server"]