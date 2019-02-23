FROM golang:alpine as build
WORKDIR /opt/puppet-server
ADD . .
RUN go build

FROM alpine
COPY --from=build /opt/puppet-server/puppet-server /usr/bin
CMD ["puppet-server"]