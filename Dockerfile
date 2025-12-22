FROM golang:1.25-alpine

#dodat pocetni image

WORKDIR /app
#trenutni radni direktorijum naseg imagea

COPY go.mod go.sum ./
RUN go mod download
#preuzima sve zavisnosti naseg servisa
COPY ./ ./

COPY swagger.yaml /app/swagger.yaml
COPY swagger /app/swagger



RUN go build -o app

EXPOSE 8080

CMD ["./app"]
