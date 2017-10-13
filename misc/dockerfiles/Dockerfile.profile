FROM library/ubuntu:16.04

ADD ./build/profile /opt/app/profile
ADD ./src/config/config.yaml /home/bv/.gvm/pkgsets/go1.8/global/src/config/config.yaml
EXPOSE 30000-65535
CMD ./opt/app/profile
