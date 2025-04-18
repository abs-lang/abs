FROM ghcr.io/charmbracelet/vhs

RUN apt-get update
ENV TZ="Asia/Dubai"
RUN apt-get install -y curl vim tzdata
RUN ln -f -s /usr/share/zoneinfo/Asia/Dubai /etc/localtime
RUN dpkg-reconfigure --frontend noninteractive tzdata

ENTRYPOINT ["abs"]
CMD ["/abs/docs/vhs/vhs.abs"]