FROM amazonlinux:2

LABEL io.buildpacks.stack.id devflow.lambda.ruby

ENV CNB_USER_ID 99
ENV CNB_GROUP_ID 99

ADD ruby2.5.tgz /

RUN yum install -y wget tar gzip openssl gcc binutils make glibc-devel gmp-devel procps-ng iproute && yum clean all && \
  rm -rf /var/cache/yum /var/lib/rpm/__db.* && \
  > /var/log/yum.log && \
  mkdir -p /root /tmp /var/task

ENV PATH=/var/lang/bin:/usr/local/bin:/usr/bin/:/bin:/opt/bin \
    LD_LIBRARY_PATH=/var/lang/lib:/lib64:/usr/lib64:/var/runtime:/var/runtime/lib:/var/task:/var/task/lib:/opt/lib \
    LANG=en_US.UTF-8 \
    TZ=:UTC

ADD buildpack /buildpack

ADD bin/builder /builder
ADD bin/lambda-runner /lambda-runner

ADD bootstrap-raw /var/runtime/bootstrap-raw

ADD prompt.sh /etc/profile.d/prompt.sh

RUN rm -rf /var/cache/yum /var/lib/rpm/__db.* && \
  > /var/log/yum.log && \
  mkdir -p /root /tmp /var/task

WORKDIR /var/task

CMD ["/builder"]
