    FROM golang:latest
    WORKDIR /app
    COPY . /app/
    RUN apt-get -y update
    RUN apt-get install -y ffmpeg
    RUN go build -o ./out/my-gin-app .
    EXPOSE 8081
    ENTRYPOINT ["./out/my-gin-app"]