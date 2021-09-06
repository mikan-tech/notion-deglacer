# build stage
FROM golang:1.14 AS build
WORKDIR /go/src/project
COPY . .
RUN go get -d
RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux go build -o bin/deglacer ./cmd/deglacer

# final stage
FROM alpine:latest
COPY --from=build /go/src/project/bin/deglacer /project/bin/deglacer
EXPOSE 8080
ENTRYPOINT ["/project/bin/deglacer"]
