FROM golang:1.13.0 AS builder
ENV GO111MODULE=on

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o app .
# RUN CGO_ENABLED=0 GOOS=linux go build -o app .


FROM alpine
COPY --from=builder /app .

# download wait so we can explicitly tell our app service which services to wait to be healthy in our docker-compose file.
ADD https://github.com/ufoscout/docker-compose-wait/releases/download/2.7.3/wait /wait
# elevate the permissions for wait
RUN chmod +x /wait

CMD /wait && /app
