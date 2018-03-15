FROM nginx:1.13

ENV HTTP_PORT=8888

RUN apt-get update && \
    apt-get install -y apache2-utils && \
    apt-get clean

RUN rm /etc/nginx/conf.d/default.conf
ADD nginx.conf /etc/nginx/nginx.conf

ADD entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

CMD ["/entrypoint.sh"] 
