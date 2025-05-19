FROM golang:latest

WORKDIR /app

ADD ./ ./

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -o sucurity_app cmd/server.go

ENV APP_PORT=9000
ENV DB_URL="postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"

RUN #echo "$DB_URL"
#RUN curl -s https://packagecloud.io/install/repositories/golang-migrate/migrate/script.deb.sh | bash
#RUN apt-get update && apt-get install -y migrate
#RUN migrate -database "$DB_URL" -path migrations up

CMD ["./security_app"]
