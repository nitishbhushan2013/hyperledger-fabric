package main

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

// it will implement chaincodde interface
type PcXchng struct {
}

//Data structure to store details on blockchain
type PC struct { // field start with capital letter for encoding/json
	Snumber string // unique key
	Series  string
	Other   string
	Status  string
}

// implement chaincode interface
// Init
func (c *PcXchng) Init(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Success(nil) // nothing to initialize
}

// implement Invoke method
func (c *PcXchng) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	// get user input, and invoke function based on argument

	/*
		// GetFunctionAndParameters returns the first argument as the function
		// name and the rest of the arguments as parameters in a string array.
		// Only use GetFunctionAndParameters if the client passes arguments intended
		// to be used as strings.
		GetFunctionAndParameters() (string, []string)
	*/
	fun, args := stub.GetFunctionAndParameters()

	switch fun {
	case "createPC": // computer is produced
		return c.createPC(stub, args)
	case "queryStock": // get the stock of computers
		return c.queryStock(stub, args)
	case "queryDetail": // get detail about computer
		return c.queryDetail(stub, args)
	case "buyPC": // buy a computer
		return c.updateStatus(stub, args, "bought")
	case "handbackPC": // returned the computer back to producer
		return c.updateStatus(stub, args, "returned")
	default: // error, display user friendly message
		return shim.Error("please provide the correct function to invoke. Available functions are createPC,queryStock,queryDetail,updateStatus,buyPC,handbackPC ")
	}
}

// lets define each of the available methods

func (c *PcXchng) createPC(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	// check the arguments
	if len(args) != 3 {
		shim.Error("createPC expects three arguments : Snumber,Series,Other ")
	}

	// create a pc and marked it 'available'
	pc := PC{
		Snumber: args[0],
		Series:  args[0],
		Other:   args[0],
		Status:  "available",
	}
	// Save it to ledger
	/*
		// PutState puts the specified `key` and `value` into the transaction's
		// writeset as a data-write proposal. PutState doesn't effect the ledger
		// until the transaction is validated and successfully committed.
		// Simple keys must not be an empty string and must not start with null
		// character (0x00), in order to avoid range query collisions with
		// composite keys, which internally get prefixed with 0x00 as composite
		// key namespace.
		PutState(key string, value []byte) error
		}
	*/
	// get the []byte
	/*
		func Marshal(v interface{}) ([]byte, error)
	*/
	pcBytes, err := json.Marshal(pc)
	if err != nil {
		return shim.Error(err.Error())
	}

	err = stub.PutState(pc.Snumber, pcBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

func (c *PcXchng) queryStock(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 2 {
		return shim.Error("queryStock expects two arguments : startIndex, endIndex")
	}

	startIndex, endIndex := " ", " "
	startIndex = args[0]
	endIndex = args[1]

	/*
		// GetStateByRange returns a range iterator over a set of keys in the
		// ledger. The iterator can be used to iterate over all keys
		// between the startKey (inclusive) and endKey (exclusive).
		// The keys are returned by the iterator in lexical order. Note
		// that startKey and endKey can be empty string, which implies unbounded range
		// query on start or end.
		// Call Close() on the returned StateQueryIteratorInterface object when done.
		// The query is re-executed during validation phase to ensure result set
		// has not changed since transaction endorsement (phantom reads detected).
		GetStateByRange(startKey, endKey string) (StateQueryIteratorInterface, error)
	*/
	recordIterator, err := stub.GetStateByRange(startIndex, endIndex)
	if err != nil {
		return shim.Error(err.Error())
	}

	defer recordIterator.Close()

	keys := "\n"
	for recordIterator.HasNext() {
		response, err := recordIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		keys += response.Key + "\n"
	}

	fmt.Println(keys)
	return shim.Success([]byte(keys))
}

func (c *PcXchng) queryDetail(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("queryDetail expects one argument : Snumber ")
	}
	/*
		func (stub *ChaincodeStub) GetState(key string) ([]byte, error)
	*/
	pcDetailsBytes, err := stub.GetState(args[0])
	if err != nil {
		shim.Error(err.Error())
	}
	/*
		To decode JSON data we use the Unmarshal function.

		func Unmarshal(data []byte, v interface{}) error
		We must first create a place where the decoded data will be stored

		var m Message
		and call json.Unmarshal, passing it a []byte of JSON data and a pointer to m

		err := json.Unmarshal(b, &m)

		}
	*/
	var pc PC
	err = json.Unmarshal(pcDetailsBytes, &pc)
	if err != nil {
		shim.Error(err.Error())
	}
	fmt.Println(pc)
	return shim.Success([]byte("pc.Snumber>>" + pc.Snumber + ",pc.Series>>" + pc.Series + ",pc.Other>>" + pc.Other + ",pc.Status>>" + pc.Status))
}

func (c *PcXchng) updateStatus(stub shim.ChaincodeStubInterface, args []string, status string) pb.Response {
	if len(args) != 1 {
		return shim.Error("queryDetail expects one argument : Snumber ")
	}
	/*
		func (stub *ChaincodeStub) GetState(key string) ([]byte, error)
	*/
	pcDetailsBytes, err := stub.GetState(args[0])
	if err != nil {
		shim.Error(err.Error())
	}

	// get the PC
	var pc PC
	err = json.Unmarshal(pcDetailsBytes, &pc)
	if err != nil {
		shim.Error(err.Error())
	}

	// update the status
	pc.Status = status

	// convert PC to bytes[] for storing in blockchain
	pcBytes, err := json.Marshal(pc)
	if err != nil {
		return shim.Error(err.Error())
	}

	// write into write set of ledger
	err = stub.PutState(pc.Snumber, pcBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

func main() {
	err := shim.Start(new(PcXchng))
	if err != nil {
		fmt.Printf("Error starting chaincode sample: %s", err)
	}
}
