FROM golang:1.24.1 as build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY ./ ./

RUN \
CGO_ENABLED=0 GOOS=linux \
go build -ldflags="-s -w" -o jitsi-welcome && \
strip jitsi-welcome


FROM scratch

WORKDIR /app

COPY --from=build /etc/ssl/certs/ca-certificates.crt \
/etc/ssl/certs/ca-certificates.crt

COPY --from=build /app/jitsi-welcome /jitsi-welcome

ENTRYPOINT ["/jitsi-welcome"]