#!/bin/bash

cp /etc/nginx/templates/${M8S_UI_TEMPLATE} /etc/nginx/conf.d/default.conf

sed -i -e "s/M8S_UI_PORT/${M8S_UI_PORT}/g" /etc/nginx/conf.d/default.conf
sed -i -e "s/M8S_UI_API_PORT/${M8S_UI_API_PORT}/g" /etc/nginx/conf.d/default.conf
sed -i -e "s/M8S_UI_OAUTH2_PORT/${M8S_UI_OAUTH2_PORT}/g" /etc/nginx/conf.d/default.conf

nginx -g "daemon off;"