FROM golang:1.12
WORKDIR /go/src/github.com/nouney/slack-flim
ADD . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o server *.go

FROM drone/ca-certs
ARG APP_PATH
WORKDIR /
COPY --from=0 /go/src/github.com/nouney/slack-flim/server .
ENTRYPOINT ["/server"]
