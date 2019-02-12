FROM golang

WORKDIR /startup
COPY . /startup

RUN go get github.com/bradfitz/gomemcache/memcache
RUN go get github.com/go-sql-driver/mysql
RUN go get google.golang.org/grpc
RUN go get google.golang.org/grpc/reflection
RUN go get github.com/Fattouche/DayTrader/golang/protobuff

WORKDIR /day_trader

ENTRYPOINT ["bash", "/startup/entrypoint.sh"]