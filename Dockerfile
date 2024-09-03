FROM golang:1.19-alpine AS build

ARG VERSION

LABEL description="An incredibly fast proxy checker & IP rotator with ease."
LABEL repository="https://github.com/kitabisa/mubeng"
LABEL maintainer="dwisiswant0"

WORKDIR /app
COPY ./go.mod .
RUN go mod download

COPY . .
RUN go build -ldflags "-s -w -X github.com/kitabisa/mubeng/common.Version=${VERSION}" \
	-o ./bin/mubeng ./cmd/mubeng 

FROM alpine:latest

COPY --from=build /app/bin/mubeng /bin/mubeng
ENV HOME /
ENTRYPOINT ["/bin/mubeng"]
