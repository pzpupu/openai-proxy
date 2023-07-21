FROM ubuntu
ENV WORKDIR /app
COPY ./main $WORKDIR/main
RUN chmod +x $WORKDIR/main
WORKDIR $WORKDIR
CMD ./main