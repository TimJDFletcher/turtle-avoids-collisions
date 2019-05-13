FROM armhf/ubuntu

WORKDIR /

COPY cavoidance ./cavoidance

RUN chmod +x ./cavoidance

CMD ["./cavoidance"]
