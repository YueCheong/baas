# 移动区块链后端接口使用实例

## 网络

### 创建网络

> 注，由于加入区块链系统底层的网络分隔管理功能，需要对系统中的命名进行修改。从此版本起，网络指代区块链系统所属的网络。每个区块链系统至少属于一个网络，一个网络中可以有多个区块链系统。网络通过网络名确定，名字不可重复。

**url**

http://192.168.203.2:8080/api/networks/Testnet

**method**

PUT

**body**

```
无
```

### 查看网络

**url**

http://192.168.203.2:8080/api/networks

**method**

GET

### 删除网络 

**url**

http://192.168.203.2:8080/api/networks/Testnet

**method**

DELETE

## 区块链系统
### 区块链系统创建流程：

1. 创建区块链系统配置

   > 区块链创建后状态为Configuring

2. 配置区块链系统

   > 添加组织和节点

3. 初始化区块链系统

   > 区块链状态转为Stop

4. 启动区块链系统

   > 区块链状态转为Running

### 创建区块链系统配置

> 从新版本开始，此API创建后的区块链将处于配置阶段。可以继续添加节点和组织。需要对区块链进行初始化操作，区块链才会被创建和启动。初始化后无法继续向区块链中添加组织和节点



**url**

http://192.168.203.2:8080/api/blockchains

**method**

PUT

**body**

> 在新版本中，增加了指定创建的区块链系统所属网络的字段。在创建区块链系统时，如果将该字段留空，则会自动为该系统创建其ID对应的独立网络。如果要指定创建的区块链系统所属的网络，则需要在该字段中填入已存在的网络名。

##### 创建区块链系统后再进行具体配置(自动创建独立的网络)

``` 
{
    "Name": "secondnet"
}
```

> 现在创建区块链系统时，可以不提供完整的网络配置。可以不提供orderer org和peer org，也可提供org但并不添加节点。



##### 创建区块链时进行配置（自动创建独立的网络）

```
{
    "ID": 1,
    "Name": "firstnet",
    "OrdererOrg": {
        "Name": "Orderer",
        "Domain": "example.com",
        "Orderer": {
            "Name": "Orderer",
            "Port": "5050"
        },
        "MSPID": "OrdererMSP"
    },
    "PeerOrg": [
        {
            "Name": "org1",
            "Domain": "example.com",
            "MSPID": "Org1MSP",
            "Peers": [
                {
                    "Name": "peer0",
                    "Port": "5051"
                },
                {
                    "Name": "peer1",
                    "Port": "6051"
                }
            ]
        }
    ]
}
```



##### 让系统加入已有的网络

``` 
{
    "Name": "secondnet",
    "Netname":"Testnet"
}
```





### 向区块链添加组织或节点

**url**

http://192.168.203.2:8080/api/blockchains

**method**

POST

**body**

##### 添加orderer组织

```json
{
    "BlockchainID":1,
    "Operation":4,
    "Name":"Orderer",
    "Domain":"example.com",
    "MSPID":"OrdererMSP"
}
```

##### 添加orderer 节点

```json
{
    "BlockchainID":1,
    "Operation":5,
    "Name":"Orderer",
    "Domain":"example.com",
    "MSPID":"OrdererMSP",
    "Nodes":[
        {
            "Name": "Orderer",
            "Port": "5050"
        }
    ]
}
```

##### 添加peer组织

```json
{
    "BlockchainID":1,
    "Operation":6,
    "Name":"org1",
    "Domain": "example.com",
    "MSPID":"Org1MSP"
}
```

##### 添加peer节点

```json
{
    "BlockchainID":1,
    "Operation":7,
    "Name":"org1",
    "Domain": "example.com",
    "MSPID":"Org1MSP",
    "Nodes":[
        {
            "Name": "peer0",
            "Port": "5051"
        },
        {
            "Name": "peer1",
            "Port": "6051"
        }
    ]
}
```



### 

### 初始化区块链

**url**

http://192.168.203.2:8080/api/blockchains

**method**

POST

**body**

```json
{
    "BlockchainID":1,
    "Operation":3
}
```



### 启动和停止区块链系统

**url**

http://192.168.203.2:8080/api/blockchains

**method**

POST

**body**

```
{
    "BlockchainID":1,
    "Operation":1
}
```

参数Operation的值 1代表启动网络，2表示停止网络。区块链系统在创建后会自动运行。可以停止正在运行的区块链系统，也可以启动停止的区块链系统。

### 删除区块链系统

**url**

http://192.168.203.2:8080/api/blockchains/1

**method**

DELETE

**body**

空



### 查看区块链系统

**url**

http://192.168.203.2:8080/api/blockchains

**method**

GET

**body**

空



## 通道

### 创建通道

**url**

http://192.168.203.2:8080/api/channels

**method**

PUT

**body**

```
{
    "Name": "mychannel",
    "Peers": [
        "peer0.org1.example.com"
    ],
    "AnchorPeers": null,
    "BlockchainID": 1,
    "Blockchainname": ""
}
```

### 修改通道

**url**

http://192.168.203.2:8080/api/channels

**method**

POST

**body**

```
{
    "ChannelName":"mychannel",
    "Operation":0,
    "Args":["peer1.org1.example.com","peer2.org1.example.com"],
    "BlockchainID":1
}
```

### 查看通道

**url**

获取所有通道：

http://192.168.203.2:8080/api/channels

获取区块链1中的通道:

http://192.168.203.2:8080/api/channels?blockchainid=1

##### 参数

blockchainid     int

如果参数留空，则返回平台中所有区块链系统的通道

输入blockchainid参数，则返回对应区块链中的通道信息

**method**

GET

## 合约

### 上传合约

**url**

http://192.168.203.2:8080/api/newcontract

**method**

POST

**body**  （如使用Postman，选择Body为form-data可上传文件）

```
{
  file: 选择文件
	BlockchainID:1
  ContractLang:0
  ContractName:test
  ChannelName:mychannel
  ContractDesc:this is a contract
  ContractVersion:1.0.0
}
```

### 安装合约（已废弃）

> 在新版本中，合约会在上传后自动安装。因此无需再次安装。此API被废弃。

**url**

http://192.168.203.2:8080/api/contracts

**method**

POST

**body**

```
{
    "ID": 1,
    "BlockchainID": 1,
    "ContractName": "contract1",
    "ContractVersion": "1.0.0",
    "Operation": 0
}
```

### 实例化合约

**url**

http://192.168.203.2:8080/api/contracts

**method**

POST

**body**

```
{
    "ID": 1,
    "BlockchainID": 1,
    "ContractName": "test",
    "ContractVersion": "1.0.0",
    "Operation": 1,
    "Args": ["init",  "a", "100",  "b",  "100"]
}
```

### 查看合约

**url**

获得所有合约：

http://192.168.203.2:8080/api/contracts

获得区块链1中的合约：

http://192.168.203.2:8080/api/contracts?blockchainid=1

获得区块链1中mychannel的合约：

http://192.168.203.2:8080/api/contracts?blockchainid=1&channelname=mychannel

##### 参数

blockchainid     int

如果参数留空，则返回平台中所有区块链系统的合约

输入blockchainid参数，则返回对应区块链中的合约



channelname    string

此参数在指定了区块链ID后才有效

如果参数留空，则返回对应区块链中所有合约

输入channelname参数，则返回对应通道中的合约

**method**

GET

### 调用合约(Query)   

**url**

http://192.168.203.2:8080/api/contractcall

**method**

post

**body**

```
{
    "ID": 1,
    "BlockchainID": 1,
    "ContractName": "test",
    "ContractVersion": "1.0.0",
    "InvokeType": 0,
    "Args": ["query",  "a"]
}
```

### 调用合约(Invoke)   

**url**

http://192.168.203.2:8080/api/contractcall

**method**

POST

**body**

```
{
    "ID": 1,
    "BlockchainID": 1,
    "ContractName": "test",
    "ContractVersion": "1.0.0",
    "InvokeType": 1,
    "Args": ["invoke", "a", "b", "30"]
}
```

### 获取某一个区块链系统的合约调用记录

**url**

http://192.168.203.2:8080/api/contractlog/1

**method **

GET

**body**

无

### 获取所有合约调用记录

**url**

http://192.168.203.2:8080/api/contractlog

**method**

GET

**body**

无




## 概览

### 获取概览

> 目前获取概览会提供Baas中的区块链系统总数，正在运行的区块链系统总数，组织总数，Orderer组织总数，Peer组织总数，通道总数，合约总数，区块总数。以及分别展示每个区块链系统的数据等。概览页面不需要将所有数据展示出来，展示一下整体的总数就行

**url**

http://192.168.203.2:8080/api/summary

**method**

GET

**body**

无
