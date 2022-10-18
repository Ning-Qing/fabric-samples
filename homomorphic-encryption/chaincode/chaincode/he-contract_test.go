package chaincode

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/tuneinsight/lattigo/v3/bfv"
	"github.com/tuneinsight/lattigo/v3/rlwe"
)

var paramDef = bfv.PN13QP218

type MockStub struct {
	shim.ChaincodeStubInterface
	mock.Mock
}

func (ms *MockStub) GetState(key string) ([]byte, error) {
	args := ms.Called(key)
	return args.Get(0).([]byte), args.Error(1)
}

func (ms *MockStub) PutState(key string, value []byte) error {
	args := ms.Called(key, value)
	return args.Error(0)
}

type MockContext struct {
	contractapi.TransactionContextInterface
	mock.Mock
}

func (mc *MockContext) GetStub() shim.ChaincodeStubInterface {
	args := mc.Called()
	return args.Get(0).(*MockStub)
}

func genKeyPair(t *testing.T) (*rlwe.SecretKey, *rlwe.PublicKey) {
	param, err := bfv.NewParametersFromLiteral(paramDef)
	if err != nil {
		fmt.Println("Generate key pair failed", err.Error())
		t.FailNow()
	}

	kgen := bfv.NewKeyGenerator(param)
	return kgen.GenKeyPair()
}

func decrypto(t *testing.T, data []byte, sk *rlwe.SecretKey) int64 {
	param, err := bfv.NewParametersFromLiteral(paramDef)
	if err != nil {
		fmt.Println("decrypt: set of BFV parameters failed", err.Error())
		t.FailNow()
	}
	ciphertext := bfv.NewCiphertext(param, 1)
	err = ciphertext.UnmarshalBinary(data)
	if err != nil {
		fmt.Println("failed to unmarshal value")
		t.FailNow()
	}
	decryptor := bfv.NewDecryptor(param, sk)
	text := decryptor.DecryptNew(ciphertext)
	encoder := bfv.NewEncoder(param)
	return encoder.DecodeIntNew(text)[0]
}

func encrypto(t *testing.T, data int64, pubkey *rlwe.PublicKey) []byte {
	param, err := bfv.NewParametersFromLiteral(paramDef)
	if err != nil {
		fmt.Println("encrypt: set of BFV parameters failed", err.Error())
		t.FailNow()
	}

	text := bfv.NewPlaintext(param)
	encoder := bfv.NewEncoder(param)
	encoder.Encode([]int64{data}, text)
	encryptor := bfv.NewEncryptor(param, pubkey)
	ciphertext := encryptor.EncryptNew(text)
	raw, err := ciphertext.MarshalBinary()
	if err != nil {
		fmt.Println("marshalBinary encodes a Ciphertext failed", err.Error())
		t.FailNow()
	}
	return raw
}

func TestInitialize(t *testing.T) {
	mc := new(MockContext)
	c := new(HEContract)

	option := c.Initialize(mc)
	assert.Equal(t, nil, option)
}

func TestSubmit(t *testing.T) {
	_, pk := genKeyPair(t)
	revenue := encrypto(t, 2000, pk)
	mc := new(MockContext)
	ms := new(MockStub)
	mc.On("GetStub").Return(ms)
	form := &Form{
		Name: "form1",
		Rows: make([]Row, 0),
	}
	form.Rows = append(form.Rows, Row{Item: "item1", Revenue: revenue})
	form1_encode, _ := json.Marshal(form)
	ms.On("GetState", "form1").Return([]byte(nil), nil)
	ms.On("PutState", "form1", form1_encode).Return(nil)

	c := new(HEContract)
	c.param, _ = bfv.NewParametersFromLiteral(paramDef)
	option := c.Submit(mc, "form1", "item1", revenue)
	assert.Equal(t, nil, option)
}

func TestQuery(t *testing.T) {
	sk, pk := genKeyPair(t)
	mc := new(MockContext)
	ms := new(MockStub)
	mc.On("GetStub").Return(ms)
	form := &Form{
		Name: "form1",
		Rows: make([]Row, 0),
	}
	form.Rows = append(form.Rows, Row{Item: "item1", Revenue: encrypto(t, 2000, pk)})
	form.Rows = append(form.Rows, Row{Item: "item2", Revenue: encrypto(t, -1000, pk)})
	form.Rows = append(form.Rows, Row{Item: "item3", Revenue: encrypto(t, 2000, pk)})
	form_encode, _ := json.Marshal(form)
	ms.On("GetState", "form1").Return(form_encode, nil)

	c := new(HEContract)
	c.param, _ = bfv.NewParametersFromLiteral(paramDef)
	payload, err := c.Query(mc, "form1")
	if err != nil {
		fmt.Println("query failed", err.Error())
		t.FailNow()
	}
	option := decrypto(t, payload, sk)
	assert.Equal(t, int64(3000), option)
}
