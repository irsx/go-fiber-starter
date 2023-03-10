FROM golang:1.20-alpine

# Set necessary environmet variables needed for our image
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# Move to working directory /build
WORKDIR /build

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN go build -o main .

WORKDIR /dist

RUN cp /build/main .
RUN mkdir -p /dist/app/controllers
RUN mkdir -p /dist/storage/logs
RUN mkdir -p /dist/storage/uploads

EXPOSE 8000

CMD ["/dist/main"]

