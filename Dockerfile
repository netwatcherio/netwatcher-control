FROM golang:1.18-bullseye

RUN go install github.com/netwatcherio/netwatcher-control@latest

ENV GO111MODULE=on
ENV GOFLAGS=-mod=vendor

COPY .env /app
ENV APP_HOME /go/src/netwatcher-control
RUN mkdir -p "$APP_HOME"

WORKDIR "$APP_HOME"
EXPOSE 3000
CMD ["bee", "run"]

# TODO