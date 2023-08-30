FROM golang:1.21-bookworm as build

WORKDIR /go/src/app
ADD . /go/src/app
RUN go build -o /go/bin/app

FROM gcr.io/distroless/base-debian12
COPY --from=build /go/bin/app /
CMD ["/app"]
