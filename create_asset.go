
package main
import "fmt"
import "github.com/hyperledger/fabric/core/chaincode/shim"
import "github.com/hyperledger/fabric/protos/peer"
import "encoding/json"  


type NewAsset struct{
	uuid string  'json:"uuid"'  
	material string  'json:"material"'  
	make string  'json:"make"'  
	material_loaction string  'json:"material_location"'  
	shipment_status string  'json:"shipment_status"'  
	product_status string  'json:"product_status"'  
	
}

// Define Status codes for the response
const (
	OK    = 200
	ERROR = 500
)

// func(n NewAsset) details(){
//     fmt.Println(n," "+"NewAsset")
// }

// func main(){
// 	p1 := NewAsset{"1","scanner","siemens","pune","in_good_condition","in_production"}
// 	p1.details()
// }
func (t *NewAsset) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte,
	error) {
	fmt.println("inside init function")
	var uuid = args[0]
	var NewAssetInput = args[1]

	err := stub.PutState(uuid, []byte(NewAssetInput))
	if err != nil {
		fmt.println("could not save the product", err)
		return nil, err
}
fmt.println("product created successfully")
return shim.Success(nil)
}

   func (t *NewAsset) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte,
	   error) {
		fn, args := stub.GetFunctionAndParameters()

		var result string
		var err error
		if fn == "set" {
				result, err = set(stub, args)
		} else { // assume 'get' even if fn is nil
				result, err = get(stub, args)
		}
		if err != nil {
				return shim.Error(err.Error())
		}
	
		// Return the result as success payload
		return shim.Success([]byte(result))
	  }
   
	  func set(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
		  fmt.println("inside create product")
		if len(args) < 2 {
			fmt.println("invalid number of arguments")
			return nil, errors.New("expected atleast two arguments for product creation")
		}
		var uuid = args[0]
		var NewAssetInput = args[1]
	
		err := stub.PutState(uuid, []byte(NewAssetInput))
		if err != nil {
			fmt.println("could not save the product", err)
			return nil, err
		}
		fmt.println("Successfully created and saved product into ledger")
		return nil, nil
	}

	func get(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
		fmt.println("inside get product details")
		if len(args) < 1 {
			fmt.println("invalid number of arguments")
			return nil, errors.New("missing product id")
		}
		var uuid = args[0]
		bytes, err := stub.GetState(uuid)
		if err != nil {
				fmt.println("could not fetch the product with id "+uuid+" from ledger", err)
		}
		return nil,err
	}
	return bytes, nil
	}
//    func (t *NewAsset) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte,
// 	error) {
// 	return nil, nil
// 	}
   

   func main() {
	   err := shim.Start(new(NewAsset))
	   if err != nil {
	   fmt.Println("Could not start SampleChaincode")
	   } else {
	   fmt.Println("SampleChaincode successfully started")
	   }
	  }