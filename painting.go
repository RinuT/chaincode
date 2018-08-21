package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

// SimpleChaincode example simple Chaincode implementation
type OrderChaincode struct {
}

type order struct {
	ObjectType      string `json:"docType"` //docType is used to distinguish the various types of objects in state database
	OrderId         string `json:"OrderId"` //the fieldtags are needed to keep case from bouncing around
	Buyer           string `json:"Buyer"`
	Seller          string `json:"Seller"`
	CurrentLocation string `json:"CurrentLocation"`
	DestinationCity string `json:"DestinationCity"` //the fieldtags are needed to keep case from bouncing around
	OriginCity      string `json:"OriginCity"`
	OrderCondition  string `json:"OrderCondition"`
	Temperature     string `json:"Temperature"`
	Humidity        string `json:"Humidity"`
	Luminosity      string `json:"Luminosity"`
}

type getTemperature struct {
	Temperature string `json:"Temperature"`
	// LastName  string `json:"last_name"`
}

func main() {
	err := shim.Start(new(OrderChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}

func (t *OrderChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	function, args := stub.GetFunctionAndParameters()
	fmt.Println("init is running " + function)

	OrderId := args[0]
	Buyer := args[1]
	Seller := args[2]
	CurrentLocation := args[3]
	DestinationCity := args[4]
	OriginCity := args[5]
	OrderCondition := args[6]
	Temperature := args[7]
	Humidity := args[8]
	Luminosity := args[9]

	// ==== Create product object and marshal to JSON ====
	objectType := "order"
	order := &order{objectType, OrderId, Buyer, Seller, CurrentLocation, DestinationCity, OriginCity, OrderCondition, Temperature, Humidity, Luminosity}
	orderJSONasBytes, err := json.Marshal(order)
	if err != nil {
		return shim.Error(err.Error())
	}

	err = stub.PutState(OrderId, orderJSONasBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	fmt.Println("- end init product")

	return shim.Success(nil)
}

func (t *OrderChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	function, args := stub.GetFunctionAndParameters()
	fmt.Println("invoke is running " + function)

	// Handle different functions
	if function == "registerOrder" {
		return t.registerOrder(stub, args)
	} else if function == "getOrderDetails" { //read a product
		return t.getOrderDetails(stub, args)
	} else if function == "updateTemparature" {
		return t.updateTemparature(stub, args)
	} else if function == "queryTemperature" {
		return t.queryTemperature(stub, args)
	} else if function == "updateHumidity" {
		return t.updateHumidity(stub, args)
	} else if function == "updateLuminosity" {
		return t.updateLuminosity(stub, args)
	}

	fmt.Println("invoke did not find func: " + function) //error
	return shim.Error("Received unknown function invocation")
}

func (t *OrderChaincode) registerOrder(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error

	if len(args) != 10 {
		return shim.Error("Incorrect number of arguments. Expecting 10")
	}

	// ==== Input sanitation ====
	fmt.Println("- start register order")

	OrderId := args[0]
	Buyer := args[1]
	Seller := args[2]
	CurrentLocation := args[3]
	DestinationCity := args[4]
	OriginCity := args[5]
	OrderCondition := args[6]
	Temperature := args[7]
	Humidity := args[8]
	Luminosity := args[9]

	// ==== Check if order already exists ====
	orderAsBytes, err := stub.GetState(OrderId)
	if err != nil {
		return shim.Error("Failed to get product: " + err.Error())
	} else if orderAsBytes != nil {
		fmt.Println("This order already exists: " + OrderId)
		return shim.Error("This order already exists: " + OrderId)
	}

	// ==== Create order object and marshal to JSON ====
	objectType := "order"
	order := &order{objectType, OrderId, Buyer, Seller, CurrentLocation, DestinationCity, OriginCity, OrderCondition, Temperature, Humidity, Luminosity}
	fmt.Println(order)
	//marshal- convert go datatypes to json format

	// productJSONasString := `{"docType":"product",  "uuid": "` + uuid + `", "make": "` + make + `", "material_location": ` + material_location + `, "product_status": "` + product_status + `, "shipment_status": ` + shipment_status + `"}`
	// productJSONasBytes := []byte(productJSONasString)

	orderJSONasBytes, err := json.Marshal(order)
	fmt.Println(orderJSONasBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	// === Save product to state ===
	err = stub.PutState(OrderId, orderJSONasBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	fmt.Println("- end register order")
	return shim.Success(nil)
}

func (t *OrderChaincode) getOrderDetails(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var OrderId, jsonResp string
	var err error

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting name of the marble to query")
	}

	OrderId = args[0]
	valAsbytes, err := stub.GetState(OrderId) //get the marble from chaincode state
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + OrderId + "\"}"
		return shim.Error(jsonResp)
	} else if valAsbytes == nil {
		jsonResp = "{\"Error\":\"Marble does not exist: " + OrderId + "\"}"
		return shim.Error(jsonResp)
	}

	return shim.Success(valAsbytes)
}

func (t *OrderChaincode) updateTemparature(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) < 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	OrderId := args[0]
	newStatus := args[1]
	fmt.Println("- update product status ", OrderId, newStatus)

	productAsBytes, err := stub.GetState(OrderId)
	if err != nil {
		return shim.Error("Failed to get product:" + err.Error())
	} else if productAsBytes == nil {
		return shim.Error("product does not exist")
	}

	productToUpdate := order{}
	err = json.Unmarshal(productAsBytes, &productToUpdate) //unmarshal it aka JSON.parse()
	if err != nil {
		return shim.Error(err.Error())
	}
	productToUpdate.Temperature = newStatus //change the owner

	productJSONasBytes, _ := json.Marshal(productToUpdate)
	err = stub.PutState(OrderId, productJSONasBytes) //rewrite the marble
	if err != nil {
		return shim.Error(err.Error())
	}

	fmt.Println("- end updateProductStatus (success)")
	return shim.Success(nil)
}

func (t *OrderChaincode) updateHumidity(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) < 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	OrderId := args[0]
	newStatus := args[1]
	fmt.Println("- update product status ", OrderId, newStatus)

	productAsBytes, err := stub.GetState(OrderId)
	if err != nil {
		return shim.Error("Failed to get product:" + err.Error())
	} else if productAsBytes == nil {
		return shim.Error("product does not exist")
	}

	productToUpdate := order{}
	err = json.Unmarshal(productAsBytes, &productToUpdate) //unmarshal it aka JSON.parse()
	if err != nil {
		return shim.Error(err.Error())
	}
	productToUpdate.Humidity = newStatus //change the owner

	productJSONasBytes, _ := json.Marshal(productToUpdate)
	err = stub.PutState(OrderId, productJSONasBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	fmt.Println("- end updateProductStatus (success)")
	return shim.Success(nil)
}

func (t *OrderChaincode) updateLuminosity(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) < 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	OrderId := args[0]
	newStatus := args[1]
	fmt.Println("- update product status ", OrderId, newStatus)

	productAsBytes, err := stub.GetState(OrderId)
	if err != nil {
		return shim.Error("Failed to get product:" + err.Error())
	} else if productAsBytes == nil {
		return shim.Error("product does not exist")
	}

	productToUpdate := order{}
	err = json.Unmarshal(productAsBytes, &productToUpdate) //unmarshal it aka JSON.parse()
	if err != nil {
		return shim.Error(err.Error())
	}
	productToUpdate.Luminosity = newStatus //change the owner

	productJSONasBytes, _ := json.Marshal(productToUpdate)
	err = stub.PutState(OrderId, productJSONasBytes) //rewrite the marble
	if err != nil {
		return shim.Error(err.Error())
	}

	fmt.Println("- end updateProductStatus (success)")
	return shim.Success(nil)
}

func (t *OrderChaincode) queryTemperature(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) < 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	OrderId := args[0]

	fmt.Printf("- start getHistoryForMarble: %s\n", OrderId)

	resultsIterator, err := stub.GetHistoryForKey(OrderId)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing historic values for the marble
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		response, err := resultsIterator.Next()
		fmt.Println(response)

		if err != nil {
			return shim.Error(err.Error())
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"TxId\":")
		buffer.WriteString("\"")
		buffer.WriteString(response.TxId)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Value\":")
		// if it was a delete operation on given key, then we need to set the
		//corresponding value null. Else, we will write the response.Value
		//as-is (as the Value itself a JSON marble)
		if response.IsDelete {
			buffer.WriteString("null")
		} else {
			buffer.WriteString(string(response.Value))
		}
		buffer.WriteString(", \"Timestamp\":")
		buffer.WriteString("\"")
		buffer.WriteString(time.Unix(response.Timestamp.Seconds, int64(response.Timestamp.Nanos)).String())
		buffer.WriteString("\"")

		buffer.WriteString(", \"IsDelete\":")
		buffer.WriteString("\"")
		buffer.WriteString(strconv.FormatBool(response.IsDelete))
		buffer.WriteString("\"")

		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true

	}
	buffer.WriteString("]")

	fmt.Printf("- returning history of the order:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}
