FROM registry.access.redhat.com/ubi9/go-toolset:latest as builder
ENV GOPATH=$APP_ROOT
COPY --chown=1001:0 . .
RUN make cmd

FROM registry.access.redhat.com/ubi9/ubi-minimal
USER root
RUN echo -e "[centos8]" \
 "\nname = centos8" \
 "\nbaseurl = http://mirror.stream.centos.org/9-stream/AppStream/x86_64/os/" \
 "\nenabled = 1" \
 "\ngpgcheck = 0" > /etc/yum.repos.d/centos.repo

RUN microdnf -y install \
  java-17-openjdk-devel\
  openssh-clients \
  unzip \
  wget \
  git \
  subversion \
  maven \
&& microdnf -y clean all

ENV HOME=/working \
    JAVA_HOME="/usr/lib/jvm/java-17-openjdk/" \
    JAVA_VENDOR="openjdk" \
    JAVA_VERSION="17"

RUN /sbin/alternatives --set java java-17-openjdk.x86_64
RUN /sbin/alternatives --set javac java-17-openjdk.x86_64

WORKDIR /working
COPY --from=builder /opt/app-root/src/bin/addon /usr/local/bin/addon
ENTRYPOINT ["/usr/local/bin/addon"]
