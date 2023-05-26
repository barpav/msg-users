FROM golang:1.20-alpine
WORKDIR /usr/src/app
COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY cmd/ cmd/
COPY internal/ internal/
RUN go build -v -o /usr/local/bin/app ./cmd
CMD ["app"]