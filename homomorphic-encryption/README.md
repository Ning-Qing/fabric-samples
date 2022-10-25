# HE scenario

The HE can provide certain types of operations on encrypted data without decrypting the encrypted data.

## Case

For some reason the data needs to be encrypted and computed on the ciphertext.  
For example, the following table:

|form|items|revenue|
|--- | --- | --- |
||item1| 2000|
||item2|-1000|
||item3|2000|
||results|3000|

## Contract

### Initialize

**Initialize()** set a default set of encryption parameters for the contract

### Create

**Create(args string)** create a form named args that can be used to save data

### Submit

**Submit(args1 string, args2 string, args3 []byte)** can record a row of data on a form named args1

|args1|items|revenue|
|--- | --- | --- |
||args2| args3|

### Query

**Query(args string)** query the result of the calculation of the data on the form named args, the ciphertext calculation is done here

## Use

can use the sdk to make contract calls or use the key tool to process the data and then make the calls.

refer to [example](./chaincode/chaincode/he-contract_test.go)
