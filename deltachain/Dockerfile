FROM docker.io/golang:1.19.9@sha256:86901bd215ca2b33bdef68f70b999140252169bb55883bb29347334fa00b5b62

#used for a readiness probe
RUN apt-get update
RUN apt install -y netcat
RUN apt install -y jq

WORKDIR /app

COPY ./deltachain/go.mod .
COPY ./deltachain/go.sum .
COPY ./sdk /sdk
RUN go mod download

# copy over c bindings (libwasmvm.x86_64.so, etc)
RUN cp -r /go/pkg/mod/github.com/!cosm!wasm/wasmvm@v1.1.1/internal/api/* /usr/lib

COPY ./deltachain .

EXPOSE 26657
EXPOSE 26656
EXPOSE 6060
EXPOSE 9090
EXPOSE 1317
EXPOSE 4500

RUN unset GOPATH

# create the dir for the default "home" config
RUN mkdir -p /root/.deltachain

ARG num_phylaxs
ENV NUM_PHYLAXS=$num_phylaxs

RUN /bin/bash /app/devnet/create-genesis.sh

RUN make client
RUN chmod +x /app/build/deltachaind


ENTRYPOINT ["/bin/bash","-c","/app/build/deltachaind start"]
