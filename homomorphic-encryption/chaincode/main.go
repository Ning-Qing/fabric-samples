/*
 * SPDX-License-Identifier: Apache-2.0
 */

package main

import (
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/hyperledger/fabric-contract-api-go/metadata"
	"github.com/hyperledger/fabric-samples/homomorphic-encryption/chaincode-go/chaincode"
)

func main() {
	HE_Contract := new(chaincode.HEContract)
	HE_Contract.Info.Version = "0.0.1"
	HE_Contract.Info.Description = "HE on fabric chaincode"
	HE_Contract.Info.License = new(metadata.LicenseMetadata)
	HE_Contract.Info.License.Name = "Apache-2.0"
	HE_Contract.Info.Contact = new(metadata.ContactMetadata)
	HE_Contract.Info.Contact.Name = "Ning-Qing"

	chaincode, err := contractapi.NewChaincode(HE_Contract)
	chaincode.Info.Title = "HE chaincode"
	chaincode.Info.Version = "0.0.1"
	if err != nil {
		panic("Could not create chaincode from HEContract." + err.Error())
	}
	err = chaincode.Start()

	if err != nil {
		panic("Failed to start chaincode. " + err.Error())
	}
}
