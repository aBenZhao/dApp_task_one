package taskGoFile

import (
	"context"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

//任务 1：区块链读写 任务目标
//使用 Sepolia 测试网络实现基础的区块链交互，包括查询区块和发送交易。
//具体任务
//环境搭建
//安装必要的开发工具，如 Go 语言环境、 go-ethereum 库。
//注册 Infura 账户，获取 Sepolia 测试网络的 API Key。
//查询区块
//编写 Go 代码，使用 ethclient 连接到 Sepolia 测试网络。
//实现查询指定区块号的区块信息，包括区块的哈希、时间戳、交易数量等。
//输出查询结果到控制台。
//发送交易
//准备一个 Sepolia 测试网络的以太坊账户，并获取其私钥。
//编写 Go 代码，使用 ethclient 连接到 Sepolia 测试网络。
//构造一笔简单的以太币转账交易，指定发送方、接收方和转账金额。
//对交易进行签名，并将签名后的交易发送到网络。
//输出交易的哈希值。

func BlockchainReadWrite() {
	client, err := ethclient.Dial(nodeURL)
	if err != nil {
		log.Fatalf("节点连接失败：%v", err)
	}
	defer client.Close()

	//查询区块
	//实现查询指定区块号的区块信息，包括区块的哈希、时间戳、交易数量等。
	//输出查询结果到控制台。
	number, err := client.BlockByNumber(context.Background(), big.NewInt(9591371))
	if err != nil {
		log.Fatalf("查询区块失败：%v", err)
	}
	log.Printf("区块哈希：%s", number.Hash().Hex())
	log.Printf("时间戳：%d", number.Time())
	log.Printf("交易数量：%d", len(number.Transactions()))

	//2025/11/09 15:19:36 区块哈希：0x70dda5b226657a35ef28a07a0a569d4001bc20c6401916a10ad115298d39aea5
	//2025/11/09 15:19:36 时间戳：1762672224
	//2025/11/09 15:19:36 交易数量：109

	//发送交易
	//准备一个 Sepolia 测试网络的以太坊账户，并获取其私钥。
	//构造一笔简单的以太币转账交易，指定发送方、接收方和转账金额。
	// 发送方
	fromAddress, privateKeyECDSA := getFromAddressByPrivateKeyHex()

	// 接收方
	toAddress := common.HexToAddress("0x0d6913C10f0F2E63b56fb6FD56E24D510c4AF538")

	// 获取交易序号
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatal(err)
	}

	// 转账金额
	amount := big.NewInt(100000)
	gasLimit := uint64(21000)
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	// 构建未签名交易
	var data []byte
	tx := types.NewTransaction(nonce, toAddress, amount, gasLimit, gasPrice, data)

	//对交易进行签名，并将签名后的交易发送到网络。
	chainId, err := client.NetworkID(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainId), privateKeyECDSA)
	if err != nil {
		log.Fatal(err)
	}

	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		log.Fatal(err)
	}
	//输出交易的哈希值。
	log.Printf("转账交易的哈希值：%s\n", signedTx.Hash().Hex())
	// 2025/11/09 16:02:34 转账交易的哈希值：0x6da7358e9609bca422f0c82d51724e75d7acee1f0da43509c3858d461f4b0c27)
}
