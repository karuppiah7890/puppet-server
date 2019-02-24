FROM golang as build
WORKDIR /opt/puppet-server
ADD . .
RUN go build -o out/puppet-server

FROM ubuntu
COPY --from=build /opt/puppet-server/out/puppet-server /usr/bin
CMD ["puppet-server"]