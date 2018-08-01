package main

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

type SimpleChaincode struct {
}

type product struct {
	ObjectType        string `json:"docType"`
	uuid              string `json:"uuid"`     //docType is used to distinguish the various types of objects in state database
	material          string `json:"material"` //the fieldtags are needed to keep case from bouncing around
	make              string `json:"make"`
	material_location string `json:"material_location"`
	shipment_status   string `json:"shipment_status"`
	product_status    string `json:"product_status"`
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
	if function == "updateShipmentStatus" {
		return t.updateShipmentStatus(stub, args)
	}else if function == "updateproductStatus" {
		return t.updateproductStatus(stub, args)
	}
	
	fmt.Println("invoke did not find func: " + function) //error
	return shim.Error("Received unknown function invocation")
}


func (t *SimpleChaincode) updateShipmentStatus(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) < 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	uuid := args[0]
	newStatus := args[1]
	fmt.Println("- update shipment status ", uuid, newStatus)

	productAsBytes, err := stub.GetState(uuid)
	if err != nil {
		return shim.Error("Failed to get product:" + err.Error())
	} else if productAsBytes == nil {
		return shim.Error("product does not exist")
	}

	productToUpdate := product{}
	err = json.Unmarshal(productAsBytes, &productToUpdate) //unmarshal it aka JSON.parse()
	if err != nil {
		return shim.Error(err.Error())
	}
	productToUpdate.shipment_status = newStatus //change the owner

	productJSONasBytes, _ := json.Marshal(productToUpdate)
	err = stub.PutState(uuid, productJSONasBytes) //rewrite the marble
	if err != nil {
		return shim.Error(err.Error())
	}

	fmt.Println("- end updateShipmentStatus (success)")
	return shim.Success(nil)
}

func (t *SimpleChaincode) updateproductStatus(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) < 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	uuid := args[0]
	newStatus := args[1]
	fmt.Println("- update product status ", uuid, newStatus)

	productAsBytes, err := stub.GetState(uuid)
	if err != nil {
		return shim.Error("Failed to get product:" + err.Error())
	} else if productAsBytes == nil {
		return shim.Error("product does not exist")
	}

	productToUpdate := product{}
	err = json.Unmarshal(productAsBytes, &productToUpdate) //unmarshal it aka JSON.parse()
	if err != nil {
		return shim.Error(err.Error())
	}
	productToUpdate.product_status = newStatus //change the owner

	productJSONasBytes, _ := json.Marshal(productToUpdate)
	err = stub.PutState(uuid, productJSONasBytes) //rewrite the marble
	if err != nil {
		return shim.Error(err.Error())
	}

	fmt.Println("- end updateProductStatus (success)")
	return shim.Success(nil)
}