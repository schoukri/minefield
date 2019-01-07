FROM golang:1.11

WORKDIR /go/src/minefield

COPY *.go /go/src/minefield/
COPY *.txt /go/src/minefield/

RUN go get -d -v ./...
RUN go install -v ./...

ENTRYPOINT ["minefield"]
