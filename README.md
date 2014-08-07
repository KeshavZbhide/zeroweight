ZeroWeight
==========
Zeroweight is a super lightweight commandline bitcoin wallet written in golang. 
It can be used as a package/api or as a stand alone wallet.

###Commands:
    # use => $zeroweight createWallet [encryptionKey|password]
    #     => $zeroweight send [toAddress] [BTC amount] [password]
    #     => $zeroweight balance [password]

use `$history -c` to eleminate password stored in bash history

###Installing as stand alone wallet:
    
    $ export $PATH=$GOPATH/bin:$PATH
    $ go get github.com/KeshavZbhide/zeroweight
    $ zeroweight 

make sure your $PATH contains $GOPATH/bin. Or else just install it 
by moving $GOPATH/bin/zeroweight to /usr/local/bin
    
    $ sudo mv $GOPATH/bin/zeroweight /usr/local/bin

###Installing as a golang package:
    $ git clone https://github.com/KeshavZbhide/zeroweight.git
    $ mv zeroweight $GOPATH/src
    $ cd $GOPATH/src/zeroweight
    $ find . -type f | xargs sed -i "" "s|package main|package zeroweight|g"
The last command changes the name of the package to zeroweight, since the initial
version is a command-line wallet and the package is named "main", now

    import "zeroweight"

should work.

####tx.go:
tx.go exports essential functions that can be used construct a minimum wallet.
tx.go does not download the block chain.

###Exported functions:
#####1. Tx(privateKey string, to string, amount float64) (string, error)
Tx builds a hex encoded transaction that can be submited to a bitcoin node.

#####2. Balance(publicKey string) (float64, error)
Balance returns the balance of the address (public key)

#####3. GenRandPrivateKey() string
genrates randomized private key. this should be called to genrate a new wallet address. 

#####4. GetPublicKey(privateKey string) string 
gets the corresponding public key of the private key, usualy genrated by GenRandomPrivateKey.

#####5. SubmitTransaction(tx string) string
Submits the output genrated from Tx to blockchain.info/pushtx. There is no
error values, but if the function is successfull it will return "Transaction submited"

refer to the comments in tx.go for simple use cases.

####warning:
`$zeroweight createWallet` creates a encrypted file 'zeroweight.wal' in user's home 
directory. It contains the private key requrired to make transactions. Losing this 
file will lead to loosing all your bitcoins. BE VERY CAREFULL.

###Licensing and Copyright
Copyright 2014 Keshav Bhide. All rights reserved.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software 
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
 
