
FROM scratch

#ENV GOPROXY https://goproxy.cn,direct
WORKDIR $GOPATH/src/gogin
COPY . $GOPATH/src/gogin
RUN #go build .

EXPOSE 8080
#ENTRYPOINT ["./gogin"]
CMD ["./gogin"]