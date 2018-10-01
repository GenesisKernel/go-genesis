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

package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/GenesisKernel/go-genesis/packages/consts"
	"github.com/GenesisKernel/go-genesis/packages/converter"
	"github.com/GenesisKernel/go-genesis/packages/language"
	"github.com/GenesisKernel/go-genesis/packages/model"

	log "github.com/sirupsen/logrus"
)

func sections(w http.ResponseWriter, r *http.Request, data *apiData, logger *log.Entry) (err error) {
	var limit int

	table := converter.EscapeName(getPrefix(data) + `_sections`)
	count, err := model.GetRecordsCountTx(nil, strings.Trim(table, `"`))
	if err != nil {
		logger.WithFields(log.Fields{"type": consts.DBError, "error": err, "table": table}).Error("Getting table records count")
		return errorAPI(w, err.Error(), http.StatusInternalServerError)
	}

	if data.params[`limit`].(int64) > 0 {
		limit = int(data.params[`limit`].(int64))
	} else {
		limit = 25
	}
	list, err := model.GetAll(`select * from `+table+` order by id desc`+
		fmt.Sprintf(` offset %d `, data.params[`offset`].(int64)), limit)
	if err != nil {
		logger.WithFields(log.Fields{"type": consts.DBError, "error": err, "table": table}).Error("Getting rows from table")
		return errorAPI(w, err, http.StatusInternalServerError)
	}
	var result []map[string]string
	for _, item := range list {
		var roles []int64
		if err := json.Unmarshal([]byte(item["roles_access"]), &roles); err != nil {
			return errorAPI(w, err, http.StatusInternalServerError)
		}
		var added bool
		if len(roles) > 0 {
			for _, v := range roles {
				if v == data.roleId {
					added = true
					break
				}
			}
		} else {
			added = true
		}
		if added {
			item["title"] = language.LangMacro(item["title"], int(data.ecosystemId),
				1, r.Header.Get(`Accept-Language`), false)
			result = append(result, item)
		}
	}
	data.result = &listResult{
		Count: converter.Int64ToStr(count), List: result,
	}
	return
}
