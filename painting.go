//peer chaincode install -p chaincodedev/chaincode/order -n order -v 0
//peer chaincode instantiate -n order -v 0 -c '{"Args":["order","1","abc","xyz","kochi","Pune","Kochi","in_good_condition","100","25","10"]}' -C myc
//peer chaincode invoke -n order -c '{"Args":["registerOrder","2","abc","xyz","kochi","Pune","Kochi","in_good_condition","100","25","10"]}' -C myc
//peer chaincode query -n order -c '{"Args":["getOrderDetails","2"]}' -C myc
//peer chaincode invoke -n order -c '{"Args":["updateTemparature","2","100"]}' -C myc
//peer chaincode invoke -n order -c '{"Args":["updateHumidity","2","123"]}' -C myc
//peer chaincode invoke -n order -c '{"Args":["updateLuminosity","2","23"]}' -C myc
//peer chaincode query -n order -c '{"Args":["queryHistory","2"]}' -C myc

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

//Initialization of the order
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

	// ==== Create order object and marshal to JSON ====
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

	fmt.Println("- end init Order")

	return shim.Success(nil)
}

//invoke function

func (t *OrderChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	function, args := stub.GetFunctionAndParameters()
	fmt.Println("invoke is running " + function)

	// Handle different functions
	if function == "registerOrder" {
		return t.registerOrder(stub, args)
	} else if function == "getOrderDetails" {
		return t.getOrderDetails(stub, args)
	} else if function == "updateTemparature" {
		return t.updateTemparature(stub, args)
	} else if function == "updateHumidity" {
		return t.updateHumidity(stub, args)
	} else if function == "updateLuminosity" {
		return t.updateLuminosity(stub, args)
	} else if function == "queryHistory" {
		return t.queryHistory(stub, args)
	} else if function == "updateOriginCity" {
		return t.updateOriginCity(stub, args)
	} else if function == "updateCurrentLocation" {
		return t.updateCurrentLocation(stub, args)
	} else if function == "updateDestinationCity" {
		return t.updateDestinationCity(stub, args)
	}

	fmt.Println("invoke did not find func: " + function) //error
	return shim.Error("Received unknown function invocation")
}

//Create an order with unique order Id
func (t *OrderChaincode) registerOrder(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error

	// if len(args) != 10 {
	// 	return shim.Error("Incorrect number of arguments. Expecting 10")
	// }

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
		return shim.Error("Failed to register order: " + err.Error())
	} else if orderAsBytes != nil {
		fmt.Println("This order already exists: " + OrderId)
		return shim.Error("This order already exists: " + OrderId)
	}

	// ==== Create order object and marshal to JSON ====
	objectType := "order"
	order := &order{objectType, OrderId, Buyer, Seller, CurrentLocation, DestinationCity, OriginCity, OrderCondition, Temperature, Humidity, Luminosity}
	fmt.Println(order)
	orderJSONasBytes, err := json.Marshal(order)
	fmt.Println(orderJSONasBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	// === Save order to state ===
	err = stub.PutState(OrderId, orderJSONasBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	fmt.Println("- end register order")
	return shim.Success(nil)
}

//search the order details using order Id
func (t *OrderChaincode) getOrderDetails(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var OrderId, jsonResp string
	var err error

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting name of the marble to query")
	}

	OrderId = args[0]
	valAsbytes, err := stub.GetState(OrderId)
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + OrderId + "\"}"
		return shim.Error(jsonResp)
	} else if valAsbytes == nil {
		jsonResp = "{\"Error\":\"Order does not exist: " + OrderId + "\"}"
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
	fmt.Println("- update temperature ", OrderId, newStatus)

	orderAsBytes, err := stub.GetState(OrderId)
	if err != nil {
		return shim.Error("Failed to get order:" + err.Error())
	} else if orderAsBytes == nil {
		return shim.Error("order does not exist")
	}

	orderToUpdate := order{}
	err = json.Unmarshal(orderAsBytes, &orderToUpdate) //unmarshal it aka JSON.parse()
	if err != nil {
		return shim.Error(err.Error())
	}
	orderToUpdate.Temperature = newStatus //change the temperature

	orderJSONasBytes, _ := json.Marshal(orderToUpdate)
	err = stub.PutState(OrderId, orderJSONasBytes) //rewrite the marble
	if err != nil {
		return shim.Error(err.Error())
	}

	fmt.Println("- end updateTemperature (success)")
	return shim.Success(nil)
}

func (t *OrderChaincode) updateHumidity(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) < 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	OrderId := args[0]
	newStatus := args[1]
	fmt.Println("- update humidity ", OrderId, newStatus)

	orderAsBytes, err := stub.GetState(OrderId)
	if err != nil {
		return shim.Error("Failed to get order:" + err.Error())
	} else if orderAsBytes == nil {
		return shim.Error("order does not exist")
	}

	orderToUpdate := order{}
	err = json.Unmarshal(orderAsBytes, &orderToUpdate) //unmarshal it aka JSON.parse()
	if err != nil {
		return shim.Error(err.Error())
	}
	orderToUpdate.Humidity = newStatus //change the humidity

	orderJSONasBytes, _ := json.Marshal(orderToUpdate)
	err = stub.PutState(OrderId, orderJSONasBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	fmt.Println("- end updateHumidity (success)")
	return shim.Success(nil)
}

func (t *OrderChaincode) updateLuminosity(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) < 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	OrderId := args[0]
	newStatus := args[1]
	fmt.Println("- update Luminosity ", OrderId, newStatus)

	orderAsBytes, err := stub.GetState(OrderId)
	if err != nil {
		return shim.Error("Failed to get order:" + err.Error())
	} else if orderAsBytes == nil {
		return shim.Error("order does not exist")
	}

	orderToUpdate := order{}
	err = json.Unmarshal(orderAsBytes, &orderToUpdate) //unmarshal it aka JSON.parse()
	if err != nil {
		return shim.Error(err.Error())
	}
	orderToUpdate.Luminosity = newStatus //change the owner

	orderJSONasBytes, _ := json.Marshal(orderToUpdate)
	err = stub.PutState(OrderId, orderJSONasBytes) //rewrite the marble
	if err != nil {
		return shim.Error(err.Error())
	}

	fmt.Println("- end updateLuminosity (success)")
	return shim.Success(nil)
}

func (t *OrderChaincode) updateCurrentLocation(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) < 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	OrderId := args[0]
	newLocation := args[1]
	fmt.Println("- update current location ", OrderId, newLocation)

	orderAsBytes, err := stub.GetState(OrderId)
	if err != nil {
		return shim.Error("Failed to get order:" + err.Error())
	} else if orderAsBytes == nil {
		return shim.Error("order does not exist")
	}

	orderToUpdate := order{}
	err = json.Unmarshal(orderAsBytes, &orderToUpdate) //unmarshal it aka JSON.parse()
	if err != nil {
		return shim.Error(err.Error())
	}
	orderToUpdate.CurrentLocation = newLocation //change the temperature

	orderJSONasBytes, _ := json.Marshal(orderToUpdate)
	err = stub.PutState(OrderId, orderJSONasBytes) //rewrite the marble
	if err != nil {
		return shim.Error(err.Error())
	}

	fmt.Println("- end updateCurrentlocation (success)")
	return shim.Success(nil)
}

func (t *OrderChaincode) updateDestinationCity(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) < 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	OrderId := args[0]
	newDestinationCity := args[1]
	fmt.Println("- update destination city ", OrderId, newDestinationCity)

	orderAsBytes, err := stub.GetState(OrderId)
	if err != nil {
		return shim.Error("Failed to get order:" + err.Error())
	} else if orderAsBytes == nil {
		return shim.Error("order does not exist")
	}

	orderToUpdate := order{}
	err = json.Unmarshal(orderAsBytes, &orderToUpdate) //unmarshal it aka JSON.parse()
	if err != nil {
		return shim.Error(err.Error())
	}
	orderToUpdate.DestinationCity = newDestinationCity //change the temperature

	orderJSONasBytes, _ := json.Marshal(orderToUpdate)
	err = stub.PutState(OrderId, orderJSONasBytes) //rewrite the marble
	if err != nil {
		return shim.Error(err.Error())
	}

	fmt.Println("- end updateDestinationCity (success)")
	return shim.Success(nil)
}

func (t *OrderChaincode) updateOriginCity(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) < 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	OrderId := args[0]
	newOriginCity := args[1]
	fmt.Println("- update temperature ", OrderId, newOriginCity)

	orderAsBytes, err := stub.GetState(OrderId)
	if err != nil {
		return shim.Error("Failed to get order:" + err.Error())
	} else if orderAsBytes == nil {
		return shim.Error("order does not exist")
	}

	orderToUpdate := order{}
	err = json.Unmarshal(orderAsBytes, &orderToUpdate) //unmarshal it aka JSON.parse()
	if err != nil {
		return shim.Error(err.Error())
	}
	orderToUpdate.OriginCity = newOriginCity //change the temperature

	orderJSONasBytes, _ := json.Marshal(orderToUpdate)
	err = stub.PutState(OrderId, orderJSONasBytes) //rewrite the marble
	if err != nil {
		return shim.Error(err.Error())
	}

	fmt.Println("- end updateOriginCity (success)")
	return shim.Success(nil)
}

func (t *OrderChaincode) queryHistory(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) < 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	OrderId := args[0]

	fmt.Printf("- start getHistoryForOrder: %s\n", OrderId)

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
