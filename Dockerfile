FROM golang:1.17-buster as builder

WORKDIR /home

COPY ./ .

RUN apt-get install --no-install-recommends git=1:2.20.1-2+deb10u3 -y
RUN GIT_COMMIT=$(git rev-list -1 HEAD); \
    CGO_ENABLED=0 CGOOS=linux GOARCH=amd64 \
    go build -o api -ldflags "-X main.version=$GIT_COMMIT" cmd/api/main.go

FROM gcr.io/distroless/static
WORKDIR /app/
COPY --from=builder /home/api .

USER nonroot
ENTRYPOINT [ "./api" ]