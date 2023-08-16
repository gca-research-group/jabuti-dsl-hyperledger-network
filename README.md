# Jabuti Hyperledger Fabirc Network

## References links

- [How Fabric networks are structured](https://hyperledger-fabric.readthedocs.io/en/release-2.5/network/network.html)
- [Hyperledger Fabric 1.4.0](https://github.com/deusimarferreira/hyperledger-fabric/tree/master)
- [Course of Hyperledger Fabric](https://github.com/blockchainempresarial/curso-hyperledger-fabric/tree/master)
- [Golang](https://hub.docker.com/_/golang)

## Configuring the environment
- Operational System supported
- [x] Ubuntu 22.04
- [x] Ubuntu 20.04

- Installing Go
```shell
sudo apt update

sudo apt install golang-go

go version
```

- Installing docker and docker-compose
```shell
sudo apt update

sudo apt install software-properties-common apt-transport-https ca-certificates lsb-release -y

curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo gpg --dearmor -o /usr/share/keyrings/docker-archive-keyring.gpg

echo "deb [arch=$(dpkg --print-architecture) signed-by=/usr/share/keyrings/docker-archive-keyring.gpg] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable" | sudo tee /etc/apt/sources.list.d/docker.list > /dev/null

sudo apt update

sudo apt install docker-ce -y

docker --version

docker-compose --version

```

- Downloading network files
```shell
cd ~
git clone https://github.com/mtelesborges/hyperledger-fabric.git
cd hyperledger-fabric
```

- Installing the binaries
```sh
# The binaries should be configures in global path
# The bast pratice the path sould be $HOME/go/src/github.com/<id of the repository on github>
mkdir -p $HOME/go/src/github.com/gca-research-group

cd $HOME/go/src/github.com/gca-research-group

# Downloading the files and add permitions
curl -sSLO https://raw.githubusercontent.com/hyperledger/fabric/main/scripts/install-fabric.sh && chmod +x install-fabric.sh

# Installing the binaries
bash ./install-fabric.sh binary

# Add binaries to the path
export PATH="$PATH:$HOME/go/src/github.com/gca-research-group/bin"
```

## Building the network
- Generating network digital certificates
```sh
cd ~/hyperledger-fabric/network

cryptogen generate --config=./crypto-config.yaml
```

- Generating genesis block
```sh
configtxgen -profile ThreeOrgsOrdererGenesis -channelID system-channel -outputBlock ./channel-artifacts/genesis.block
configtxgen -profile ThreeOrgsChannel -outputCreateChannelTx ./channel-artifacts/channel.tx -channelID marketplace
configtxgen -profile ThreeOrgsChannel -outputAnchorPeersUpdate ./channel-artifacts/Org1MSPanchors.tx -channelID marketplace -asOrg Org1MSP
configtxgen -profile ThreeOrgsChannel -outputAnchorPeersUpdate ./channel-artifacts/Org2MSPanchors.tx -channelID marketplace -asOrg Org2MSP
configtxgen -profile ThreeOrgsChannel -outputAnchorPeersUpdate ./channel-artifacts/Org3MSPanchors.tx -channelID marketplace -asOrg Org3MSP
```

## Configuring the network and publishing the chaincode

```shell
### Running the containers ###
cd ~/hyperledger-fabric/network

sudo docker-compose -f docker-compose-cli-couchdb.yaml up -d

# Configuring the channel
docker exec -it cli sh

export CHANNEL_NAME=marketplace

peer channel create -o orderer.gca.unijui.edu.br:7050 -c $CHANNEL_NAME -f ./channel-artifacts/channel.tx --tls true --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/gca.unijui.edu.br/orderers/orderer.gca.unijui.edu.br/msp/tlscacerts/tlsca.gca.unijui.edu.br-cert.pem

### Adding organizations to channel ###
docker exec -it cli sh

# Adding Org1
peer channel join -b marketplace.block

# Adding Org2
CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org2.gca.unijui.edu.br/users/Admin@org2.gca.unijui.edu.br/msp CORE_PEER_ADDRESS=peer0.org2.gca.unijui.edu.br:7051 CORE_PEER_LOCALMSPID="Org2MSP" CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org2.gca.unijui.edu.br/peers/peer0.org2.gca.unijui.edu.br/tls/ca.crt peer channel join -b marketplace.block

# Adding Org3
CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org3.gca.unijui.edu.br/users/Admin@org3.gca.unijui.edu.br/msp CORE_PEER_ADDRESS=peer0.org3.gca.unijui.edu.br:7051 CORE_PEER_LOCALMSPID="Org3MSP" CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org3.gca.unijui.edu.br/peers/peer0.org3.gca.unijui.edu.br/tls/ca.crt peer channel join -b marketplace.block

### Packing the chaincode ###
export CHAINCODE_NAME=foodcontrol
export CHAINCODE_VERSION=1
export CC_RUNTIME_LANGUAGE=golang
export CC_SRC_PATH=/opt/gopath/src/github.com/chaincode/foodcontrol/
export ORDERER_CA=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/gca.unijui.edu.br/orderers/orderer.gca.unijui.edu.br/msp/tlscacerts/tlsca.gca.unijui.edu.br-cert.pem

peer lifecycle chaincode package ${CHAINCODE_NAME}.tar.gz --path ${CC_SRC_PATH} --lang ${CC_RUNTIME_LANGUAGE} --label ${CHAINCODE_NAME}_${CHAINCODE_VERSION} >&log.txt

### Installing the chaincode ###

# Installing on the Org1
peer lifecycle chaincode install ${CHAINCODE_NAME}.tar.gz

export CC_PACKAGEID=<Id generated for the chaincode> 

# Installing on the Org2
CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org2.gca.unijui.edu.br/users/Admin@org2.gca.unijui.edu.br/msp CORE_PEER_ADDRESS=peer0.org2.gca.unijui.edu.br:7051 CORE_PEER_LOCALMSPID="Org2MSP" CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org2.gca.unijui.edu.br/peers/peer0.org2.gca.unijui.edu.br/tls/ca.crt peer lifecycle chaincode install ${CHAINCODE_NAME}.tar.gz

# Installing on the Org3
CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org3.gca.unijui.edu.br/users/Admin@org3.gca.unijui.edu.br/msp CORE_PEER_ADDRESS=peer0.org3.gca.unijui.edu.br:7051 CORE_PEER_LOCALMSPID="Org3MSP" CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org3.gca.unijui.edu.br/peers/peer0.org3.gca.unijui.edu.br/tls/ca.crt peer lifecycle chaincode install ${CHAINCODE_NAME}.tar.gz

### Creating chaincode approval policies in organizations ###

# Org 1
peer lifecycle chaincode approveformyorg --tls --cafile $ORDERER_CA --channelID $CHANNEL_NAME --name $CHAINCODE_NAME --version $CHAINCODE_VERSION --sequence 1 --waitForEvent --signature-policy "OR ('Org1MSP.peer','Org3MSP.peer')" --package-id foodcontrol_1:$CC_PACKAGEID

# Org 3
CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org3.gca.unijui.edu.br/users/Admin@org3.gca.unijui.edu.br/msp CORE_PEER_ADDRESS=peer0.org3.gca.unijui.edu.br:7051 CORE_PEER_LOCALMSPID="Org3MSP" CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org3.gca.unijui.edu.br/peers/peer0.org3.gca.unijui.edu.br/tls/ca.crt peer lifecycle chaincode approveformyorg --tls --cafile $ORDERER_CA --channelID $CHANNEL_NAME --name $CHAINCODE_NAME --version $CHAINCODE_VERSION --sequence 1 --waitForEvent --signature-policy "OR ('Org1MSP.peer','Org3MSP.peer')" --package-id foodcontrol_1:$CC_PACKAGEID

# Commiting policies
peer lifecycle chaincode commit -o orderer.gca.unijui.edu.br:7050 --tls --cafile $ORDERER_CA --peerAddresses peer0.org1.gca.unijui.edu.br:7051 --tlsRootCertFiles /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org1.gca.unijui.edu.br/peers/peer0.org1.gca.unijui.edu.br/tls/ca.crt --peerAddresses peer0.org3.gca.unijui.edu.br:7051 --tlsRootCertFiles /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org3.gca.unijui.edu.br/peers/peer0.org3.gca.unijui.edu.br/tls/ca.crt --channelID $CHANNEL_NAME --name $CHAINCODE_NAME --version $CHAINCODE_VERSION --sequence 1 --signature-policy "OR ('Org1MSP.peer','Org3MSP.peer')"

### Executing the chaincode ###
peer chaincode invoke -o orderer.gca.unijui.edu.br:7050 --tls --cafile $ORDERER_CA -C  $CHANNEL_NAME  -n $CHAINCODE_NAME -c '{"Args":["Set","did:3","josé","banana"]}'
```
