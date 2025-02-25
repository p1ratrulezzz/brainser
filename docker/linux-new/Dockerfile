FROM phusion/baseimage:noble-1.0.0

USER root

RUN apt-get update \
    && apt-get -y install libasound2-dev git golang build-essential gcc pkg-config libwayland-dev libx11-dev libx11-xcb-dev libxkbcommon-x11-dev libgles2-mesa-dev libegl1-mesa-dev libffi-dev libxcursor-dev libvulkan-dev

RUN git config --global --add safe.directory /app

ENV PATH="/root/go/bin:${PATH}"
RUN go install golang.org/dl/go1.23.6@latest && go1.23.6 download && ln -sf /root/go/bin/go1.23.6 /usr/bin/go

