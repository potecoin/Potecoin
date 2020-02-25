# GO-PTC
---

### About
The PTC repository is based on go-ethereum which contains protocol changes to support the PTC protocol and a few other distinct features. This implements the PTC cryptocurrency, which maintains a separate ledger from the ETHEREUM network, for several reasons, the most immediate of which is that the consensus protocol is different.

### Highlights

+ High-performace TPS
+ Highly-regulated network hierarchy

+ 高性能 TPS
+ 一对多交易

### Documentation Guide

+ If you want to know more about PTC MINING , please check out ['Guide to MINING'](http://potecoin.org/)

+ If you want to know more about PTC WALLET, please check out ['Guide to BlockChain WALLET'](http://http://69.30.224.136:81)


### 文档指引

+ 如果你想了解 PTC 详细挖矿操作方法，请查看[《PTC 挖矿教程》] (http://potecoin.org/)
+ 如果你想了解 PTC 网页钱包的操作方法，请查看[《PTC 网页钱包教程》] (http://http://69.30.224.136:81)


### Getting Started
Welcome! This guide is intended to get you running on the PTC testnet. To ensure your client behaves gracefully throughout the setup process, please check your system meets the following requirements:


| OS      | Windows, Mac, Linux                                 |
|---------|----------------------------------------------|
| CPU     | 4 Core (Intel(R) Xeon(R) CPU X5670 @2.93GHz) |
| RAM     | 8G                                           |
| Free HD | 500G                                         |



### Build from Source

First of all, you need to clone the source code from PTC repository:

Git clone https://github.com/potecoin/Potecoin.git, or

wget https://github.com/potecoin/Potecoin/archive/master.zip

Building gptc requires both a Go (version 1.7 or later) and a C compiler. You can install them using your favourite package manager. Once the dependencies are installed, run:

    cd Potecoin
    make gptc
or, to build the full suite of utilities:

    make all


### Docker Quick Start

One of the quickest ways to get PTC up and running on your machine is by using Docker:

    docker build –t image_name

    docker run -d  -p 8341:8341 –p 50505:50505 –p 50505:50505 imagename --fast --cache=512

This will start gptc in fast-sync mode with a DB memory allowance of 1GB just as the above command does.  It will also create a persistent volume in your home directory for saving your blockchain as well as map the default ports. There is also an `alpine` tag available for a slim version of the image.

Do not forget `--rpcaddr 0.0.0.0`, if you want to access RPC from other containers and/or hosts. By default, `gptc` binds to the local interface and RPC endpoints is not accessible from the outside.


### Operating a private network

Maintaining your own private network is more involved as a lot of configurations taken for granted in
the official networks need to be manually set up.

#### Defining the private genesis state

./gptc --datadir ./chaindata/ init ./root/genesis.json
```

#### Starting up your member nodes

step 1: With the bootnode operational and externally reachable (you can try `telnet <ip> <port>` to ensure
it's indeed reachable), start every subsequent gptc node pointed to the bootnode for peer discovery
via the `--bootnodes` flag. It will probably also be desirable to keep the data directory of your
private network separated, so do also specify a custom `--datadir` flag.

```
$ gptc --datadir=path/to/custom/data/folder --bootnodes=<bootnode-enode-url-from-above>

step 2: start up all validator nodes （see genesis.json for configurations of validators）

step 3: start up all broadcast nodes (see man.json for configurations of broadcast nodes)

step 4: start up all miner nodes (see genesis.json for configurations of validators)


#### Execute the following command:

./gptc --identity "MyNodeName" --datadir ./chaindata/ --rpc --rpcaddr 0.0.0.0 --rpccorsdomain "*" --networkid 1 --password ./chaindata/password.txt

Note:  password.txt contains your password of the wallet, which can also be placed under /chaindata.



#### License
Copyright 2020 The PTC Authors

The go-ptc library is licensed under MIT.

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
