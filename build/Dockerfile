FROM golang:1.20.5-alpine as builder

WORKDIR /root/whaletail
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN go build -trimpath -o bin/whaletail cmd/whaletail/*.go


FROM scratch

COPY --from=builder /root/whaletail/bin/whaletail /usr/local/bin/whaletail
ENTRYPOINT [ "whaletail" ]
