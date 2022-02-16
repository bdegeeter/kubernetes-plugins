FROM golang
ENV PORTER_VERSION v1.0.0-alpha.9
RUN curl --silent -L https://cdn.porter.sh/${PORTER_VERSION}/install-linux.sh | bash >/dev/null 2>&1 && \
    curl --silent -LO "https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/linux/amd64/kubectl" && \
    install -o root -g root -m 0755 kubectl /usr/local/bin/kubectl
ENV PATH "$PATH:/root/.porter"
RUN mkdir -p /root/.porter/plugins/kubernetes  &&\
    mkdir -p /root/.porter/credentials
COPY / /app
RUN install /app/bin/plugins/kubernetes/kubernetes /root/.porter/plugins/kubernetes/kubernetes
WORKDIR /app
COPY ./tests/integration/scripts/config-*.toml /root/.porter/
#TODO: update to use porter credentials apply 
COPY ./tests/testdata/kubernetes-plugin-test-*.json /root/.porter/credentials/
ENTRYPOINT ["tail","-f","/dev/null"]