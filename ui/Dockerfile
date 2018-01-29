FROM nginx:latest

ENV M8S_UI_PORT=80
ENV M8S_UI_TEMPLATE=default.conf
ENV M8S_UI_API_PORT=8080
ENV M8S_UI_OAUTH2_PORT=4180

ADD conf/default.conf /etc/nginx/templates/default.conf
ADD conf/oauth2_proxy.conf /etc/nginx/templates/oauth2_proxy.conf

ADD entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

ADD src/client/public /var/www/html

CMD ["/entrypoint.sh"]
