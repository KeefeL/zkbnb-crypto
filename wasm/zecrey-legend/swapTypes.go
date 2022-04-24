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

package zecrey_legend

import (
	"bytes"
	"encoding/json"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr/mimc"
	"hash"
	"log"
	"math/big"
)

type SwapSegmentFormat struct {
	FromAccountIndex       int64  `json:"from_account_index"`
	ToAccountIndex         int64  `json:"to_account_index"`
	PairIndex              int64  `json:"pair_index"`
	AssetAId               int64  `json:"asset_a_id"`
	AssetAAmount           string `json:"asset_a_amount"`
	AssetBId               int64  `json:"asset_b_id"`
	AssetBMinAmount        string `json:"asset_b_min_amount"`
	AssetBAmountDelta      string `json:"asset_b_amount_delta"`
	PoolAAmount            string `json:"pool_a_amount"`
	PoolBAmount            string `json:"pool_b_amount"`
	FeeRate                int64  `json:"fee_rate"`
	TreasuryAccountIndex   int64  `json:"treasury_account_index"`
	TreasuryRate           int64  `json:"treasury_rate"`
	TreasuryFeeAmountDelta int64  `json:"treasury_fee_amount_delta"`
	GasAccountIndex        int64  `json:"gas_account_index"`
	GasFeeAssetId          int64  `json:"gas_fee_asset_id"`
	GasFeeAssetAmount      int64  `json:"gas_fee_asset_amount"`
	Nonce                  int64  `json:"nonce"`
}

func ConstructSwapTxInfo(sk *PrivateKey, segmentStr string) (txInfo *SwapTxInfo, err error) {
	var segmentFormat *SwapSegmentFormat
	err = json.Unmarshal([]byte(segmentStr), &segmentFormat)
	if err != nil {
		log.Println("[ConstructSwapTxInfo] err info:", err)
		return nil, err
	}
	assetAAmount, err := StringToBigInt(segmentFormat.AssetAAmount)
	if err != nil {
		log.Println("[ConstructBuyNftTxInfo] unable to convert string to big int:", err)
		return nil, err
	}
	assetBMinAmount, err := StringToBigInt(segmentFormat.AssetBMinAmount)
	if err != nil {
		log.Println("[ConstructBuyNftTxInfo] unable to convert string to big int:", err)
		return nil, err
	}
	assetBAmountDelta, err := StringToBigInt(segmentFormat.AssetBAmountDelta)
	if err != nil {
		log.Println("[ConstructBuyNftTxInfo] unable to convert string to big int:", err)
		return nil, err
	}
	txInfo = &SwapTxInfo{
		FromAccountIndex:       segmentFormat.FromAccountIndex,
		ToAccountIndex:         segmentFormat.ToAccountIndex,
		PairIndex:              segmentFormat.PairIndex,
		AssetAId:               segmentFormat.AssetAId,
		AssetAAmount:           assetAAmount,
		AssetBId:               segmentFormat.AssetBId,
		AssetBMinAmount:        assetBMinAmount,
		AssetBAmountDelta:      assetBAmountDelta,
		PoolAAmount:            ZeroBigInt,
		PoolBAmount:            ZeroBigInt,
		FeeRate:                segmentFormat.FeeRate,
		TreasuryAccountIndex:   segmentFormat.TreasuryAccountIndex,
		TreasuryRate:           segmentFormat.TreasuryRate,
		TreasuryFeeAmountDelta: segmentFormat.TreasuryFeeAmountDelta,
		GasAccountIndex:        segmentFormat.GasAccountIndex,
		GasFeeAssetId:          segmentFormat.GasFeeAssetId,
		GasFeeAssetAmount:      segmentFormat.GasFeeAssetAmount,
		Nonce:                  segmentFormat.Nonce,
		Sig:                    nil,
	}
	hFunc := mimc.NewMiMC()
	msgHash := ComputeSwapMsgHash(txInfo, hFunc)
	hFunc.Reset()
	sigBytes, err := sk.Sign(msgHash, hFunc)
	if err != nil {
		log.Println("[ConstructSwapTxInfo] unable to sign:", err)
		return nil, err
	}
	txInfo.Sig = sigBytes
	return txInfo, nil
}

type SwapTxInfo struct {
	FromAccountIndex       int64
	ToAccountIndex         int64
	PairIndex              int64
	AssetAId               int64
	AssetAAmount           *big.Int
	AssetBId               int64
	AssetBMinAmount        *big.Int
	AssetBAmountDelta      *big.Int
	PoolAAmount            *big.Int
	PoolBAmount            *big.Int
	FeeRate                int64
	TreasuryAccountIndex   int64
	TreasuryRate           int64
	TreasuryFeeAmountDelta int64
	GasAccountIndex        int64
	GasFeeAssetId          int64
	GasFeeAssetAmount      int64
	Nonce                  int64
	Sig                    []byte
}

func ComputeSwapMsgHash(txInfo *SwapTxInfo, hFunc hash.Hash) (msgHash []byte) {
	/*
		hFunc.Write(
			tx.FromAccountIndex,
			tx.ToAccountIndex,
			tx.PairIndex,
			tx.AssetAId,
			tx.AssetAAmount,
			tx.AssetBId,
			tx.AssetBMinAmount,
			tx.FeeRate,
			tx.TreasuryAccountIndex,
			tx.TreasuryRate,
			tx.GasAccountIndex,
			tx.GasFeeAssetId,
			tx.GasFeeAssetAmount,
		)
		hFunc.Write(nonce)
	*/
	hFunc.Reset()
	var buf bytes.Buffer
	writeInt64IntoBuf(&buf, txInfo.FromAccountIndex)
	writeInt64IntoBuf(&buf, txInfo.ToAccountIndex)
	writeInt64IntoBuf(&buf, txInfo.PairIndex)
	writeInt64IntoBuf(&buf, txInfo.AssetAId)
	writeBigIntIntoBuf(&buf, txInfo.AssetAAmount)
	writeInt64IntoBuf(&buf, txInfo.AssetBId)
	writeBigIntIntoBuf(&buf, txInfo.AssetBMinAmount)
	writeInt64IntoBuf(&buf, txInfo.FeeRate)
	writeInt64IntoBuf(&buf, txInfo.TreasuryAccountIndex)
	writeInt64IntoBuf(&buf, txInfo.TreasuryRate)
	writeInt64IntoBuf(&buf, txInfo.GasAccountIndex)
	writeInt64IntoBuf(&buf, txInfo.GasFeeAssetId)
	writeInt64IntoBuf(&buf, txInfo.GasFeeAssetAmount)
	writeInt64IntoBuf(&buf, txInfo.Nonce)
	hFunc.Write(buf.Bytes())
	msgHash = hFunc.Sum(nil)
	return msgHash
}