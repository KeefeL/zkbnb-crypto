/*
 * Copyright © 2021 Zecrey Protocol
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

package legendTxTypes

import (
	"bytes"
	"encoding/json"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr/mimc"
	"github.com/consensys/gnark-crypto/ecc/bn254/twistededwards/eddsa"
	"hash"
	"log"
	"math/big"
)

type AtomicMatchSegmentFormat struct {
	AccountIndex      int64
	BuyOffer          string
	SellOffer         string
	GasAccountIndex   int64
	GasFeeAssetId     int64
	GasFeeAssetAmount string
	Nonce             int64
}

/*
	ConstructMintNftTxInfo: construct mint nft tx, sign txInfo
*/
func ConstructAtomicMatchTxInfo(sk *PrivateKey, segmentStr string) (txInfo *AtomicMatchTxInfo, err error) {
	var segmentFormat *AtomicMatchSegmentFormat
	err = json.Unmarshal([]byte(segmentStr), &segmentFormat)
	if err != nil {
		log.Println("[ConstructMintNftTxInfo] err info:", err)
		return nil, err
	}
	gasFeeAmount, err := StringToBigInt(segmentFormat.GasFeeAssetAmount)
	if err != nil {
		log.Println("[ConstructBuyNftTxInfo] unable to convert string to big int:", err)
		return nil, err
	}
	var (
		buyOffer, sellOffer *OfferTxInfo
	)
	err = json.Unmarshal([]byte(segmentFormat.BuyOffer), buyOffer)
	if err != nil {
		log.Println("[ConstructBuyNftTxInfo] unable to unmarshal offer", err.Error())
		return nil, err
	}
	err = json.Unmarshal([]byte(segmentFormat.SellOffer), sellOffer)
	if err != nil {
		log.Println("[ConstructBuyNftTxInfo] unable to unmarshal offer", err.Error())
		return nil, err
	}
	txInfo = &AtomicMatchTxInfo{
		AccountIndex:      segmentFormat.AccountIndex,
		BuyOffer:          buyOffer,
		SellOffer:         sellOffer,
		GasAccountIndex:   segmentFormat.GasAccountIndex,
		GasFeeAssetId:     segmentFormat.GasFeeAssetId,
		GasFeeAssetAmount: gasFeeAmount,
		Nonce:             segmentFormat.Nonce,
		Sig:               nil,
	}
	// compute call data hash
	hFunc := mimc.NewMiMC()
	// compute msg hash
	msgHash, err := ComputeAtomicMatchMsgHash(txInfo, hFunc)
	if err != nil {
		log.Println("[ConstructMintNftTxInfo] unable to compute hash: ", err.Error())
		return nil, err
	}
	// compute signature
	hFunc.Reset()
	sigBytes, err := sk.Sign(msgHash, hFunc)
	if err != nil {
		log.Println("[ConstructMintNftTxInfo] unable to sign:", err)
		return nil, err
	}
	txInfo.Sig = sigBytes
	return txInfo, nil
}

type AtomicMatchTxInfo struct {
	AccountIndex      int64
	BuyOffer          *OfferTxInfo
	SellOffer         *OfferTxInfo
	GasAccountIndex   int64
	GasFeeAssetId     int64
	GasFeeAssetAmount *big.Int
	Nonce             int64
	Sig               []byte
}

func ComputeAtomicMatchMsgHash(txInfo *AtomicMatchTxInfo, hFunc hash.Hash) (msgHash []byte, err error) {
	hFunc.Reset()
	var buf bytes.Buffer
	packedBuyAmount, err := ToPackedAmount(txInfo.BuyOffer.AssetAmount)
	if err != nil {
		log.Println("[ComputeTransferMsgHash] unable to packed amount:", err.Error())
		return nil, err
	}
	packedSellAmount, err := ToPackedAmount(txInfo.SellOffer.AssetAmount)
	if err != nil {
		log.Println("[ComputeTransferMsgHash] unable to packed amount:", err.Error())
		return nil, err
	}
	packedFee, err := ToPackedFee(txInfo.GasFeeAssetAmount)
	if err != nil {
		log.Println("[ComputeTransferMsgHash] unable to packed amount:", err.Error())
		return nil, err
	}
	WriteInt64IntoBuf(&buf, txInfo.AccountIndex)
	WriteInt64IntoBuf(&buf, txInfo.BuyOffer.Type)
	WriteInt64IntoBuf(&buf, txInfo.BuyOffer.OfferId)
	WriteInt64IntoBuf(&buf, txInfo.BuyOffer.AccountIndex)
	WriteInt64IntoBuf(&buf, txInfo.BuyOffer.NftIndex)
	WriteInt64IntoBuf(&buf, txInfo.BuyOffer.AssetId)
	WriteInt64IntoBuf(&buf, packedBuyAmount)
	WriteInt64IntoBuf(&buf, txInfo.BuyOffer.ListedAt)
	WriteInt64IntoBuf(&buf, txInfo.BuyOffer.ExpiredAt)
	var (
		buyerSig, sellerSig = new(eddsa.Signature), new(eddsa.Signature)
	)
	_, err = buyerSig.SetBytes(txInfo.BuyOffer.Sig)
	if err != nil {
		log.Println("[ComputeAtomicMatchMsgHash] unable to convert to sig: ", err.Error())
		return nil, err
	}
	buf.Write(buyerSig.R.X.Marshal())
	buf.Write(buyerSig.R.Y.Marshal())
	buf.Write(buyerSig.S[:])
	WriteInt64IntoBuf(&buf, txInfo.SellOffer.Type)
	WriteInt64IntoBuf(&buf, txInfo.SellOffer.OfferId)
	WriteInt64IntoBuf(&buf, txInfo.SellOffer.AccountIndex)
	WriteInt64IntoBuf(&buf, txInfo.SellOffer.NftIndex)
	WriteInt64IntoBuf(&buf, txInfo.SellOffer.AssetId)
	WriteInt64IntoBuf(&buf, packedSellAmount)
	WriteInt64IntoBuf(&buf, txInfo.SellOffer.ListedAt)
	WriteInt64IntoBuf(&buf, txInfo.SellOffer.ExpiredAt)
	_, err = sellerSig.SetBytes(txInfo.SellOffer.Sig)
	if err != nil {
		log.Println("[ComputeAtomicMatchMsgHash] unable to convert to sig: ", err.Error())
		return nil, err
	}
	buf.Write(sellerSig.R.X.Marshal())
	buf.Write(sellerSig.R.Y.Marshal())
	buf.Write(sellerSig.S[:])
	WriteInt64IntoBuf(&buf, txInfo.GasAccountIndex)
	WriteInt64IntoBuf(&buf, txInfo.GasFeeAssetId)
	WriteInt64IntoBuf(&buf, packedFee)
	WriteInt64IntoBuf(&buf, txInfo.Nonce)
	hFunc.Write(buf.Bytes())
	msgHash = hFunc.Sum(nil)
	return msgHash, nil
}