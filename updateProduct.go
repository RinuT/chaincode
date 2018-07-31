package main

import (
	"encoding/json"
	"fmt"
	"strings"

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
	}else if function == "searchPro" { 
		return t.searchPro(stub, args)
	}else if function == "updateShipmentStatus" { 
		return t.updateShipmentStatus(stub, args)
	}

	fmt.Println("invoke did not find func: " + function) //error
	return shim.Error("Received unknown function invocation")
}

func (t *SimpleChaincode) updateShipmentStatus(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) < 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	uuid := args[0]
	newStatus := strings.ToHigher(args[1])
	fmt.Println("- update shipment status ", uuid, newStatus)

	productAsBytes, err := stub.GetState(uuid)
	if err != nil {
		return shim.Error("Failed to get product:" + err.Error())
	} else if productAsBytes == nil {
		return shim.Error("product does not exist")
	}

	productToUpdate := marble{}
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
