FROM ubuntu:22.04

USER root

RUN apt-get update \
    && apt-get -y install wget libasound2-dev git build-essential gcc pkg-config libwayland-dev libx11-dev libx11-xcb-dev libxkbcommon-x11-dev libgles2-mesa-dev libegl1-mesa-dev libffi-dev libxcursor-dev libvulkan-dev

RUN wget https://go.dev/dl/go1.23.6.linux-amd64.tar.gz  \
    && rm -rf /usr/local/go && tar -C /usr/local -xzf go1.23.6.linux-amd64.tar.gz && ln -sf /usr/local/go/bin/go /usr/bin/go

#ENV PATH=$PATH:/usr/local/go/bin
#RUN go install golang.org/dl/go1.23.6@latest
#RUN go1.23.6 download && ln -sf /root/go/bin/go1.23.6 /usr/bin/go

RUN git config --global --add safe.directory /app
