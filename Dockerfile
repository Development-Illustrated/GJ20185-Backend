FROM golang:latest 
RUN mkdir /app 
ADD ./main /app/ 
WORKDIR /app 
RUN go get -u github.com/gorilla/mux && go get github.com/googollee/go-socket.io && go get github.com/rs/cors
RUN go build -o main .
EXPOSE 6969
CMD ["/app/main"]
