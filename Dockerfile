FROM golang:1.23-alpine3.20 as DEPS

WORKDIR /app

COPY ./go.* .

RUN go mod download

FROM golang:1.23-alpine3.20 as BUILDER

ARG GIT_SHA

WORKDIR /app

COPY --from=DEPS /go/pkg /go/pkg
COPY . .

RUN go build \
  -o notifier-server \
  -ldflags="-s -w -X 'github.com/reynn/notifier/internal/constants.AppVersion=$GIT_SHA'" \
  -trimpath \
  cmd/server/main.go

FROM alpine:3.20

COPY --from=BUILDER /app/notifier-server /app

ENTRYPOINT [ "/app" ]

