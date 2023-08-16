package main

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type SmartContract struct {
	contractapi.Contract
}

type Food struct {
	Farmer  string `json:"farmer"`
	Variety string `json:"variety"`
}

func (s *SmartContract) Set(ctx contractapi.TransactionContextInterface, foodId string, farmer string, variety string) error {

	food := Food{
		Farmer:  farmer,
		Variety: variety,
	}

	foodAsBytes, err := json.Marshal(food)

	if err != nil {
		fmt.Errorf("Marshal error: %s", err.Error())
		return err
	}

	return ctx.GetStub().PutState(foodId, foodAsBytes)
}

func (s *SmartContract) Query(ctx contractapi.TransactionContextInterface, foodId string) (*Food, error) {
	foodAsBytes, err := ctx.GetStub().GetState(foodId)

	if err != nil {
		return nil, fmt.Errorf("failed to read from world state. %s", err.Error())
	}

	if foodAsBytes == nil {
		return nil, fmt.Errorf("%s does not exist", foodId)
	}

	food := new(Food)

	err = json.Unmarshal(foodAsBytes, food)

	if err != nil {
		return nil, fmt.Errorf("marshal error: %s", err.Error())
	}

	return food, nil
}

func main() {
	chainconde, err := contractapi.NewChaincode(new(SmartContract))

	if err != nil {
		fmt.Errorf("error create foodcontrol chaincode: %s", err.Error())
		return
	}

	if err := chainconde.Start(); err != nil {
		fmt.Errorf("error create foodcontrol chaincode: %s", err.Error())
	}

}
