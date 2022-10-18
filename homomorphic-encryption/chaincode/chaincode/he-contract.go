/*
 * SPDX-License-Identifier: Apache-2.0
 */

package chaincode

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/tuneinsight/lattigo/v3/bfv"
	"github.com/tuneinsight/lattigo/v3/rlwe"
)

// HEContract contract for managing HE operations.
type HEContract struct {
	contractapi.Contract
	// encryption parameters
	param bfv.Parameters
}

// Initialize set a default set of encryption parameters for the contract.
// PN13QP218 is a set of default parameters (128 bit security)
func (c *HEContract) Initialize(ctx contractapi.TransactionContextInterface) error {
	var err error
	// BFV parameters (128 bit security) with plaintext modulus 65929217
	paramDef := bfv.PN13QP218
	c.param, err = bfv.NewParametersFromLiteral(paramDef)
	if err != nil {
		return fmt.Errorf("initialization contract failed: %s", err.Error())
	}
	return nil
}

// Create create an instance of the form that can be used to record the submitted data.
// param {string} name The name of the form
func (c *HEContract) Create(ctx contractapi.TransactionContextInterface, name string) error {
	form := &Form{
		Name: name,
		Rows: make([]Row, 0),
	}
	form_encode, _ := json.Marshal(form)
	return ctx.GetStub().PutState(name, form_encode)
}

// Submit submit a row of data to a form.
// param {string} name The name of the form
// param {string} item The item of the row
// param {bytes} revenue The revenue of the row
func (c *HEContract) Submit(ctx contractapi.TransactionContextInterface, name, item string, revenue []byte) error {
	var err error
	form_encode, err := ctx.GetStub().GetState(name)
	if err != nil {
		return fmt.Errorf("get form [%s] failed: %s", name, err.Error())
	}
	form := &Form{}
	if form_encode != nil {
		json.Unmarshal(form_encode, &form)
	} else {
		form.Name = name
		form.Rows = make([]Row, 0)
	}
	// add a row to the form
	form.Rows = append(form.Rows, Row{Item: item, Revenue: revenue})
	form_encode, _ = json.Marshal(form)
	return ctx.GetStub().PutState(name, form_encode)
}

// Query perform calculations on the form and return the results.
// The returned result is a ciphertext,which needs to be decrypted using a key.
// Return an error if an unintended condition occurs.
// param {string} name The name of the form
func (c *HEContract) Query(ctx contractapi.TransactionContextInterface, name string) ([]byte, error) {
	var err error
	form_encode, err := ctx.GetStub().GetState(name)
	if err != nil {
		return nil, fmt.Errorf("get form [%s] failed: %s", name, err.Error())
	}
	form := &Form{}
	if form_encode != nil {
		json.Unmarshal(form_encode, &form)
	} else {
		return nil, fmt.Errorf("form [%s] not exist", name)
	}

	evaluator := bfv.NewEvaluator(c.param, rlwe.EvaluationKey{})
	results := bfv.NewCiphertext(c.param, 1)
	for _, row := range form.Rows {
		op := bfv.NewCiphertext(c.param, 1)
		err = op.UnmarshalBinary(row.Revenue)
		if err != nil {
			return nil, fmt.Errorf("calculation failed: %s", err.Error())
		}
		// results = results + op
		evaluator.Add(results, op, results)
	}
	return results.MarshalBinary()
}
