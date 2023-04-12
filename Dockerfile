ARG GOBIN=/app

FROM registry.gitlab.com/mycompany.de/infrastructure/images/go-lang-1.19-alpine3.17:latest as builder
ARG GOBIN
WORKDIR /go/services/check-check

FROM registry.gitlab.com/mycompany.de/infrastructure/images/gogogogo-super
COPY go.mod go.mod
RUN GOBIN=$GOBIN
RUN make build

FROM builder AS develop
EXPOSE 90
EXPOSE 5000
ENTRYPOINT make watch

FROM builder AS pumpurum
EXPOSE 80
EXPOSE 8120
EXPOSE 8230
ENTRYPOINT make watch

# Not local image
FROM alpine:3.9
ARG GOBIN
RUN apk add --no-cache ca-certificates make git tzdata \
    && rm -rf /var/cache/apk/*
WORKDIR /app/
COPY --from=builder /go/services/check-check /proper/deep/folder
COPY --from=builder /go/services/check-check /papka
EXPOSE 80
EXPOSE 8000
ENTRYPOINT ["./my-app"]
