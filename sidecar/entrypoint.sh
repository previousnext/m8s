#!/bin/bash

if [[ -n $BASIC_AUTH_USER || -n $BASIC_AUTH_PASS ]]; then
  htpasswd -b -c /etc/nginx/.htpasswd $BASIC_AUTH_USER $BASIC_AUTH_PASS
  sed -i -e "s/# BASIC_AUTH_PLACEHOLDER/auth_basic_user_file \/etc\/nginx\/.htpasswd;/g" /etc/nginx/nginx.conf
fi

echo "Setting Nginx port: ${HTTP_PORT}"
sed -i -e "s/HTTP_PORT/${HTTP_PORT}/g" /etc/nginx/nginx.conf

nginx -g "daemon off;"