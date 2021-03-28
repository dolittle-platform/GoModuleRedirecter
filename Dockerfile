FROM golang:latest as build
WORKDIR /build
COPY . /build
ENV CGO_ENABLED=0
RUN go build .

FROM gcr.io/distroless/static:nonroot
COPY --from=build /build/redirecter /bin/redirecter
COPY .docker/group /etc/group
COPY .docker/passwd /etc/passwd
COPY --chown=redirecter:redirecter .docker/config.yaml /var/lib/redirecter/.redirecter.yaml
COPY --chown=redirecter:redirecter .docker/go-get.html /var/lib/redirecter/go-get.html
COPY --chown=redirecter:redirecter .docker/user.html /var/lib/redirecter/user.html
COPY --chown=redirecter:redirecter .docker/not-found.html /var/lib/redirecter/not-found.html
USER redirecter
ENTRYPOINT [ "/bin/redirecter" ]
CMD [ "serve" ]