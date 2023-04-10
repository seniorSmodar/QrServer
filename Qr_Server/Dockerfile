FROM golang:1.19.3

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

RUN go build -o /qr_gl_serv

EXPOSE 6000

CMD [ "/qr_gl_serv" ]