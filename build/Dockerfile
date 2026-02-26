FROM lscr.io/linuxserver/webtop:arch-kde

COPY survey /usr/local/bin/survey
RUN chmod +x /usr/local/bin/survey

COPY setup-desktop.sh /custom-cont-init.d/01-setup-desktop.sh
RUN chmod +x /custom-cont-init.d/01-setup-desktop.sh