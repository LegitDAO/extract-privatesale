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

// our contract:
const fromAddress = "0x0000000000000000000000000000000000000000"
const toAddress = "0x22c32f56f1e98cbdbc97761d471da3d986686378"
const userClaimable = "0xe22d4f4a000000000000000000000000"
const userInvestment = "0xc52c5c88000000000000000000000000"

// proxy contract:
const proxyFromAddress = "0x0000000000000000000000000000000000000000"
const proxyToAddress = "0xe4bc5abd68ffefe84421d811514e3605e4a07e06"
const referralPerUser = "0x8b56ad56000000000000000000000000"

type jsonParams struct {
	From string `json:"from"`
	To   string `json:"to"`
	Data string `json:"data"`
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
		units := fetchAmount(userClaimable, oneAddress)
		investment := fetchAmount(userInvestment, oneAddress)
		if units.String() == "0" && investment.String() == "0" {
			continue
		}

		referrals := fetchReferral(referralPerUser, oneAddress)
		fmt.Printf("%d) %s\nunits: %d\nreferrals: %v\n----\n", index, oneAddress, units, strings.Join(referrals, ",\n"))
	}
}

func fetchReferral(fnAddr string, userAddress string) []string {
	address := strings.Replace(userAddress, "0x", "", 1)
	rpcClient := jsonrpc.NewClient(endpoint)
	response, err := rpcClient.Call("eth_call", []interface{}{
		&jsonParams{
			From: proxyFromAddress,
			To:   proxyToAddress,
			Data: fmt.Sprintf("%s%s", fnAddr, address),
		},
		"latest",
	})

	if err != nil {
		panic(err)
	}

	if response.Error != nil {
		panic(errors.New(response.Error.Message))
	}

	list := []string{}
	addrLength := 64
	cleaned := strings.TrimLeft(response.Result.(string), "0x")[66:]
	amount := len(cleaned) / addrLength
	for i := 0; i < amount; i++ {
		from := i * addrLength
		to := (i + 1) * addrLength
		str := cleaned[from:to]
		list = append(list, str[24:])
	}

	return list
}

func fetchAmount(fnAddr string, userAddress string) *big.Int {
	address := strings.Replace(userAddress, "0x", "", 1)
	rpcClient := jsonrpc.NewClient(endpoint)
	response, err := rpcClient.Call("eth_call", []interface{}{
		&jsonParams{
			From: fromAddress,
			To:   toAddress,
			Data: fmt.Sprintf("%s%s", fnAddr, address),
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

//curl https://bsc-mainnet.web3api.com/v1/EVC3765X3YTKGKR77SUIYSEW6Y5TP5AX9Y -X POST -H "Content-Type: application/json"  -d "{\"jsonrpc\":\"2.0\",\"method\":\"eth_call\",\"params\": [{\"from\": \"0x0000000000000000000000000000000000000000\",\"to\": \"0x22c32f56f1e98cbdbc97761d471da3d986686378\",\"data\": \"0xc52c5c8800000000000000000000000002107c1794c7074513874f4c608e2d00cdee84f5\"}, \"latest\"],\"id\":1}"
