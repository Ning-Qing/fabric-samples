/*
 * SPDX-License-Identifier: Apache-2.0
 */
 
package chaincode

type Form struct {
	// form name
	Name string `json:"name"`
	// form record items
	Rows []Row `json:"rows"`
}

type Row struct {
	// form item
	Item string `json:"iten"`
	// the revenue of this form item
	Revenue []byte `json:"revenue"`
}
