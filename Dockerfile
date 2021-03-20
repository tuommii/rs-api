FROM golang:alpine as BUILD-STEP

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

RUN mkdir /build
WORKDIR /build

COPY . ./
RUN go build -o run cmd/rs/*.go

WORKDIR /dist

RUN cp /build/run .


FROM scratch

COPY --from=BUILD-STEP /build/run /
COPY .env .env

ENTRYPOINT ["/run"]
