FROM golang:1.15.6-alpine3.12

WORKDIR /go/src/app

ENV GOlllMODULE=on

RUN apk add --update --no-cache \
    alpine-sdk \
    git \
    bash \
    bash-completion \
    sudo \
#デフォルトのシェルをbashに変更
&& sudo sed -e 's;/bin/ash$;/bin/bash;g' -i /etc/passwd \
&& echo -e "\n\
       PS1='\[\e[1;32m\]\u@\h:\[\e[0m\]\w\[\e[1;32m\]$ \[\e[0m\]' \n\
       " >> ~/.bashrc \
&& echo -e "\n\
       if [ '\${SHELL}x' = '/bin/bashx' ]; then \n\
         if [ -f '\${HOME}/.bashrc' ]; then \n\
           . '\${HOME}/.bashrc' \n\
         fi \n\
       fi \n\
       ">> ~/.profile

EXPOSE 80