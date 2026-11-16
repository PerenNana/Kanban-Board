#--------------- Stage 1: Build the Go application ---------------
FROM golang:1.27-alpine AS build

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -o /out/server .

#--------------- Stage 2: Run the final image ---------------
FROM alpine:3

WORKDIR /app

COPY --from=build /out/server ./server
COPY templates ./templates

EXPOSE 80

CMD ["./server"]
