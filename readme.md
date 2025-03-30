### Yt-downloader

Скачиватель музыки в формате mp3 из Youtube. Веб инфтерфейс для [yt-dlp](https://github.com/yt-dlp/yt-dlp).

Сейчас уже не работает. Требуется добавить поддержку прокси серверов для возобновления работы.


### Nginx

Настройки сервера nginx для отдельного домена

`user@24641 ~> cat /etc/nginx/sites-available/yt`

```
server {

    server_name yt.daybydayz.ru;

    location / {
        proxy_pass              http://127.0.0.1:10542;
        proxy_set_header        Host $host;
        proxy_set_header        X-Real-IP $remote_addr;
        proxy_set_header        X-Forwarded-For $proxy_add_x_forwarded_for;
    }

    # включаем поддержку http2
    http2 on;
    # включаем поддержку http3
    http3 on;
    # разрешаем GSO
    quic_gso on;
    # разрешаем проверку адреса
    quic_retry on;
    # для перенаправления браузеров в quic-порт
    add_header Alt-Svc 'h3=":443";max=86400' always;
    add_header Alt-Svc 'h3=":443"; ma=2592000, h3-29=":443"; ma=2592000, h3-Q050=":443"; ma=2592000, h3-Q046=":443"; ma=2592000, h3-Q043=":443"; ma=2592000, quic=":443"; ma=2592000;' always;

    listen 443; # managed by Certbot
    ssl_certificate /etc/letsencrypt/live/yt.daybydayz.ru/fullchain.pem; # managed by Certbot
    ssl_certificate_key /etc/letsencrypt/live/yt.daybydayz.ru/privkey.pem; # managed by Certbot
    include /etc/letsencrypt/options-ssl-nginx.conf; # managed by Certbot
    ssl_dhparam /etc/letsencrypt/ssl-dhparams.pem; # managed by Certbot

}

```
