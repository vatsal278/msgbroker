FROM golang:1.17 as builder

ENV GO111MODULE=on

WORKDIR /app

COPY . .

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main cmd/main.go

FROM scratch
ARG PORT
EXPOSE $PORT
COPY --from=builder /app/main .

#COPY --from=builder /app/.env .

CMD [ "/main" ]