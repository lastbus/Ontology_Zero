/*
 * Copyright (C) 2018 Onchain <onchain@onchain.com>
 *
 * This file is part of The ontology_Zero.
 *
 * The ontology_Zero is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Lesser General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * The ontology_Zero is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Lesser General Public License for more details.
 *
 * You should have received a copy of the GNU Lesser General Public License
 * along with The ontology_Zero.  If not, see <http://www.gnu.org/licenses/>.
 */

package httpjsonrpc

import (
	. "github.com/Ontology/common"
	"github.com/Ontology/core/asset"
	. "github.com/Ontology/core/contract"
	. "github.com/Ontology/core/transaction"
	"github.com/Ontology/core/transaction/payload"
	"bytes"
)

type PayloadInfo interface{}

//implement PayloadInfo define BookKeepingInfo
type BookKeepingInfo struct {
	Nonce  uint64
	Issuer IssuerInfo
}

//implement PayloadInfo define DeployCodeInfo
type FunctionCodeInfo struct {
	Code           string
	ParameterTypes string
	ReturnType     uint8
}
type InvokeCodeInfo struct {
	CodeHash string
	Code     string
}
type DeployCodeInfo struct {
	Code        *FunctionCodeInfo
	Name        string
	CodeVersion string
	Author      string
	Email       string
	Description string
}

//implement PayloadInfo define IssueAssetInfo
type IssueAssetInfo struct {
}

type IssuerInfo struct {
	X, Y string
}

//implement PayloadInfo define RegisterAssetInfo
type RegisterAssetInfo struct {
	Asset      *asset.Asset
	Amount     Fixed64
	Issuer     IssuerInfo
	Controller string
}

//implement PayloadInfo define TransferAssetInfo
type TransferAssetInfo struct {
}

type RecordInfo struct {
	RecordType string
	RecordData string
}

type BookkeeperInfo struct {
	PubKey     string
	Action     string
	Issuer     IssuerInfo
	Controller string
}

type DataFileInfo struct {
	IPFSPath string
	Filename string
	Note     string
	Issuer   IssuerInfo
}

type PrivacyPayloadInfo struct {
	PayloadType uint8
	Payload     string
	EncryptType uint8
	EncryptAttr string
}

func TransPayloadToHex(p Payload) PayloadInfo {
	switch object := p.(type) {
	case *payload.BookKeeping:
		obj := new(BookKeepingInfo)
		obj.Nonce = object.Nonce
		return obj
	case *payload.BookKeeper:
		obj := new(BookkeeperInfo)
		encodedPubKey, _ := object.PubKey.EncodePoint(true)
		obj.PubKey = ToHexString(encodedPubKey)
		if object.Action == payload.BookKeeperAction_ADD {
			obj.Action = "add"
		} else if object.Action == payload.BookKeeperAction_SUB {
			obj.Action = "sub"
		} else {
			obj.Action = "nil"
		}
		obj.Issuer.X = object.Issuer.X.String()
		obj.Issuer.Y = object.Issuer.Y.String()

		return obj
	case *payload.IssueAsset:
		obj := new(IssueAssetInfo)
		return obj
	case *payload.TransferAsset:
		obj := new(TransferAssetInfo)
		return obj
	case *payload.InvokeCode:
		obj := new(InvokeCodeInfo)
		obj.CodeHash = ToHexString(object.CodeHash.ToArray())
		obj.Code = ToHexString(object.Code)
		return obj
	case *payload.DeployCode:
		obj := new(DeployCodeInfo)
		obj.Code = new(FunctionCodeInfo)
		obj.Code.Code = ToHexString(object.Code.Code)
		obj.Code.ParameterTypes = ToHexString(ContractParameterTypeToByte(object.Code.ParameterTypes))
		obj.Code.ReturnType = uint8(object.Code.ReturnType)
		obj.Name = object.Name
		obj.CodeVersion = object.CodeVersion
		obj.Author = object.Author
		obj.Email = object.Email
		obj.Description = object.Description
		return obj
	case *payload.RegisterAsset:
		obj := new(RegisterAssetInfo)
		obj.Asset = object.Asset
		obj.Amount = object.Amount
		obj.Issuer.X = object.Issuer.X.String()
		obj.Issuer.Y = object.Issuer.Y.String()
		obj.Controller = ToHexString(object.Controller.ToArray())
		return obj
	case *payload.Record:
		obj := new(RecordInfo)
		obj.RecordType = object.RecordType
		obj.RecordData = ToHexString(object.RecordData)
		return obj
	case *payload.PrivacyPayload:
		obj := new(PrivacyPayloadInfo)
		obj.PayloadType = uint8(object.PayloadType)
		obj.Payload = ToHexString(object.Payload)
		obj.EncryptType = uint8(object.EncryptType)
		bytesBuffer := bytes.NewBuffer([]byte{})
		object.EncryptAttr.Serialize(bytesBuffer)
		obj.EncryptAttr = ToHexString(bytesBuffer.Bytes())
		return obj
	case *payload.DataFile:
		obj := new(DataFileInfo)
		obj.IPFSPath = object.IPFSPath
		obj.Filename = object.Filename
		obj.Note = object.Note
		obj.Issuer.X = object.Issuer.X.String()
		obj.Issuer.Y = object.Issuer.Y.String()
		return obj
	}
	return nil
}
