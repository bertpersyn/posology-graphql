FROM swipl

WORKDIR /data

RUN apt update
RUN apt-get install curl git -y
RUN curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.35.3/install.sh | bash

RUN bash -c 'source $HOME/.nvm/nvm.sh   && \
    nvm install node                    && \
    npm install -g doctoc urchin eclint dockerfile_lint && \
    npm install --prefix "$HOME/.nvm/"'

RUN git clone https://github.com/fnogatz/xsd2json.git /tmp/xsd2json

RUN bash -c 'source $HOME/.nvm/nvm.sh && \
            cd /tmp/xsd2json && npm run preinstall && \
            cp /tmp/xsd2json/lib-pl/cli.exe /usr/bin'

ENTRYPOINT [ "/usr/bin/cli.exe" ]