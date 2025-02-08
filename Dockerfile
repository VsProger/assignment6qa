FROM golang:1.20.1-alpine3.16 as base

RUN apk add build-base 
WORKDIR /web
COPY . .
RUN go build -o forum ./cmd/web/

FROM alpine:3.16
WORKDIR /web
COPY --from=base /web/ /web/

CMD ["./forum"]
