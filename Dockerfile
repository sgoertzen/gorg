FROM golang:1.7
RUN go get golang.org/x/oauth2
RUN go get github.com/stretchr/testify/assert
COPY ./ /go/src/github.com/sgoertzen/gorg/
RUN cd /go/src/github.com/sgoertzen/gorg/ && go get -t ./...
RUN go build /go/src/github.com/sgoertzen/gorg/*.go