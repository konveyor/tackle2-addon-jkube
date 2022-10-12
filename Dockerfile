FROM registry.access.redhat.com/ubi8/go-toolset:1.16.12 as builder
ENV GOPATH=$APP_ROOT
COPY --chown=1001:0 . .
RUN make cmd

FROM registry.access.redhat.com/ubi8/ubi-minimal
USER root
RUN echo -e "[centos8]" \
 "\nname = centos8" \
 "\nbaseurl = http://mirror.centos.org/centos/8-stream/AppStream/x86_64/os/" \
 "\nenabled = 1" \
 "\ngpgcheck = 0" > /etc/yum.repos.d/centos.repo
RUN echo -e "[WandiscoSVN]" \
 "\nname=Wandisco SVN Repo" \
 "\nbaseurl=http://opensource.wandisco.com/centos/6/svn-1.9/RPMS/$basearch/" \
 "\nenabled=1" \
 "\ngpgcheck=0" > /etc/yum.repos.d/wandisco.repo

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
    JAVA_HOME="/usr/lib/jvm/jre-17" \
    JAVA_VENDOR="openjdk" \
    JAVA_VERSION="17"

RUN /sbin/alternatives --set java /usr/lib/jvm/java-17-openjdk-17.0.4.1.1-2.el8_6.x86_64/bin/java
RUN /sbin/alternatives --set javac /usr/lib/jvm/java-17-openjdk-17.0.4.1.1-2.el8_6.x86_64/bin/javac


WORKDIR /working
COPY --from=builder /opt/app-root/src/bin/addon /usr/local/bin/addon
ENTRYPOINT ["/usr/local/bin/addon"]