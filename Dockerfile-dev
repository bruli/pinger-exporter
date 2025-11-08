FROM golang:1.25.4

RUN mkdir /app
WORKDIR /app

RUN apt-get update && apt-get upgrade -y && apt-get install -y make git

RUN go install -v github.com/cespare/reflex@latest
EXPOSE 8080
ENTRYPOINT ["reflex", "-c", "./reflex.conf"]
