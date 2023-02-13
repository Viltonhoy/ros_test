# syntax=docker/dockerfile:1
ARG GO_VERSION=1.17

FROM golang:${GO_VERSION}-alpine AS builder
ENV GO111MODULE=on
WORKDIR /src      
COPY go.mod ./
COPY . .
RUN CGO_ENABLED=0 go build \
 -installsuffix `static` \
 -v -o /app cmd/server/main.go

FROM gcr.io/distroless/static-debian11
WORKDIR /
COPY --from=builder /app /app
EXPOSE 8080 8080
USER nobody:nobody
ENTRYPOINT ["/app"]