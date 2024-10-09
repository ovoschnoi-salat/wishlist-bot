# syntax=docker/dockerfile:1

# Build the application from source
FROM golang:1.22 AS build-stage

WORKDIR /build

COPY locales ./locales
COPY pg ./pg
COPY repository ./repository
COPY storage ./storage
COPY go.mod go.sum *.go ./
RUN go mod tidy
RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -o /app

# Deploy the application binary into a lean image
#FROM gcr.io/distroless/base-debian11 AS build-release-stage
FROM debian AS build-release-stage

RUN addgroup --gid 1000 groupcontainer
RUN adduser -u 1000 -G groupcontainer -h /home/containeruser -D containeruser

USER containeruser


WORKDIR /home/containeruser

COPY --from=build-stage /app /home/containeruser/app

ENTRYPOINT ["./app"]