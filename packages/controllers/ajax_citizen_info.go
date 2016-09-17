// Copyright 2016 The go-daylight Authors
// This file is part of the go-daylight library.
//
// The go-daylight library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-daylight library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-daylight library. If not, see <http://www.gnu.org/licenses/>.

package controllers

import (
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/DayLightProject/go-daylight/packages/utils"
)

const ACitizenInfo = `ajax_citizen_info`

type FieldInfo struct {
	Name     string `json:"name"`
	HtmlType string `json:"htmlType"`
	TxType   string `json:"txType"`
	Title    string `json:"title"`
}

type CitizenInfoJson struct {
	Result bool   `json:"result"`
	Error  string `json:"error"`
}

func init() {
	newPage(ACitizenInfo, `json`)
}

func (c *Controller) AjaxCitizenInfo() interface{} {
	var (
		result CitizenInfoJson
		err    error
		data   map[string]string
	)
	c.w.Header().Add("Access-Control-Allow-Origin", "*")
	stateCode := utils.StrToInt64(c.r.FormValue(`stateId`))
	statePrefix, err := c.GetStatePrefix(stateCode)

	field, err := c.Single(`SELECT value FROM ` + statePrefix + `_state_settings where parameter='citizen_fields'`).String()
	vals := make(map[string]string)
	time := c.r.FormValue(`time`)
	walletId := c.r.FormValue(`walletId`)

	if err == nil {
		var (
			fields    []FieldInfo
			sign      []byte
			checkSign bool
		)
		if err = json.Unmarshal([]byte(field), &fields); err == nil {
			for _, ifield := range fields {
				vals[ifield.Name] = c.r.FormValue(ifield.Name)
			}

			data, err = c.OneRow("SELECT public_key_0, public_key_1, public_key_2 FROM dlt_wallets WHERE wallet_id = ?", walletId).String()
			if err == nil {
				var PublicKeys [][]byte
				PublicKeys = append(PublicKeys, []byte(data["public_key_0"]))
				forSign := fmt.Sprintf("CitizenInfo,%s,%s", time, walletId)
				sign, err = hex.DecodeString(c.r.FormValue(`signature1`))

				if err == nil {
					checkSign, err = utils.CheckSign(PublicKeys, forSign, sign, true)
					if err == nil && !checkSign {
						err = fmt.Errorf(`incorrect signature`)
					}
				}
			}
		}
	}
	if err == nil {
		data, err = c.OneRow(`SELECT * FROM `+statePrefix+`_citizenship_requests WHERE dlt_wallet_id = ? order by request_id desc`, walletId).String()
		if err != nil || data == nil || len(data) == 0 {
			err = fmt.Errorf(`unknown request for wallet %s`, walletId)
		} else {
			var (
				fval []byte
			)
			if fval, err = json.Marshal(vals); err == nil {
				err = c.ExecSql(`INSERT INTO `+statePrefix+`_citizens_requests_private ( request_id, fields, public ) VALUES ( ?, ?, [hex] )`,
					data[`request_id`], fval, c.r.FormValue(`publicKey`))
			}
		}
	}
	if err != nil {
		result.Error = err.Error()
	} else {
		result.Result = true
	}

	return result
}
