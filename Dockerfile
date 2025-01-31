FROM ubuntu:latest
LABEL authors="hero"

ENTRYPOINT ["top", "-b"]