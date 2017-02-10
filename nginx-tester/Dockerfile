FROM bprashanth/nginx-test:1.0

COPY ./nginx.conf /usr/local/nginx/conf/nginx.conf
EXPOSE 80 443
ENTRYPOINT ["nginx", "-g", "daemon off;"]
