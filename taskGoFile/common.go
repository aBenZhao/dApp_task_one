package taskGoFile

import (
	"crypto/ecdsa"
	"log"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

// 配置参数（按需替换）
const (
	privateKeyHex = "私钥"
	nodeURL       = "https://sepolia.infura.io/v3/8985d62b6a814e89ad9097ceeeef7d14"
)

func getFromAddressByPrivateKeyHex() (common.Address, *ecdsa.PrivateKey) {
	privateKeyECDSA, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		log.Fatal("私钥解析失败：%v", err)
	}
	publicKey := privateKeyECDSA.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("无法获取公钥")
	}

	//构造一笔简单的以太币转账交易，指定发送方、接收方和转账金额。
	// 发送方
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	return fromAddress, privateKeyECDSA
}
