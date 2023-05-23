#!/bin/bash
set -e

GO_VER="1.19.3"
GOROOT="/usr/local/go"
GOPATH="/home/go"
GOCACHE="/home/go/cache"
GOBIN="${GOROOT}/bin"

SEKAI_BRANCH="v0.3.4.28"
INTERX_BRANCH="v0.4.21"
TOOLS_VERSION="v0.2.20" 
COSIGN_VERSION="v1.7.2"

DEFAULT_INTERX_PORT=11000

apt update && apt upgrade -y || "Upgraded"

mkdir ~/tmp && cd ~/tmp || echo "Entered tmp dir"

wget "https://go.dev/dl/go${GO_VER}.linux-amd64.tar.gz"
rm -rf /usr/local/go && tar -C /usr/local -xzf "go${GO_VER}.linux-amd64.tar.gz" || echo "Go installed" 

curl -fsSL https://get.docker.com -o get-docker.sh
sh ./get-docker.sh

echo "export PATH=$PATH:/usr/local/go/bin" >> ~/.profile
source ~/.profile

version=$(go version)

apt update && apt upgrade -y || "Upgraded"

apt-get install build-essential -y
apt-get install jq -y

# install cosign
if [[ "$(uname -m)" == *"ar"* ]] ; then ARCH="arm64"; else ARCH="amd64" ; fi && echo $ARCH && \
PLATFORM=$(uname) && FILE=$(echo "cosign-${PLATFORM}-${ARCH}" | tr '[:upper:]' '[:lower:]') && \
 wget https://github.com/sigstore/cosign/releases/download/${COSIGN_VERSION}/$FILE && chmod +x -v ./$FILE && \
 mv -fv ./$FILE /usr/local/bin/cosign && cosign version

# save KIRA public cosign key
KEYS_DIR="/usr/keys" && KIRA_COSIGN_PUB="${KEYS_DIR}/kira-cosign.pub" && \
mkdir -p $KEYS_DIR  && cat > ./cosign.pub << EOL
-----BEGIN PUBLIC KEY-----
MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAE/IrzBQYeMwvKa44/DF/HB7XDpnE+
f+mU9F/Qbfq25bBWV2+NlYMJv3KvKHNtu3Jknt6yizZjUV4b8WGfKBzFYw==
-----END PUBLIC KEY-----
EOL

# download desired files and the corresponding .sig file from: https://github.com/KiraCore/tools/releases

# verify signature of downloaded files
# cosign verify-blob --key=$KIRA_COSIGN_PUB--signature=./<file>.sig ./<file>

mkdir -p /usr/keys && FILE_NAME="bash-utils.sh" && \
 if [ -z "$KIRA_COSIGN_PUB" ] ; then KIRA_COSIGN_PUB=/usr/keys/kira-cosign.pub ; fi && \
 echo -e "-----BEGIN PUBLIC KEY-----\nMFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAE/IrzBQYeMwvKa44/DF/HB7XDpnE+\nf+mU9F/Qbfq25bBWV2+NlYMJv3KvKHNtu3Jknt6yizZjUV4b8WGfKBzFYw==\n-----END PUBLIC KEY-----" > $KIRA_COSIGN_PUB && \
 wget "https://github.com/KiraCore/tools/releases/download/$TOOLS_VERSION/${FILE_NAME}" -O ./$FILE_NAME && \
 wget "https://github.com/KiraCore/tools/releases/download/$TOOLS_VERSION/${FILE_NAME}.sig" -O ./${FILE_NAME}.sig && \
 cosign verify-blob --key="$KIRA_COSIGN_PUB" --signature=./${FILE_NAME}.sig ./$FILE_NAME && \
 chmod -v 555 ./$FILE_NAME && ./$FILE_NAME bashUtilsSetup "/var/kiraglob" && . /etc/profile && \
 echoInfo "Installed bash-utils $(bashUtilsVersion)"

set -x

BIN_DEST="/usr/local/bin/sekai-utils.sh" && \
wget "https://github.com/KiraCore/sekai/releases/download/$SEKAI_BRANCH/sekai-utils.sh" -O ./sekai-utils.sh \
&& chmod -v 755 ./sekai-utils.sh && ./sekai-utils.sh sekaiUtilsSetup && chmod -v 755 $BIN_DEST && . /etc/profile

FILE=/usr/local/bin/sekai-env.sh && \
wget "https://github.com/KiraCore/sekai/releases/download/$SEKAI_BRANCH/sekai-env.sh" -O $FILE \
&& chmod -v 755 $FILE && echo "source $FILE" >> /etc/profile && . /etc/profile
    
cd $HOME && rm -fvr ./sekai && \
git clone https://github.com/KiraCore/sekai.git -b "release/$SEKAI_BRANCH" && \
cd ./sekai && chmod -R 777 ./scripts && \
make install && echo "SUCCESS installed sekaid $(sekaid version)" || echo "FAILED" 
    
    
cd $HOME && rm -fvr ./interx && \
git clone https://github.com/KiraCore/interx.git -b "release/$INTERX_BRANCH" && \
cd ./interx && chmod -R 777 ./scripts && \

set -e
make install && echo "SUCCESS" || echo "FAILED" 
make test && echo "SUCCESS" || echo "FAILED" 
make test-local && echo "SUCCESS" || echo "FAILED" 

. /etc/profile

TEST_NAME="NETWORK-START" && timerStart $TEST_NAME
echoInfo "INFO: $TEST_NAME - Integration Test - START"

echoInfo "INFO: Ensuring essential dependencies are installed & up to date"
SYSCTRL_DESTINATION=/usr/local/bin/systemctl

TEST_NAME="NETWORK-START" && timerStart $TEST_NAME

ARCHITECURE=$(getArch)
PLATFORM="$(getPlatform)"
DEFAULT_GRPC_PORT=9090
DEFAULT_RPC_PORT=26657
DEFAULT_INTERX_PORT=11000
PING_TARGET="127.0.0.1"
CFG_grpc="dns:///$PING_TARGET:$DEFAULT_GRPC_PORT"
CFG_rpc="http://$PING_TARGET:$DEFAULT_RPC_PORT"


echoInfo "INFO: Environment cleanup...."
NETWORK_NAME="localnet-1"
setGlobEnv SEKAID_HOME ~/.sekaid-$NETWORK_NAME
setGlobEnv INTERXD_HOME ~/.interxd-$NETWORK_NAME
setGlobEnv NETWORK_NAME $NETWORK_NAME
loadGlobEnvs

rm -rfv "$SEKAID_HOME" "$INTERXD_HOME"
mkdir -p "$SEKAID_HOME" "$INTERXD_HOME/cache"

cp /sekaid $GOBIN && echo "SUCCESS sekaid copied" || echo "FAILED sekaid copied" 
cp /interxd $GOBIN && echo "SUCCESS interxd copied" || echo "FAILED interxd copied" 

echoInfo "INFO: Starting new network..."
$GOBIN/sekaid init --overwrite --chain-id=$NETWORK_NAME "KIRA TEST LOCAL VALIDATOR NODE" --home=$SEKAID_HOME
addAccount validator


echo $(addAccount validator | jq .mnemonic | xargs) > $SEKAID_HOME/sekai.mnemonic
echo $(addAccount interx | jq .mnemonic | xargs) > $INTERXD_HOME/interx.mnemonic
echo $(addAccount faucet | jq .mnemonic | xargs) > $INTERXD_HOME/faucet.mnemonic
$GOBIN/sekaid add-genesis-account $(showAddress validator) 150000000000000ukex,300000000000000test,2000000000000000000000000000samolean,1000000lol --keyring-backend=test --home=$SEKAID_HOME
$GOBIN/sekaid add-genesis-account $(showAddress faucet) 150000000000000ukex,300000000000000test,2000000000000000000000000000samolean,1000000lol --keyring-backend=test --home=$SEKAID_HOME
$GOBIN/sekaid gentx-claim validator --keyring-backend=test --moniker="GENESIS VALIDATOR" --home=$SEKAID_HOME

cat > /etc/systemd/system/sekai.service << EOL
[Unit]
Description=Local KIRA Test Network
After=network.target
[Service]
MemorySwapMax=0
Type=simple
User=root
WorkingDirectory=/root
ExecStart=$GOBIN/sekaid start --home=$SEKAID_HOME --trace
Restart=always
RestartSec=5
LimitNOFILE=4096
[Install]
WantedBy=default.target
EOL

systemctl enable sekai 
systemctl start sekai

echoInfo "INFO: Waiting for network to start..." && sleep 3

echoInfo "INFO: Checking network status..."
NETWORK_STATUS_CHAIN_ID=$(showStatus | jq .NodeInfo.network | xargs)

if [ "$NETWORK_NAME" != "$NETWORK_STATUS_CHAIN_ID" ] ; then
    echoErr "ERROR: Incorrect chain ID from the status query, expected '$NETWORK_NAME', but got $NETWORK_STATUS_CHAIN_ID"
fi

echoInfo "INFO: Initalizing interxd..."

$GOBIN/interxd init --cache_dir="$INTERXD_HOME/cache" --home="$INTERXD_HOME" --grpc="$CFG_grpc" --rpc="$CFG_rpc" --port="$INTERNAL_API_PORT" \
    --signing_mnemonic="$INTERXD_HOME/interx.mnemonic" \
    --faucet_mnemonic="$INTERXD_HOME/faucet.mnemonic" \
    --port="$DEFAULT_INTERX_PORT" \
    --node_type="validator" \
    --seed_node_id="" \
    --sentry_node_id="" \
    --validator_node_id="$(globGet validator_node_id)" \
    --addrbook="$(globFile KIRA_ADDRBOOK)" \
    --faucet_time_limit=30 \
    --faucet_amounts="100000ukex,20000000test,300000000000000000samolean,1lol" \
    --faucet_minimum_amounts="1000ukex,50000test,250000000000000samolean,1lol" \
    --fee_amounts="ukex 1000ukex,test 500ukex,samolean 250ukex,lol 100ukex"

    
cat > /etc/systemd/system/interx.service << EOL
[Unit]
Description=Local Interx
After=network.target
[Service]
MemorySwapMax=0
Type=simple
User=root
WorkingDirectory=/root
ExecStart=$GOBIN/interxd start --home="$INTERXD_HOME"
Restart=always
RestartSec=5
LimitNOFILE=4096
[Install]
WantedBy=default.target
EOL

systemctl enable interx.service 
systemctl start interx.service

echoInfo "INFO: Waiting for interx to start..." && sleep 3

INTERX_GATEWAY="127.0.0.1:$DEFAULT_INTERX_PORT"

echoInfo "INFO: Waiting for next block to be produced..."
BLOCK_HEIGHT=$(curl --fail $INTERX_GATEWAY/api/status | jsonParse "interx_info.latest_block_height" || echo "0")
timeout 60 sekai-utils awaitBlocks 2
NEXT_BLOCK_HEIGHT=$(curl --fail $INTERX_GATEWAY/api/status | jsonParse "interx_info.latest_block_height" || echo "0")

if [ $BLOCK_HEIGHT -ge $NEXT_BLOCK_HEIGHT ] ; then
    echoErr "ERROR: INTERX failed to catch up with the latest sekai block height, stuck at $BLOCK_HEIGHT"
fi

echoInfo "INFO: Printing interx status..."
curl --fail $INTERX_GATEWAY/api/status | jq

sudo iptables -t nat -A OUTPUT -j DNAT --to-destination 127.0.0.1
sudo iptables -t nat -A POSTROUTING -j MASQUERADE


set +x
echoInfo "INFO: SEKAID $(sekaid version) is running"
echoInfo "INFO: INTERXD $(interxd version) is running"
echoInfo "INFO: NETWORK-START - Integration Test - END, elapsed: $(prettyTime $(timerSpan $TEST_NAME))"




HASH="v0.11.16" && \
 cd /tmp && wget https://ipfs.kira.network/ipfs/$HASH/init.sh -O ./i.sh && \
 chmod +x -v ./i.sh && ./i.sh --infra-src="$HASH" --init-mode="interactive"
