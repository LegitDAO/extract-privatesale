package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"math/big"
	"strings"

	"github.com/ybbus/jsonrpc/v2"
)

//use the rigth web3api endpoint
const endpoint = "https://bsc-mainnet.web3api.com"

// some constants:
const fromAddress = "0x0000000000000000000000000000000000000000"
const toAddress = "0x22c32f56f1e98cbdbc97761d471da3d986686378"
const userClaimable = "0xe22d4f4a000000000000000000000000"
const userInvestment = "0xc52c5c88000000000000000000000000"

type jsonParams struct {
	From string `json:"from"`
	To   string `json:"to"`
	Data string `json:"data"`
}

func fetchAmount(addrType string, userAddress string) *big.Int {
	address := strings.Replace(userAddress, "0x", "", 1)
	rpcClient := jsonrpc.NewClient(endpoint)
	response, err := rpcClient.Call("eth_call", []interface{}{
		&jsonParams{
			From: fromAddress,
			To:   toAddress,
			Data: fmt.Sprintf("%s%s", addrType, address),
		},
		"latest",
	})

	if err != nil {
		panic(err)
	}

	if response.Error != nil {
		panic(errors.New(response.Error.Message))
	}

	cleaned := strings.TrimLeft(response.Result.(string), "0x")
	bigInt := new(big.Int)
	bigInt.SetString(cleaned, 16)
	if err != nil {
		panic(err)
	}

	return bigInt
}

func main() {
	data, err := ioutil.ReadFile("addresses.data")
	if err != nil {
		panic(err)
	}

	uniqueMap := map[string]string{}
	content := string(data)
	addressesList := strings.Split(content, "\n")
	for _, oneAddress := range addressesList {
		if oneAddress == "" {
			continue
		}

		uniqueMap[oneAddress] = oneAddress
	}

	uniqueList := []string{}
	for _, oneAddress := range uniqueMap {
		uniqueList = append(uniqueList, oneAddress)
	}

	fmt.Printf("\n")
	for index, oneAddress := range uniqueList {
		claimable := fetchAmount(userClaimable, oneAddress)
		investment := fetchAmount(userInvestment, oneAddress)
		if claimable.String() == "0" && investment.String() == "0" {
			continue
		}

		fmt.Printf("%d) %s\nclaimable: %d\ninvestment: %d\n----\n", index, oneAddress, claimable, investment)
	}
}

//curl https://bsc-mainnet.web3api.com/v1/EVC3765X3YTKGKR77SUIYSEW6Y5TP5AX9Y -X POST -H "Content-Type: application/json"  -d "{\"jsonrpc\":\"2.0\",\"method\":\"eth_call\",\"params\": [{\"from\": \"0x0000000000000000000000000000000000000000\",\"to\": \"0x22c32f56f1e98cbdbc97761d471da3d986686378\",\"data\": \"0xc52c5c8800000000000000000000000002107c1794c7074513874f4c608e2d00cdee84f5\"}, \"latest\"],\"id\":1}"
