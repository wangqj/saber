FROM golang:1.9

ENV GOPATH /gopath
ENV SABER  ${GOPATH}/src/saber
ENV PATH   ${GOPATH}/bin:${PATH}:${SABER}/bin
COPY . ${SABER}
#ADD . .
#RUN make -C ${SABER} distclean
#RUN make -C ${SABER} build-all

WORKDIR /saber
#CMD ["ls -l"]
#CMD ["saber-proxy -c=config.toml"]