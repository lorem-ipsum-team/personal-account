FROM golang:1.23-alpine AS build
WORKDIR /usr/src

COPY go.mod go.sum ./
RUN go mod download -x

COPY cmd ./cmd
COPY internal ./internal
COPY pkg ./pkg

RUN CGO_ENABLED=0 GOOS=linux go build -v -o ./out/service ./cmd/service/main.go


FROM alpine:3 AS run
RUN apk add --no-cache ca-certificates tzdata curl
COPY --from=build /usr/src/out/service /usr/bin/service
CMD [ "/usr/bin/service", "-config", "config.yaml" ]
