# This is a multi-stage build. First we are going to compile and then
# create a small image for runtime.
FROM golang:1.11.1 as builder

RUN mkdir -p /build

ENV GO111MODULE on

WORKDIR /build
RUN useradd -u 10001 app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

FROM scratch

COPY --from=builder /build/main /main
COPY --from=builder /build/static /static
COPY --from=builder /etc/passwd /etc/passwd
USER app

EXPOSE 3000
CMD ["/main"]
