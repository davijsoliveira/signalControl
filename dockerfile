# Use a imagem base do Golang
FROM golang:1.16

# Defina o diretório de trabalho dentro do contêiner
WORKDIR /app

# Copie o código fonte para o diretório de trabalho
COPY . .

# Compile o código Go
RUN go build -o main signalControl


# Execute o aplicativo quando o contêiner for iniciado
CMD ["./main"]