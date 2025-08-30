FROM golang:latest
LABEL maintainer="Ryan S"

RUN apt-get update -qq

RUN apt-get install -y -qq libtesseract-dev libleptonica-dev
ENV TESSDATA_PREFIX=/usr/share/tesseract-ocr/5/tessdata/
# ENV GO111MODULE=on

# Load languages.
RUN apt-get install -y -qq \
  tesseract-ocr-eng \
  tesseract-ocr-deu \
  tesseract-ocr-jpn

WORKDIR /app

COPY src ./
RUN go get -t github.com/otiai10/gosseract/v2
RUN go mod tidy
# RUN go test -v ./...

RUN go build -o receipt_manager

CMD [./receipt_manager]