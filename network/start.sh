WORKSPACE=`pwd`
FABRIC_HOME=$GOPATH"/src/github.com/hyperledger/fabric"
FABRIC_GO_SDK=$GOPATH"/src/github.com/hyperledger/fabric-sdk-go"
CHANNEL_ID="tracechain"
ORGANISATION_MSP="ManufacturerMSP"
FABRIC_VERSION="master"
clear(){
    rm -rf crypto-config artifacts
    cd $FABRIC_HOME
    echo "切换fabric版本为master"
    git checkout $FABRIC_VERSION
}

cryptoGenerate(){
    echo "进入"$WORKSPACE
    cd $WORKSPACE
    ../bin/cryptogen generate --config=crypto-config.yaml
}

configtxGenerate(){
    cd $WORKSPACE
    mkdir artifacts
    ../bin/configtxgen --profile ManufacturerOrdererGenesis -outputBlock ./artifacts/genesis.block 
    echo "生成创世块./artifacts/genesis.block"

    ../bin/configtxgen --profile ManufacturerChannel -outputCreateChannelTx ./artifacts/${CHANNEL_ID}.tx -channelID $CHANNEL_ID
    echo "创建通道配置交易./artifacts/${CHANNEL_ID}.tx"

    ../bin/configtxgen --profile ManufacturerChannel -outputAnchorPeersUpdate ./artifacts/${ORGANISATION_MSP}Anchors.tx -channelID $CHANNEL_ID -asOrg $ORGANISATION_MSP
    echo "创建瞄节点配置交易./artifacts/${ORGANISATION_MSP}Anchors.tx"
}

startNetwork(){
    docker-compose up -d
}

downNetwork(){
    echo "downNetwork() nothing to do"
}

ARG1=$1
if ["${ARG1}"=="start"];then
    startNetwork
elif ["${ARG1}"=="generate"];then
    cryptoGenerate
    configtxGenerate
elif ["{$ARG1}"=="down"];then
    downNetwork
else
    downNetwork
    clear
    cryptoGenerate
    configtxGenerate
    startNetwork
fi


