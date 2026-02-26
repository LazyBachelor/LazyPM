FROM lscr.io/linuxserver/webtop:latest

RUN apk add --no-cache exo gcompat libc6-compat bash bash-completion alacritty

COPY pm /usr/local/bin/pm
COPY pm_bash.sh /etc/bash_completion.d/pm
COPY survey_bash.sh /etc/bash_completion.d/survey
RUN chmod +x /usr/local/bin/pm

COPY survey /usr/local/bin/survey
RUN chmod +x /usr/local/bin/survey

COPY setup-desktop.sh /custom-cont-init.d/01-setup-desktop.sh
RUN chmod +x /custom-cont-init.d/01-setup-desktop.sh