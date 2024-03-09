FROM golang:1.21.8-alpine3.19 AS build

LABEL authors="Benyamin Mahmoudyan"

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download && go mod verify

COPY . .

RUN go build -o /app/gameapp


FROM  alpine:3.19 AS RUN

COPY --from=build /app/gameapp /bin/gameapp

EXPOSE 1986

CMD ["/bin/gameapp"]
