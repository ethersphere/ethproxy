FROM golang:1.18 AS build

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download
COPY . ./

RUN make binary

FROM debian:11.5-slim

ENV DEBIAN_FRONTEND noninteractive

RUN apt-get update && apt-get install -y --no-install-recommends \
    ca-certificates; \
    apt-get clean; \
    rm -rf /var/lib/apt/lists/*; 


COPY --from=build /src/dist/proxy /usr/local/bin/proxy

EXPOSE 6000 6100

ENTRYPOINT ["proxy"]
