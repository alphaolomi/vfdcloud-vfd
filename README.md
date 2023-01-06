
[vfdcloud/vfd](https://github.com/vfdcloud/vfd) a fully compliant golang implementation of the VFD API as specified by the
Tanzania Revenue Authority (TRA).

The library currently supports

- Client Registration
- Token Fetching
- Receipt Posting
- Z Report Posting

You can also generate the receipt/z report in form of xml file as specified by the TRA. Also you can specify the location of the xml file
and post it to TRA.

The library is used to power multiple platform in real world scenarios. So yes it has been tested in some real world environment. Still some improvements can be done and contributions
are welcome.


### Installation

```bash
 go get github.com/vfdcloud/vfd
```


### Usage

```go
package main

import (
    "fmt"
    "github.com/vfdcloud/vfd"
	"github.com/vfdcloud/vfd/pkg/env"
)

// Fetching Access Token
func main(){
	tokenURL := vfd.RequestURL(env.DEV,vfd.FetchTokenAction)
	request := &vfd.TokenRequest{
		Username:  "",
		Password:  "",
		GrantType: "",
	}
	response, err := vfd.FetchToken(context.Background(), tokenURL, request)
	if err != nil {
		fmt.Printf("failed to fetch token: %v", err)
		os.Exit(1)
	}

	fmt.Printf("token: %+v", response)
}

```