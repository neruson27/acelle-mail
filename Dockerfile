# BUILD ENVIRONMENT
# -----------------
FROM golang:1.16-alpine as build_environment

WORKDIR /src
ENV CGO_ENABLED=0
COPY go.* ./
RUN apk add upx
RUN go mod download
COPY . .
RUN  go build -ldflags "-s -w -extldflags '-static'" -o acelle-mail && upx ./acelle-mail


# DEPLOYMENT ENVIRONMENT
# -----------------
FROM alpine

RUN apk update && apk add --no-cache bash
WORKDIR /app
COPY --from=build_environment /src/acelle-mail /app/

ENTRYPOINT ["./acelle-mail"]