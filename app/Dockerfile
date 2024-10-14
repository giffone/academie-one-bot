# academ_stats
FROM golang:1.22.5

ENV TZ=Asia/Almaty

WORKDIR /app

COPY . .

RUN go mod download
RUN go build -o academie_one_bot_app cmd/ac_bot/main.go

EXPOSE 80

CMD ["./academie_one_bot_app"]