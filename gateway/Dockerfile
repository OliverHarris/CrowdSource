FROM golang:1.8-alpine3.6
WORKDIR /go/src
# install git
RUN apk --update add \
	git openssl \
	&& rm /var/cache/apk/*



ADD . /go/src

RUN go-wrapper download

CMD ["go", "run", "main.go"]

EXPOSE 80
