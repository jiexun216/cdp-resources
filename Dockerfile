FROM alpine:latest

ADD cdp-resources /cdp-resources
ENTRYPOINT ["./cdp-resources"]