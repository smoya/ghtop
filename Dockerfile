FROM centos:latest

COPY build/ghtop /opt/ghtop
ENTRYPOINT ["/opt/ghtop"]
