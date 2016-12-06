package main

import (
	"errors"
	"fmt"
	"strconv"

	"hyperledger/cci/appinit"
	"hyperledger/cci/org/hyperledger/chaincode/example02"
	"hyperledger/ccs"

	"github.com/golang/protobuf/proto"
	"github.com/hyperledger/fabric/core/chaincode/shim"
)

type ChaincodeExample struct {
}

// Called to initialize the chaincode
func (t *ChaincodeExample) Init(stub shim.ChaincodeStubInterface, param *appinit.Init) error {

	var err error

	fmt.Printf("Aval = %d, Bval = %d\n", param.PartyA.Value, param.PartyB.Value)

	// Write the state to the ledger
	err = t.PutState(stub, param.PartyA)
	if err != nil {
		return err
	}

	err = t.PutState(stub, param.PartyB)
	if err != nil {
		return err
	}

	return nil
}

// Transaction makes payment of X units from A to B
func (t *ChaincodeExample) MakePayment(stub shim.ChaincodeStubInterface, param *example02.PaymentParams) error {

	var err error

	// Get the state from the ledger
	src, err := t.GetState(stub, param.PartySrc)
	if err != nil {
		return err
	}

	dst, err := t.GetState(stub, param.PartyDst)
	if err != nil {
		return err
	}

	// Perform the execution
	X := int(param.Amount)
	src = src - X
	dst = dst + X
	fmt.Printf("Aval = %d, Bval = %d\n", src, dst)

	// Write the state back to the ledger
	err = stub.PutState(param.PartySrc, []byte(strconv.Itoa(src)))
	if err != nil {
		return err
	}

	err = stub.PutState(param.PartyDst, []byte(strconv.Itoa(dst)))
	if err != nil {
		return err
	}

	return nil
}

// Deletes an entity from state
func (t *ChaincodeExample) DeleteAccount(stub shim.ChaincodeStubInterface, param *example02.Entity) error {

	// Delete the key from the state in ledger
	err := stub.DelState(param.Id)
	if err != nil {
		return errors.New("Failed to delete state")
	}

	return nil
}

// Query callback representing the query of a chaincode
func (t *ChaincodeExample) CheckBalance(stub shim.ChaincodeStubInterface, param *example02.Entity) (*example02.BalanceResult, error) {
	var err error

	// Get the state from the ledger
	val, err := t.GetState(stub, param.Id)
	if err != nil {
		return nil, err
	}

	fmt.Printf("Query Response: %d\n", val)
	return &example02.BalanceResult{Balance: *proto.Int32(int32(val))}, nil
}

func main() {
	self := &ChaincodeExample{}
	interfaces := ccs.Interfaces{
		"org.hyperledger.chaincode.example02": self,
		"appinit": self,
	}

	err := ccs.Start(interfaces) // Our one instance implements both Transactions and Queries interfaces
	if err != nil {
		fmt.Printf("Error starting example chaincode: %s", err)
	}
}

//-------------------------------------------------
// Helpers
//-------------------------------------------------
func (t *ChaincodeExample) PutState(stub shim.ChaincodeStubInterface, party *appinit.Party) error {
	return stub.PutState(party.Entity, []byte(strconv.Itoa(int(party.Value))))
}

func (t *ChaincodeExample) GetState(stub shim.ChaincodeStubInterface, entity string) (int, error) {
	bytes, err := stub.GetState(entity)
	if err != nil {
		return 0, errors.New("Failed to get state")
	}
	if bytes == nil {
		return 0, errors.New("Entity not found")
	}

	val, _ := strconv.Atoi(string(bytes))
	return val, nil
}

    Contact GitHub API Training Shop Blog About 

    © 2016 GitHub, Inc. Terms Privacy Security Status Help 

