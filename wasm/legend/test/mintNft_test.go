package test

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/bnb-chain/zkbas-crypto/wasm/legend/legendTxTypes"
	"log"
	"testing"
)

func TestMintNftSegmentFormat(t *testing.T){

	var segmentFormat *legendTxTypes.MintNftSegmentFormat
	segmentFormat = &legendTxTypes.MintNftSegmentFormat{
		CreatorAccountIndex: 15,
		ToAccountIndex:      1,
		ToAccountNameHash:   "ddc6171f9fe33153d95c8394c9135c277eb645401b85eb499393a2aefe6422a6",
		NftContentHash:      "7eb645401b85eb499393a2aefe6422a6ddc6171f9fe33153d95c8394c9135c27",
		NftCollectionId:     65,
		CreatorTreasuryRate: 30,
		GasAccountIndex:     1,
		GasFeeAssetId:       3,
		GasFeeAssetAmount:   "3",
		ExpiredAt:           1654656781000, // milli seconds
		Nonce:               1,
	}

	res, err := json.Marshal(segmentFormat)
	assert.Nil(t, err)

	log.Println(string(res))
}
