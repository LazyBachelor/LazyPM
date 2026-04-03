FROM lscr.io/linuxserver/webtop:latest

ENV SELKIES_MANUAL_WIDTH=1920
ENV SELKIES_MANUAL_HEIGHT=1080
ENV MAX_RESOLUTION=1920x1080
ENV PIXELFLUX_WAYLAND=true
ENV TITLE="LazyOS - Project Management Survey"

RUN apk add --no-cache \
	ca-certificates \
	exo \
	gcompat \
	libc6-compat \
	bash \
	bash-completion && \
	update-ca-certificates

COPY pm /usr/local/bin/pm
COPY pm_bash.sh /etc/bash_completion.d/pm
RUN chmod +x /usr/local/bin/pm

COPY setup-desktop.sh /custom-cont-init.d/01-setup-desktop.sh
RUN chmod +x /custom-cont-init.d/01-setup-desktop.sh

COPY PARTICIPANT_README.md /config/Desktop/README.md