FROM ubuntu
RUN mkdir -p /k8s/log
WORKDIR  /k8s
ENV PATH /k8s/:$PATH
ADD k8s-load-simulator ./
CMD ["k8s-load-simulator","-h"]
