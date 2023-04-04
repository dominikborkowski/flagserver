FROM golang:1.19-alpine as builder
LABEL org.opencontainers.image.authors="dominik.borkowski@gmail.com"

# compile flagserver code
WORKDIR /app
COPY ./flagserver/go.mod ./
RUN go mod download
COPY ./flagserver/*.go ./
RUN go build -o /flagserver

# final container
FROM --platform=$TARGETPLATFORM alpine:3.17

ARG USER=user
ENV HOME /home/$USER
RUN adduser -D $USER

COPY --from=builder /flagserver /usr/local/bin/flagserver
COPY ./flag.txt $HOME/flag.txt

USER $USER
WORKDIR $HOME

ENTRYPOINT ["/usr/local/bin/flagserver","-h", "0.0.0.0","-p", "9999", "-f", "/home/user/flag.txt"]