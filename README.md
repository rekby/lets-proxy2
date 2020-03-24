[![Build Status](https://travis-ci.org/rekby/lets-proxy2.svg?branch=master)](https://travis-ci.org/rekby/lets-proxy2)
[![Coverage Status](https://coveralls.io/repos/github/rekby/lets-proxy2/badge.svg?branch=master)](https://coveralls.io/github/rekby/lets-proxy2?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/rekby/lets-proxy2)](https://goreportcard.com/report/github.com/rekby/lets-proxy2)
[![GolangCI](https://golangci.com/badges/github.com/rekby/lets-proxy2.svg)](https://golangci.com/r/github.com/rekby/lets-proxy2)
[![Mentioned in Awesome Go](https://awesome.re/mentioned-badge.svg)](https://github.com/avelino/awesome-go)  


Русскоязычное описание ниже (Russian below).

English description
===================
Home page: https://github.com/rekby/lets-proxy2

Features:
* http-01 and tls-alpn-01 validation
* HTTPS (with certificate autoissue) and HTTP reverse proxy
* Zero config for start usage
* Time limit for issue certificate
* Auto include subdomains in certificate (default: domain and www.domain)
* Logging for stderr and/or file
* Self rotate log files (can disable by config)
* Can configure backend in dependence of incoming connection IP:Port
* Custom headers to backend
* Self check domain before issue cert (prevent DoS cert issue attack by requests with bad domains)
* Blacklist/whitelist of domains
* Lock certificates (force to use manual issued certificate without internal checks)
* Optional access to internal metrics with Prometheus format

It is next generation of https://github.com/rekby/lets-proxy, rewrited from scratch.

It is improove logging and internal structure for better test coverage and simpler support.
Add config support instead of many command line flags.

A reverse-proxy server to handle https requests transparently. By default Lets-proxy handles
https requests to port 443 and proxies them as http to port 80 on the same IP address.

Lets-proxy adds the http headers, `X-Forwarded-For` which contains the IP address.
It obtains valid TLS certificates from Let's Encrypt and handles https for free, in an automated way, 
including certificate renewal, and without warning in browsers.

The program was created for shared hosting and can handle many thousands of domains per server.
It is simple to implement and doesn't need settings to start the program on personal server/vps.

Quick start:

    ./lets-proxy or lets-proxy.exe
    
Use --help key for details:

    ./lets-proxy --help or lets-proxy.exe --help

Русский (Russian):
==================
Сайт программы: https://github.com/rekby/lets-proxy2

Сейчас это тестовая версия, программа в процессе разработчик и она не готова для реального использования.

Возможности:
* Авторизация доменов по протоколам http-01 and tls-alpn-01
* Проксирование HTTPS (с автовыпуском сертификата) and HTTP
* Начать использование можно без настроек
* Ограничение времени на получение сертификата
* Автоматическое получение сертификата для домена и поддоменов (default: domain and www.domain)
* Вывод логов в файл и/или на стандартный вывод ошибок
* Самостоятельная ротация лог-файлов (отключается в настройках)
* Можно настроить адрес перенаправления запроса в заивисмости от адреса приема запроса.
* Настраиваемые дополнительные заголовки для передачи на внутренний сервер
* Самостоятельная проверка возможности выпуска сертификата перед его запросов (для исключения DoS-атак путем запросов с неправильными доменами)
* Белый/чёрный списки доменов для выпуска сертификатов
* Фиксированный сертификат (возможность использовать самостоятельно полученный сертификат, без внутренних проверок и автообновления)
* Опциональный доступ к внутренним метрикам в формате Prometheus


Эта программа - следующая итерация после https://github.com/rekby/lets-proxy, переписанная с нуля.

Улучшено логирование, внутреннее устройство кода - для упрощения тестирования и поддержки/развития.
Добавлена поддержка файла настроек вместо огромного списка флагов.

Реверс-прокси сервер для прозрачной обработки https-запросов. Для начала использования достаточно просто запустить его на сервере с 
запущенным http-сервером. При этом lets-proxy начнёт слушать порт 433 и передавать запросы на порт 80 с тем же IP-адресом.
К запросу будет добавляться заголовок `X-Forwarded-For` с IP-адресом источника запроса.
Сертификаты для работы https получаются в реальном времени от letsencrypt.org. Это правильные
(не самоподписанные) бесплатные сертификаты, которым доверяют браузеры.

Программа разрабатывается для использования на виртуальном хостинге и может работать с тысячами доменов
на каждом сервере.
С другой стороны она проста и не требует начальных настроек для запуска на персональном сервере.

Быстрый старт:

    ./lets-proxy или lets-proxy.exe


Used libraries (alphabet ordered):
==================================

* http://github.com/gojuno/minimock - for tests
* http://github.com/kardianos/minwinsvc - for run as windows service
* http://github.com/maxatome/go-testdeep - for tests
* http://github.com/mitchellh/gox - for multiply binaries build
* http://github.com/pelletier/go-toml - for config file.
* http://github.com/rekby/zapcontext - for pass logger to/from context
* http://github.com/satori/go.uuid - for generate uuid
* http://go.uber.org/zap - for logging
* http://golang.org/x/crypto - use acme part - for access for lets encrypt server.
* http://golang.org/x/net - use idna part for log domain names
* http://gopkg.in/natefinch/lumberjack.v2 - use for self log rotation.
