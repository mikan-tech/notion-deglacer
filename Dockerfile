# build stage
FROM golang:1.14 AS build
WORKDIR /go/src/project
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o bin/deglacer ./cmd/deglacer

# final stage
FROM alpine:latest
COPY --from=build /go/src/project/bin/deglacer /project/bin/deglacer
EXPOSE 8000
ENTRYPOINT ["/project/bin/deglacer"]
