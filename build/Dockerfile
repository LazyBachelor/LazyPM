FROM lscr.io/linuxserver/webtop:latest

RUN apk add --no-cache exo gcompat libc6-compat bash bash-completion alacritty

COPY pm /usr/local/bin/pm
COPY pm_bash.sh /etc/bash_completion.d/pm
RUN chmod +x /usr/local/bin/pm

COPY setup-desktop.sh /custom-cont-init.d/01-setup-desktop.sh
RUN chmod +x /custom-cont-init.d/01-setup-desktop.sh