FROM ubuntu:latest
#解决在容器中报x509:certificate signed by unknown authority
RUN apt-get update && apt-get install -y ca-certificates
ENV WORKDIR /app
COPY ./build/main-linux $WORKDIR/main
RUN chmod +x $WORKDIR/main
WORKDIR $WORKDIR
CMD ./main