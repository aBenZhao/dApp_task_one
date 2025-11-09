package taskGoFile

import (
	"context"
	"dApp_task_one/counter"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

//任务 2：合约代码生成 任务目标
//使用 abigen 工具自动生成 Go 绑定代码，用于与 Sepolia 测试网络上的智能合约进行交互。
// 具体任务
//编写智能合约
//使用 Solidity 编写一个简单的智能合约，例如一个计数器合约。
//编译智能合约，生成 ABI 和字节码文件。
//使用 abigen 生成 Go 绑定代码
//安装 abigen 工具。
//使用 abigen 工具根据 ABI 和字节码文件生成 Go 绑定代码。
//使用生成的 Go 绑定代码与合约交互
//编写 Go 代码，使用生成的 Go 绑定代码连接到 Sepolia 测试网络上的智能合约。
//调用合约的方法，例如增加计数器的值。
//输出调用结果。

func ContractCodeGeneration() {
	client, err := ethclient.Dial(nodeURL)
	if err != nil {
		log.Fatalf("节点连接失败：%v", err)
	}

	fromAddress, privateKeyECDSA := getFromAddressByPrivateKeyHex()

	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatal(err)
	}

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	chainId, err := client.NetworkID(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	auth, err := bind.NewKeyedTransactorWithChainID(privateKeyECDSA, chainId)
	if err != nil {
		log.Fatal(err)
	}
	auth.Nonce = big.NewInt(int64(nonce))
	auth.GasPrice = gasPrice
	auth.Value = big.NewInt(0)
	auth.GasLimit = uint64(300000)

	contractAddr, txHash, instance, err := counter.DeployCounter(auth, client)
	if err != nil {
		log.Fatal(err)
	}

	// 打印部署结果
	log.Printf("合约部署交易哈希：%s", txHash.Hash().Hex())
	log.Printf("合约地址：%s", contractAddr.Hex()) // 后续可通过该地址交互

	// 持续等待交易被打包进区块，上链后返回收据（主动等待，阻塞直到有结果）
	receipt, err := bind.WaitMined(context.Background(), client, txHash)
	if err != nil {
		log.Fatalf("获取交易收据失败（可能交易超时）：%v", err)
	}

	// 验证部署是否成功
	if receipt.Status == types.ReceiptStatusSuccessful {
		log.Println("✅ 合约部署成功！")
		log.Printf("合约地址（收据中获取，最可靠）：%s", receipt.ContractAddress.Hex())
		log.Printf("部署区块号：%d", receipt.BlockNumber.Uint64())
		log.Printf("实际消耗 Gas：%d", receipt.GasUsed)
	} else {
		log.Fatalf("❌ 合约部署失败！交易状态：%d", receipt.Status)
	}

	nonce, err = client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatal(err)
	}
	auth.Nonce = big.NewInt(int64(nonce))

	txIn, err := instance.Increment(auth)
	if err != nil {
		log.Fatal(err)
	}

	receiptWait, errWait := bind.WaitMined(context.Background(), client, txIn)
	if errWait != nil {
		log.Fatalf("获取交易收据失败（可能交易超时）：%v", err)
	}

	if receiptWait.Status == types.ReceiptStatusSuccessful {
		log.Println("✅ Increment成功！")
	} else {
		log.Fatalf("❌ Increment成功失败！交易状态：%d", receipt.Status)
	}

	count, err := instance.Count(&bind.CallOpts{})
	if err != nil {
		log.Fatalf("查询计数失败：%v", err)
	}
	log.Printf("当前计数：%d", count.Uint64())

}

//2025/11/09 19:52:53 合约部署交易哈希：0x32118ca1cd5d3db67677007b1917d2e304a1a3607c57bb4d60f189e1e28cbcaf
//2025/11/09 19:52:53 合约地址：0x44c146AE3b9D9ea23798b9F6c7371bcCa2588972
//2025/11/09 19:53:02 ✅ 合约部署成功！
//2025/11/09 19:53:02 合约地址（收据中获取，最可靠）：0x44c146AE3b9D9ea23798b9F6c7371bcCa2588972
//2025/11/09 19:53:02 部署区块号：9592682
//2025/11/09 19:53:02 实际消耗 Gas：120951
//2025/11/09 19:53:13 ✅ Increment成功！
//2025/11/09 19:53:14 当前计数：1
