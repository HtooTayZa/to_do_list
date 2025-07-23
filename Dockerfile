FROM golang:1.22-alpine AS build
WORKDIR /src
COPY . .
RUN go build -o /keep .

FROM alpine:3.20
COPY --from=build /keep /usr/local/bin/keep
ENV KEEP_ADDR=:8787
EXPOSE 8787
ENTRYPOINT ["/usr/local/bin/keep"]
