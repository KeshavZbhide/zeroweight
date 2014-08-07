ZeroWeight
==========
Zeroweight is a super lightweight bitcoin wallet written in golang. 
It can be used as a package/api or as a stand alone wallet.

### Installing as stand alone wallet
    $export $PATH=$GOPATH/bin:$PATH
    $go get github.com/KeshavZbhide/zeroweight
    $zeroweight 


find . -type f | xargs sed -i "" "s|package zeroweight|package zeroweight|g"
