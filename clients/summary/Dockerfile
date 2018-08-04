# add the necessary instructions
# to create a Docker container image
# for your summary client using the
# NGINX web server image as the base
FROM nginx
ADD default.conf /etc/nginx/conf.d/default.conf
COPY build/index.html /usr/share/nginx/html
COPY build/static/ /usr/share/nginx/html/static
