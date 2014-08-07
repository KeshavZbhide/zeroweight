ZeroWeight
==========
Zeroweight is a super lightweight bitcoin wallet written in golang. 
It can be used as a package/api or as a stand alone wallet.

### Installing as stand alone wallet.
    $ export $PATH=$GOPATH/bin:$PATH
    $ go get github.com/KeshavZbhide/zeroweight
    $ zeroweight 
make sure your $PATH contains $GOPATH/bin. Else
    $ mv $GOPATH/bin/zeroweight /usr/local/bin

###Installing as a golang package
    $ git clone https://github.com/KeshavZbhide/zeroweight.git
    $ mv zeroweight $GOPATH/src
    $ cd $GOPATH/src/zeroweight
    $ find . -type f | xargs sed -i "" "s|package main|package zeroweight|g"
The last command changes the name of the package to zeroweight, since the initial
version is a command line wallet and is named "main"
