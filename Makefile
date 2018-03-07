CC=go build

all: chat_messenger

chat_messenger: main.go

	$(CC) main.go 

clean:

	rm main
