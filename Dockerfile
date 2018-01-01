FROM golang:alpine as builder
RUN apk update && apk add --update git
RUN mkdir /go/src/app
ADD . /go/src/app/
WORKDIR /go/src/app
ENV CGO_ENABLED=0 GOOS=linux
RUN go get -v ./...
RUN cd cmd && go build -o /mail-proxy-service

FROM alpine
RUN apk update && apk add --update ca-certificates
COPY --from=builder /mail-proxy-service /bin/
CMD /bin/mail-proxy-service -primaryProvider $PRIMARY_PROVIDER -privateAPIKey $PRIVATE_API_KEY -publicAPIKey $PUBLIC_API_KEY -senderDomain $SENDER_DOMAIN -serverAddr :$PORT
