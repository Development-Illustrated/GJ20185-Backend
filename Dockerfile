FROM golang:latest 
RUN mkdir /app 
ADD ./main /app/ 
WORKDIR /app 
RUN go get github.com/gorilla/mux && go get -u github.com/rs/cors && go get -u github.com/gorilla/websocket
RUN go build -o main .
EXPOSE 6969
EXPOSE 8000
CMD ["/app/main"]
