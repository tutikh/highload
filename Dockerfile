FROM debian:stretch

WORKDIR /root

RUN apt-get update \
    && DEBIAN_FRONTEND=noninteractive apt-get install -y sqlite3 wget


# Устанавливаем Go, создаем workspace и папку проекта
RUN wget https://storage.googleapis.com/golang/go1.9.linux-amd64.tar.gz &&\
	tar -C /usr/local -xzf go1.9.linux-amd64.tar.gz && \
    mkdir go && mkdir go/src && mkdir go/bin && mkdir go/pkg && \
    mkdir go/src/dumb

RUN apt-get update
RUN apt-get install unzip gcc -y

ENV PATH=${PATH}:/usr/local/go/bin GOROOT=/usr/local/go GOPATH=/root/go

# Копируем наш исходный main.go внутрь контейнера, в папку go/src/dumb
ADD load/data /tmp/data/
RUN ls /tmp/data/
ADD trav.db go/src/highload/hl/load/
ADD .. /  go/src/highload/
ADD ../src /root/go/src/

WORKDIR /root/go/src/highload/

# Компилируем и устанавливаем наш сервер
RUN go get
RUN go build


#RUN unzip -o /root/go/src/highload/hl/load/data/data.zip -d /root/go/src/highload/hl/load/data

# Открываем 80-й порт наружу
EXPOSE 80


CMD ./highload
