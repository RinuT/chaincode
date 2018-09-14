//peer chaincode instantiate -n mycc -v 0 -c '{"Args":["Shipment","1","Buyer","Seller","CurrentLocation", "DestinationCity", "OriginCity", "ShipmentCondition", "Temperature", "Humidity", "Luminosity"]}' -C myc
//peer chaincode invoke -n mycc -c '{"Args":["registerShipment","2","abc","xyz","NY","CA","NY","good_condition","100","",""]}' -C myc
//peer chaincode invoke -n mycc -c '{"Args":["registerShipment","4","abc","xyz","kochi","Pune","Bangalore","good_condition","100","","28"]}' -C myc
//peer chaincode invoke -n mycc -c '{"Args":["updateTemparature","2","100"]}' -C myc
//peer chaincode invoke -n mycc -c '{"Args":["updateShipmentStatus","2","tampered"]}' -C myc
//peer chaincode invoke -n mycc -c '{"Args":["updateHumidity","2","25"]}' -C myc
//peer chaincode invoke -n mycc -c '{"Args":["updateLuminosity","2","18"]}' -C myc
//peer chaincode invoke -n mycc -c '{"Args":["updateOriginCity","2","kochi"]}' -C myc
// peer chaincode invoke -n mycc -c '{"Args":["updateCurrentLocation","2","Pune"]}' -C myc
//peer chaincode invoke -n mycc -c '{"Args":["updateDestinationCity","2","Banagalore"]}' -C myc
//peer chaincode query -n mycc -c '{"Args":["queryHistory","2"]}' -C myc
//peer chaincode query -n mycc -c '{"Args":["getShipmentDetails","2"]}' -C myc


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

type ShipmentChaincode struct {
}

type Shipment struct {
	ObjectType      string `json:"docType"` 
	ShipmentId      string `json:"ShipmentId"` 
	Buyer           string `json:"Buyer"`
	Seller          string `json:"Seller"`
	CurrentLocation string `json:"CurrentLocation"`
	DestinationCity string `json:"DestinationCity"` 
	OriginCity      string `json:"OriginCity"`
	ShipmentCondition  string `json:"ShipmentCondition"`
	Temperature     string `json:"Temperature"`
	Humidity        string `json:"Humidity"`
	Luminosity      string `json:"Luminosity"`
}

func main() {
	err := shim.Start(new(ShipmentChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}

//Initialization of the Shipment
func (t *ShipmentChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	function, args := stub.GetFunctionAndParameters()
	fmt.Println("init is running " + function)

	ShipmentId := args[0]
	Buyer := args[1]
	Seller := args[2]
	CurrentLocation := args[3]
	DestinationCity := args[4]
	OriginCity := args[5]
	ShipmentCondition := args[6]
	Temperature := args[7]
	Humidity := args[8]
	Luminosity := args[9]

	// ==== Create Shipment object and marshal to JSON ====
	objectType := "Shipment"
	Shipment := &Shipment{objectType, ShipmentId, Buyer, Seller, CurrentLocation, DestinationCity, OriginCity, ShipmentCondition, Temperature, Humidity, Luminosity}
	ShipmentJSONasBytes, err := json.Marshal(Shipment)
	if err != nil {
		return shim.Error(err.Error())
	}

	err = stub.PutState(ShipmentId, ShipmentJSONasBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	fmt.Println("- end init of shipment")

	return shim.Success(nil)
}

//invoke function

func (t *ShipmentChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	function, args := stub.GetFunctionAndParameters()
	fmt.Println("invoke is running " + function)

	// Handle different functions
	if function == "registerShipment" {
		return t.registerShipment(stub, args)
	} else if function == "getShipmentDetails" {
		return t.getShipmentDetails(stub, args)
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
	} else if function == "updateShipmentStatus" {
		return t.updateShipmentStatus(stub, args)
	}

	fmt.Println("invoke did not find func: " + function) //error
	return shim.Error("Received unknown function invocation")
}

//Create an shipment with unique shipment Id
func (t *ShipmentChaincode) registerShipment(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error

	// ==== Input sanitation ====
	fmt.Println("- start registering shipment")

	ShipmentId := args[0]
	Buyer := args[1]
	Seller := args[2]
	CurrentLocation := args[3]
	DestinationCity := args[4]
	OriginCity := args[5]
	ShipmentCondition := args[6]
	Temperature := args[7]
	Humidity := args[8]
	Luminosity := args[9]

	// ==== Check if shipment already exists ====
	ShipmentAsBytes, err := stub.GetState(ShipmentId)
	if err != nil {
		return shim.Error("Failed to register shipment: " + err.Error())
	} else if ShipmentAsBytes != nil {
		fmt.Println("This shipment already exists: " + ShipmentId)
		return shim.Error("This shipment already exists: " + ShipmentId)
	}

	// ==== Create Shipment object and marshal to JSON ====
	if (Humidity == "undefined" || Humidity == "" || Humidity == "null" || Luminosity == "undefined" || Luminosity == "" || Luminosity == "null") {
		objectType := "Shipment"
		Shipment := &Shipment{objectType, ShipmentId, Buyer, Seller, CurrentLocation, DestinationCity, OriginCity, ShipmentCondition, Temperature, Humidity, Luminosity}
		fmt.Println(Shipment)
		ShipmentJSONasBytes, err := json.Marshal(Shipment)
		fmt.Println(ShipmentJSONasBytes)
		if err != nil {
			return shim.Error(err.Error())
		}

		// === Save Shipment to state ===
		err = stub.PutState(ShipmentId, ShipmentJSONasBytes)
		if err != nil {
			return shim.Error(err.Error())
		}
	}
	fmt.Println("- end registering shimpment")
	return shim.Success(nil)

}

//search the Shipment details using Shipment Id
func (t *ShipmentChaincode) getShipmentDetails(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var ShipmentId, jsonResp string
	var err error

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting shipment Id to query")
	}

	ShipmentId = args[0]
	valAsbytes, err := stub.GetState(ShipmentId)
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + ShipmentId + "\"}"
		return shim.Error(jsonResp)
	} else if valAsbytes == nil {
		jsonResp = "{\"Error\":\"shipment does not exist: " + ShipmentId + "\"}"
		return shim.Error(jsonResp)
	}

	return shim.Success(valAsbytes)
}

func (t *ShipmentChaincode) updateTemparature(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) < 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	ShipmentId := args[0]
	newStatus := args[1]
	fmt.Println("- update temperature ", ShipmentId, newStatus)

	ShipmentAsBytes, err := stub.GetState(ShipmentId)
	if err != nil {
		return shim.Error("Failed to get shipment details:" + err.Error())
	} else if ShipmentAsBytes == nil {
		return shim.Error("shipment does not exist")
	}

	ShipmentToUpdate := Shipment{}
	err = json.Unmarshal(ShipmentAsBytes, &ShipmentToUpdate) //unmarshal it aka JSON.parse()
	if err != nil {
		return shim.Error(err.Error())
	}
	ShipmentToUpdate.Temperature = newStatus //change the temperature

	ShipmentJSONasBytes, _ := json.Marshal(ShipmentToUpdate)
	err = stub.PutState(ShipmentId, ShipmentJSONasBytes) //rewrite the marble
	if err != nil {
		return shim.Error(err.Error())
	}

	fmt.Println("- end updateTemperature (success)")
	return shim.Success(nil)
}

func (t *ShipmentChaincode) updateHumidity(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) < 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	ShipmentId := args[0]
	newStatus := args[1]
	fmt.Println("- update humidity ", ShipmentId, newStatus)

	ShipmentAsBytes, err := stub.GetState(ShipmentId)
	if err != nil {
		return shim.Error("Failed to get shipment details:" + err.Error())
	} else if ShipmentAsBytes == nil {
		return shim.Error("shimoent does not exist")
	}

	ShipmentToUpdate := Shipment{}
	err = json.Unmarshal(ShipmentAsBytes, &ShipmentToUpdate) //unmarshal it aka JSON.parse()
	if err != nil {
		return shim.Error(err.Error())
	}
	ShipmentToUpdate.Humidity = newStatus //change the humidity

	ShipmentJSONasBytes, _ := json.Marshal(ShipmentToUpdate)
	err = stub.PutState(ShipmentId, ShipmentJSONasBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	fmt.Println("- end updateHumidity (success)")
	return shim.Success(nil)
}

func (t *ShipmentChaincode) updateLuminosity(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) < 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	ShipmentId := args[0]
	newStatus := args[1]
	fmt.Println("- update Luminosity ", ShipmentId, newStatus)

	ShipmentAsBytes, err := stub.GetState(ShipmentId)
	if err != nil {
		return shim.Error("Failed to get shipment details:" + err.Error())
	} else if ShipmentAsBytes == nil {
		return shim.Error("shipment does not exist")
	}

	ShipmentToUpdate := Shipment{}
	err = json.Unmarshal(ShipmentAsBytes, &ShipmentToUpdate) //unmarshal it aka JSON.parse()
	if err != nil {
		return shim.Error(err.Error())
	}
	ShipmentToUpdate.Luminosity = newStatus //change the Luminosity

	ShipmentJSONasBytes, _ := json.Marshal(ShipmentToUpdate)
	err = stub.PutState(ShipmentId, ShipmentJSONasBytes) //rewrite the shipment
	if err != nil {
		return shim.Error(err.Error())
	}

	fmt.Println("- end updateLuminosity (success)")
	return shim.Success(nil)
}

func (t *ShipmentChaincode) updateCurrentLocation(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) < 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	ShipmentId := args[0]
	newLocation := args[1]
	fmt.Println("- update current location ", ShipmentId, newLocation)

	ShipmentAsBytes, err := stub.GetState(ShipmentId)
	if err != nil {
		return shim.Error("Failed to get shipment details:" + err.Error())
	} else if ShipmentAsBytes == nil {
		return shim.Error("shipment does not exist")
	}

	ShipmentToUpdate := Shipment{}
	err = json.Unmarshal(ShipmentAsBytes, &ShipmentToUpdate) //unmarshal it aka JSON.parse()
	if err != nil {
		return shim.Error(err.Error())
	}
	ShipmentToUpdate.CurrentLocation = newLocation 

	ShipmentJSONasBytes, _ := json.Marshal(ShipmentToUpdate)
	err = stub.PutState(ShipmentId, ShipmentJSONasBytes) 
	if err != nil {
		return shim.Error(err.Error())
	}

	fmt.Println("- end updateCurrentlocation (success)")
	return shim.Success(nil)
}

func (t *ShipmentChaincode) updateDestinationCity(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) < 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	ShipmentId := args[0]
	newDestinationCity := args[1]
	fmt.Println("- update destination city ", ShipmentId, newDestinationCity)

	ShipmentAsBytes, err := stub.GetState(ShipmentId)
	if err != nil {
		return shim.Error("Failed to get shipment details:" + err.Error())
	} else if ShipmentAsBytes == nil {
		return shim.Error("shipment does not exist")
	}

	ShipmentToUpdate := Shipment{}
	err = json.Unmarshal(ShipmentAsBytes, &ShipmentToUpdate) 
	if err != nil {
		return shim.Error(err.Error())
	}
	ShipmentToUpdate.DestinationCity = newDestinationCity 

	ShipmentJSONasBytes, _ := json.Marshal(ShipmentToUpdate)
	err = stub.PutState(ShipmentId, ShipmentJSONasBytes) 
	if err != nil {
		return shim.Error(err.Error())
	}

	fmt.Println("- end updateDestinationCity (success)")
	return shim.Success(nil)
}

func (t *ShipmentChaincode) updateOriginCity(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) < 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	ShipmentId := args[0]
	newOriginCity := args[1]
	fmt.Println("- update temperature ", ShipmentId, newOriginCity)

	ShipmentAsBytes, err := stub.GetState(ShipmentId)
	if err != nil {
		return shim.Error("Failed to get shipment details:" + err.Error())
	} else if ShipmentAsBytes == nil {
		return shim.Error("shipment does not exist")
	}

	ShipmentToUpdate := Shipment{}
	err = json.Unmarshal(ShipmentAsBytes, &ShipmentToUpdate) 
	if err != nil {
		return shim.Error(err.Error())
	}
	ShipmentToUpdate.OriginCity = newOriginCity 

	ShipmentJSONasBytes, _ := json.Marshal(ShipmentToUpdate)
	err = stub.PutState(ShipmentId, ShipmentJSONasBytes) 
	if err != nil {
		return shim.Error(err.Error())
	}

	fmt.Println("- end updateOriginCity (success)")
	return shim.Success(nil)
}

func (t *ShipmentChaincode) updateShipmentStatus(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) < 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	ShipmentId := args[0]
	newStatus := args[1]
	fmt.Println("- update Shipment status ", ShipmentId, newStatus)

	ShipmentAsBytes, err := stub.GetState(ShipmentId)
	if err != nil {
		return shim.Error("Failed to get shipment details:" + err.Error())
	} else if ShipmentAsBytes == nil {
		return shim.Error("shipment does not exist")
	}

	ShipmentToUpdate := Shipment{}
	err = json.Unmarshal(ShipmentAsBytes, &ShipmentToUpdate) 
	if err != nil {
		return shim.Error(err.Error())
	}
	ShipmentToUpdate.ShipmentCondition = newStatus 

	ShipmentJSONasBytes, _ := json.Marshal(ShipmentToUpdate)
	err = stub.PutState(ShipmentId, ShipmentJSONasBytes) 
	if err != nil {
		return shim.Error(err.Error())
	}

	fmt.Println("- end updateShipmentStatus (success)")
	return shim.Success(nil)
}

func (t *ShipmentChaincode) queryHistory(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) < 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	ShipmentId := args[0]

	fmt.Printf("- start getHistoryForShipment: %s\n", ShipmentId)

	resultsIterator, err := stub.GetHistoryForKey(ShipmentId)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing historic values for the shipment
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

	fmt.Printf("- returning history of the shipment:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}