FROM golang:1.17-alpine as builder

WORKDIR /opt/action
COPY . .

RUN CGO_ENABLED=0 go build -ldflags '-s -w' -o bin/ ./prepare-img-tag


FROM scratch as release

WORKDIR /opt/action

COPY --from=builder /opt/action/bin .

ENTRYPOINT [ "/opt/action/prepare-img-tag" ]
