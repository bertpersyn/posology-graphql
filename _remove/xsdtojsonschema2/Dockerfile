FROM swipl

WORKDIR /data

RUN apt update
RUN apt-get install curl git -y
RUN curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.35.3/install.sh | bash

RUN bash -c 'source $HOME/.nvm/nvm.sh   && \
    nvm install node                    && \
    npm install -g doctoc urchin eclint dockerfile_lint && \
    npm install --prefix "$HOME/.nvm/"'
    
RUN git clone https://github.com/MGDIS/xsd2json /tmp/xsd2json

RUN bash -c 'source $HOME/.nvm/nvm.sh && \
            npm install -g xsd2json2'

RUN where xsd2json2
ENTRYPOINT [ "/usr/bin/xsd2json2" ]