package main

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

type SimpleChaincode struct {
}

type product struct {

	ObjectType string `json:"docType"`
	uuid string `json:"uuid"` //docType is used to distinguish the various types of objects in state database
	material       string `json:"material"`    //the fieldtags are needed to keep case from bouncing around
	make      string `json:"make"`
	material_location       string    `json:"material_location"`
	shipment_status      string `json:"shipment_status"`
	product_status      string `json:"product_status"`
}
func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}

func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	function, args := stub.GetFunctionAndParameters()
	fmt.Println("init is running " + function)

	uuid := args[0]
	material := args[1]
	make := args[2]
	material_location := args[3]
	shipment_status := args[4]
	product_status := args[5]


	// ==== Create product object and marshal to JSON ====
	objectType := "product"
	product := &product{objectType, uuid, material, make, material_location, shipment_status, product_status}
	productJSONasBytes, err := json.Marshal(product)
	if err != nil {
		return shim.Error(err.Error())
	}


	err = stub.PutState(uuid, productJSONasBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	fmt.Println("- end init product")

	return shim.Success(nil)
}

func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	function, args := stub.GetFunctionAndParameters()
	fmt.Println("invoke is running " + function)

	// Handle different functions
	if function == "createProduct" { 
		return t.createProduct(stub, args)
	}else if function == "searchProduct" { //read a product
		return t.searchProduct(stub, args)
	} 

	fmt.Println("invoke did not find func: " + function) //error
	return shim.Error("Received unknown function invocation")
}

func (t *SimpleChaincode) createProduct(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error

	if len(args) != 6 {
		return shim.Error("Incorrect number of arguments. Expecting 6")
	}

	// ==== Input sanitation ====
	fmt.Println("- start init product")

	uuid := args[0]
	material := args[1]
	make := args[2]
	material_location := args[3]
	shipment_status := args[4]
	product_status := args[5]


	// ==== Check if product already exists ====
	productAsBytes, err := stub.GetState(uuid)
	if err != nil {
		return shim.Error("Failed to get product: " + err.Error())
	} else if productAsBytes != nil {
		fmt.Println("This product already exists: " + uuid)
		return shim.Error("This product already exists: " + uuid)
	}

	// ==== Create product object and marshal to JSON ====
	objectType := "product"
	product := &product{objectType, uuid, material, make, material_location, shipment_status, product_status}
	productJSONasBytes, err := json.Marshal(product)
	if err != nil {
		return shim.Error(err.Error())
	}
	

	// === Save product to state ===
	err = stub.PutState(uuid, productJSONasBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	fmt.Println("- end init product")
	return shim.Success(nil)
}

func (t *SimpleChaincode) searchProduct(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	fmt.println("invoking search function")
	var uuid, jsonResp string
	var err error

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting product uuid")
	}

	uuid = args[0]
	valAsbytes, err := stub.GetState(uuid) //get the product from chaincode state
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + uuid + "\"}"
		return shim.Error(jsonResp)
	} else if valAsbytes == nil {
		jsonResp = "{\"Error\":\"product does not exist: " + uuid + "\"}"
		return shim.Error(jsonResp)
	}
	fmt.println("Succesffully searched product")
	return shim.Success(valAsbytes)
}