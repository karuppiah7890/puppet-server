FROM golang:alpine as build
WORKDIR /opt/puppet-server
ADD . .
RUN go build -o out/puppet-server

FROM alpine
COPY --from=build /opt/puppet-server/out/puppet-server /usr/bin
CMD ["puppet-server"]