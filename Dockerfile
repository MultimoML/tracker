FROM golang:alpine AS build

WORKDIR /usr/src/app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .

RUN go install github.com/swaggo/swag/cmd/swag@latest && swag init

RUN CGO_ENABLED=0 go build -o /tracker main.go

FROM gcr.io/distroless/static-debian11:latest

COPY --from=build /tracker /tracker

ENV PORT=6003
EXPOSE $PORT

ENTRYPOINT ["/tracker"]
