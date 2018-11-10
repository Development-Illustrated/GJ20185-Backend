FROM golang:latest 
RUN mkdir /app 
ADD ./main /app/ 
WORKDIR /app 
RUN go get -u github.com/gorilla/mux
RUN go build -o main . 
CMD ["/app/main"]
