GophKeeper

make gen - сгенерировать grpc-структуры из proto-файлов
make build - собрать бинарные файлы клиента и сервера

Генерация сертификатов и ключа сервера:
```
openssl genrsa -out rootCAKey.pem 2048
openssl req -x509 -sha256 -new -nodes -key rootCAKey.pem -days 3650 -out rootCACert.pem
```