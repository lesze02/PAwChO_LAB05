#ETAP 1
##Użyte polecenia
Kompilowanie pliku w WSL:
root@Szymon:~# cp /mnt/c/Users/szyme/Documents/Dockerlab5/main.go /home
root@Szymon:/home# CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o goapp main.go
root@Szymon:/home# cp goapp /mnt/c/Users/szyme/Documents/DockerLab5/

docker build -t lab5_obraz .
docker run -d -p 8080:8080 --name lab5_kontener lab5_obraz

Polecenie uruchamiające serwer "/goapp"

##Polecenia potwierdzające poprawne działanie

PS C:\Users\szyme\Documents\DockerLab5> docker logs lab5_kontener
2025/04/03 11:40:33 Starting app on :8080
2025/04/03 11:41:05 path: /
2025/04/03 11:41:05 path: /favicon.ico
2025/04/03 11:42:34 Starting app on :8080
2025/04/03 11:42:38 path: /

C:\Windows\System32>curl localhost:8080
Hostname: c9bf6eee36fd
Version: 1.0
IP: 172.17.0.2

##Zawartość pliku Dockerfile
FROM scratch

ADD alpine-minirootfs-3.21.3-x86_64.tar.gz /

ARG VERSION=1.0
ENV VERSION=$VERSION

COPY goapp /goapp

CMD ["/goapp"]


#ETAP 2
##Użyte polecenia
docker build -t lab5_obraz2 .
docker run -d -p 80:80 --name lab5_kontener2 lab5_obraz2

##Zawartość pliku Dockerfile
FROM scratch AS builder

ADD alpine-minirootfs-3.21.3-x86_64.tar.gz /
ARG VERSION=1.0
ENV VERSION=$VERSION

COPY goapp /goapp

FROM nginx:alpine

COPY --from=builder /goapp /usr/local/bin/goapp

COPY nginx.conf /etc/nginx/nginx.conf

CMD /usr/local/bin/goapp & nginx -g 'daemon off;'

HEALTHCHECK --interval=10s --timeout=5s \
  CMD curl -f http://localhost || exit 1

##Zawartosc pliku nginx.conf
events {}

http {
    server {
    listen 80;
    server_name localhost;

    location / {
        proxy_pass http://127.0.0.1:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
}