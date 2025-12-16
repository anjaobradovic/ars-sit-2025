FROM golang:1.25-alpine

#dodat pocetni image

WORKDIR /app
#trenutni radni direktorijum naseg imagea

COPY go.mod go.sum ./
RUN go mod download
#preuzima sve zavisnosti naseg servisa
COPY ./ ./

RUN go build -o config-service

EXPOSE 8080

CMD ["./config-service"]
