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
	"fmt"
	"net/url"
	"strings"
	"testing"

	"github.com/GenesisKernel/go-genesis/packages/crypto"

	"github.com/stretchr/testify/assert"
)

type smartParams struct {
	Params  map[string]string
	Results map[string]string
}

type smartContract struct {
	Name   string
	Value  string
	Params []smartParams
}

func TestUpperName(t *testing.T) {
	if err := keyLogin(1); err != nil {
		t.Error(err)
		return
	}
	rnd := crypto.RandSeq(4)
	form := url.Values{"Name": {"testTable" + rnd}, "Columns": {`[{"name":"num","type":"text",   "conditions":"true"},
	{"name":"text", "type":"text","conditions":"true"}]`},
		"Permissions": {`{"insert": "true", "update" : "true", "new_column": "true"}`}}
	err := postTx(`NewTable`, &form)
	if err != nil {
		t.Error(err)
		return
	}
	form = url.Values{`Value`: {`contract AddRow` + rnd + ` {
		data {
		}
		conditions {
		}
		action {
		   DBInsert("testTable` + rnd + `", "num, text", "fgdgf", "124234") 
		}
	}`}, `Conditions`: {`true`}}
	if err := postTx(`NewContract`, &form); err != nil {
		t.Error(err)
		return
	}
	if err := postTx(`AddRow`+rnd, &url.Values{}); err != nil {
		t.Error(err)
		return
	}
}

func TestSmartFields(t *testing.T) {

	if err := keyLogin(1); err != nil {
		t.Error(err)
		return
	}
	var cntResult getContractResult
	err := sendGet(`contract/MainCondition`, nil, &cntResult)
	if err != nil {
		t.Error(err)
		return
	}
	if len(cntResult.Fields) != 0 {
		t.Error(`MainCondition fields must be empty`)
		return
	}
	if cntResult.Name != `@1MainCondition` {
		t.Errorf(`MainCondition name is wrong: %s`, cntResult.Name)
		return
	}
	if err := postTx(`MainCondition`, &url.Values{}); err != nil {
		t.Error(err)
		return
	}
}

func TestMoneyTransfer(t *testing.T) {
	if err := keyLogin(1); err != nil {
		t.Error(err)
		return
	}

	form := url.Values{`Amount`: {`53330000`}, `Recipient`: {`0005-2070-2000-0006-0200`}}
	if err := postTx(`MoneyTransfer`, &form); err != nil {
		t.Error(err)
		return
	}
	form = url.Values{`Amount`: {`2440000`}, `Recipient`: {`1109-7770-3360-6764-7059`}, `Comment`: {`Test`}}
	if err := postTx(`MoneyTransfer`, &form); err != nil {
		t.Error(err)
		return
	}
	form = url.Values{`Amount`: {`53330000`}, `Recipient`: {`0005207000`}}
	if err := postTx(`MoneyTransfer`, &form); cutErr(err) != `{"type":"error","error":"Recipient 0005207000 is invalid"}` {
		t.Error(err)
		return
	}
}

func TestPage(t *testing.T) {

	if err := keyLogin(1); err != nil {
		t.Error(err)
		return
	}
	name := randName(`page`)
	menuname := randName(`menu`)
	menu := `government`
	value := `P(test,test paragraph)`

	form := url.Values{"Name": {name}, "Value": {`Param Value`},
		"Conditions": {`ContractConditions("MainCondition")`}}
	err := postTx(`NewParameter`, &form)
	if err != nil {
		t.Error(err)
		return
	}
	err = postTx(`NewParameter`, &form)
	if cutErr(err) != fmt.Sprintf(`{"type":"warning","error":"Parameter %s already exists"}`, name) {
		t.Error(err)
		return
	}

	form = url.Values{"Name": {menuname}, "Value": {`first
			second
			third`}, "Title": {`My Menu`},
		"Conditions": {`true`}}
	err = postTx(`NewMenu`, &form)
	if err != nil {
		t.Error(err)
		return
	}
	err = postTx(`NewMenu`, &form)
	if cutErr(err) != fmt.Sprintf(`{"type":"warning","error":"Menu %s already exists"}`, menuname) {
		t.Error(err)
		return
	}

	form = url.Values{"Id": {`7123`}, "Value": {`New Param Value`},
		"Conditions": {`ContractConditions("MainCondition")`}}
	err = postTx(`EditParameter`, &form)
	if cutErr(err) != `{"type":"panic","error":"Item 7123 has not been found"}` {
		t.Error(err)
		return
	}
	form = url.Values{"Id": {`16`}, "Value": {`Changed Param Value`},
		"Conditions": {`ContractConditions("MainCondition")`}}
	err = postTx(`EditParameter`, &form)
	if err != nil {
		t.Error(err)
		return
	}

	name = randName(`page`)
	form = url.Values{"Name": {name}, "Value": {value},
		"Menu": {menu}, "Conditions": {"ContractConditions(`MainCondition`)"}}
	err = postTx(`NewPage`, &form)
	if err != nil {
		t.Error(err)
		return
	}
	err = postTx(`NewPage`, &form)
	if cutErr(err) != fmt.Sprintf(`{"type":"warning","error":"Page %s already exists"}`, name) {
		t.Error(err)
		return
	}

	form = url.Values{"Name": {name}, "Value": {value},
		"Conditions": {"ContractConditions(`MainCondition`)"}}
	err = postTx(`NewBlock`, &form)
	if err != nil {
		t.Error(err)
		return
	}
	err = postTx(`NewBlock`, &form)
	if err.Error() != fmt.Sprintf(`{"type":"warning","error":"Block %s already exists"}`, name) {
		t.Error(err)
		return
	}
	form = url.Values{"Id": {`1`}, "Name": {name}, "Value": {value},
		"Conditions": {"ContractConditions(`MainCondition`)"}}
	err = postTx(`EditBlock`, &form)
	if err != nil {
		t.Error(err)
		return
	}

	form = url.Values{"Id": {`1`}, "Value": {value + `Span(Test)`},
		"Menu": {menu}, "Conditions": {"ContractConditions(`MainCondition`)"}}
	err = postTx(`EditPage`, &form)
	if err != nil {
		t.Error(err)
		return
	}
	form = url.Values{"Id": {`1112`}, "Value": {value + `Span(Test)`},
		"Menu": {menu}, "Conditions": {"ContractConditions(`MainCondition`)"}}
	err = postTx(`EditPage`, &form)
	if cutErr(err) != `{"type":"panic","error":"Item 1112 has not been found"}` {
		t.Error(err)
		return
	}

	form = url.Values{"Id": {`1`}, "Value": {`Span(Append)`}}
	err = postTx(`AppendPage`, &form)
	if err != nil {
		t.Error(err)
		return
	}
}

func TestNewTable(t *testing.T) {
	if err := keyLogin(1); err != nil {
		t.Error(err)
		return
	}
	name := randName(`tbl`)
	form := url.Values{"Name": {`1_` + name}, "Columns": {`[{"name":"MyName","type":"varchar", 
		"conditions":"true"},
	  {"name":"Name", "type":"varchar","index": "0", "conditions":"true"}]`},
		"Permissions": {`{"insert": "true", "update" : "true", "new_column": "true"}`}}
	err := postTx(`NewTable`, &form)
	if err != nil {
		t.Error(err)
		return
	}
	form = url.Values{"TableName": {`1_` + name}, "Name": {`newCol`},
		"Type": {"varchar"}, "Index": {"0"}, "Permissions": {"true"}}
	err = postTx(`NewColumn`, &form)
	if err != nil {
		t.Error(err)
		return
	}
	form = url.Values{`Value`: {`contract sub` + name + ` {
		action {
			DBInsert("1_` + name + `", "name", "ok")
			DBUpdate("1_` + name + `", 1, "name", "test value" )
			$result = DBFind("1_` + name + `").Columns("name").WhereId(1).One("name")
		}
	}`}, `Conditions`: {`true`}}
	err = postTx(`NewContract`, &form)
	if err != nil {
		t.Error(err)
		return
	}
	_, msg, err := postTxResult(`sub`+name, &url.Values{})
	if err != nil {
		t.Error(err)
		return
	}
	if msg != `test value` {
		t.Errorf("wrong result %s", msg)
		return
	}

	form = url.Values{"Name": {name}, "Columns": {`[{"name":"MyName","type":"varchar", "index": "1", 
	  "conditions":"true"},
	{"name":"Amount", "type":"number","index": "0", "conditions":"true"},
	{"name":"Doc", "type":"json","index": "0", "conditions":"true"},	
	{"name":"Active", "type":"character","index": "0", "conditions":"true"}]`},
		"Permissions": {`{"insert": "true", "update" : "true", "new_column": "true"}`}}
	err = postTx(`NewTable`, &form)
	if err != nil {
		t.Error(err)
		return
	}
	err = postTx(`NewTable`, &form)
	if err.Error() != fmt.Sprintf(`{"type":"panic","error":"table %s exists"}`, name) {
		t.Error(err)
		return
	}
	form = url.Values{"Name": {name},
		"Permissions": {`{"insert": "ContractConditions(\"MainCondition\")",
				"update" : "true", "new_column": "ContractConditions(\"MainCondition\")"}`}}
	err = postTx(`EditTable`, &form)
	if err != nil {
		t.Error(err)
		return
	}
	form = url.Values{"TableName": {name}, "Name": {`newDoc`},
		"Type": {"json"}, "Index": {"0"}, "Permissions": {"true"}}
	err = postTx(`NewColumn`, &form)
	if err != nil {
		t.Error(err)
		return
	}
	form = url.Values{"TableName": {name}, "Name": {`newCol`},
		"Type": {"varchar"}, "Index": {"0"}, "Permissions": {"true"}}
	err = postTx(`NewColumn`, &form)
	if err != nil {
		t.Error(err)
		return
	}
	err = postTx(`NewColumn`, &form)
	if err.Error() != `{"type":"panic","error":"column newcol exists"}` {
		t.Error(err)
		return
	}
	form = url.Values{"TableName": {name}, "Name": {`newCol`},
		"Permissions": {"ContractConditions(\"MainCondition\")"}}
	err = postTx(`EditColumn`, &form)
	if err != nil {
		t.Error(err)
		return
	}
	upname := strings.ToUpper(name)
	form = url.Values{"TableName": {upname}, "Name": {`UPCol`},
		"Type": {"varchar"}, "Index": {"0"}, "Permissions": {"true"}}
	err = postTx(`NewColumn`, &form)
	if err != nil {
		t.Error(err)
		return
	}
	form = url.Values{"TableName": {upname}, "Name": {`upCOL`},
		"Permissions": {"ContractConditions(\"MainCondition\")"}}
	err = postTx(`EditColumn`, &form)
	if err != nil {
		t.Error(err)
		return
	}
	form = url.Values{"Name": {upname},
		"Permissions": {`{"insert": "ContractConditions(\"MainCondition\")", 
			"update" : "true", "new_column": "ContractConditions(\"MainCondition\")"}`}}
	err = postTx(`EditTable`, &form)
	if err != nil {
		t.Error(err)
		return
	}
	var ret tablesResult
	err = sendGet(`tables`, nil, &ret)
	if err != nil {
		t.Error(err)
		return
	}
}

type invalidPar struct {
	Name  string
	Value string
}

func TestUpdateSysParam(t *testing.T) {
	if err := keyLogin(1); err != nil {
		t.Error(err)
		return
	}
	form := url.Values{"Name": {`max_columns`}, "Value": {`49`}}
	err := postTx(`UpdateSysParam`, &form)
	if err != nil {
		t.Error(err)
		return
	}
	var sysList ecosystemParamsResult
	err = sendGet(`systemparams?names=max_columns`, nil, &sysList)
	if err != nil {
		t.Error(err)
		return
	}
	if len(sysList.List) != 1 || sysList.List[0].Value != `49` {
		t.Error(`Wrong max_column value`)
		return
	}
	name := randName(`test`)
	form = url.Values{"Name": {name}, "Value": {`contract ` + name + ` {
		action { 
			var costlen int
			costlen = SysParamInt("extend_cost_len") + 1
			UpdateSysParam("Name,Value","max_columns","51")
			DBUpdateSysParam("extend_cost_len", Str(costlen), "true" )
			if SysParamInt("extend_cost_len") != costlen {
				error "Incorrect updated value"
			}
			DBUpdateSysParam("max_indexes", "4", "false" )
		}
		}`},
		"Conditions": {`ContractConditions("MainCondition")`}}
	err = postTx("NewContract", &form)
	if err != nil {
		t.Error(err)
		return
	}
	err = postTx(name, &form)
	if err != nil {
		if err.Error() != `{"type":"panic","error":"Access denied"}` {
			t.Error(err)
			return
		}
	}
	err = sendGet(`systemparams?names=max_columns,max_indexes`, nil, &sysList)
	if err != nil {
		t.Error(err)
		return
	}
	if len(sysList.List) != 2 || !((sysList.List[0].Value == `51` && sysList.List[1].Value == `4`) ||
		(sysList.List[0].Value == `4` && sysList.List[1].Value == `51`)) {
		t.Error(`Wrong max_column or max_indexes value`)
		return
	}
	err = postTx(name, &form)
	if err == nil || err.Error() != `{"type":"panic","error":"Access denied"}` {
		t.Error(`incorrect access to system parameter`)
		return
	}
	notvalid := []invalidPar{
		{`gap_between_blocks`, `100000`},
		{`rb_blocks_1`, `-1`},
		{`page_price`, `-20`},
		{`max_block_size`, `0`},
		{`max_fuel_tx`, `20string`},
		{`fuel_rate`, `string`},
		{`fuel_rate`, `[test]`},
		{`fuel_rate`, `[["name", "100"]]`},
		{`commission_wallet`, `[["1", "0"]]`},
		{`commission_wallet`, `[{"1", "50"}]`},
		{`full_nodes`, `[["", "100", "c1a9e7b2fb8cea2a272e183c3e27e2d59a3ebe613f51873a46885c9201160bd263ef43b583b631edd1284ab42483712fd2ccc40864fe9368115ceeee47a7"]]`},
		{`full_nodes`, `[["127.0.0.1", "0", "c1a9e7b2fb8cea2a272e183c3e27e2d59a3ebe613f51873a46885c9201160bd263ef43b583b631edd1284ab42483712fd2ccc40864fe9368115ceeee47a7c7d0"]]`},
	}
	for _, item := range notvalid {
		err = postTx(`UpdateSysParam`, &url.Values{`Name`: {item.Name}, `Value`: {item.Value}})
		if err == nil {
			t.Error(`must be invalid ` + item.Value)
			return
		}
		err = sendGet(`systemparams?names=`+item.Name, nil, &sysList)
		if err != nil {
			t.Error(err)
			return
		}
		if len(sysList.List) != 1 {
			t.Error(`have got wrong parameter ` + item.Name)
			return
		}
		err = postTx(`UpdateSysParam`, &url.Values{`Name`: {item.Name}, `Value`: {sysList.List[0].Value}})
		if err != nil {
			fmt.Println(item.Name, sysList.List[0].Value, sysList.List[0])
			t.Error(err)
			return
		}
	}
}

func TestValidateConditions(t *testing.T) {
	if err := keyLogin(1); err != nil {
		t.Error(err)
		return
	}

	baseForm := url.Values{"Id": {"1"}, "Value": {"Test"}, "Conditions": {"incorrectConditions"}}
	contracts := map[string]url.Values{
		"EditContract":  baseForm,
		"EditParameter": baseForm,
		"EditMenu":      baseForm,
		"EditPage":      url.Values{"Id": {"1"}, "Value": {"Test"}, "Conditions": {"incorrectConditions"}, "Menu": {"1"}},
	}
	expectedErr := `{"type":"panic","error":"unknown identifier incorrectConditions"}`

	for contract, form := range contracts {
		err := postTx(contract, &form)
		if err.Error() != expectedErr {
			t.Errorf("contract %s expected '%s' got '%s'", contract, expectedErr, err)
			return
		}
	}
}

func TestDBMetric(t *testing.T) {
	if err := keyLogin(1); err != nil {
		t.Error(err)
		return
	}

	name := randName("Metrics")
	form := url.Values{
		"Value": {`
			contract ` + name + ` {
				data {}
				conditions {}
				action {
					DBSelectMetrics("ecosystem_pages", "1 days", "max")
				}
			}`},
		"Conditions": {"true"},
	}
	if err := postTx("NewContract", &form); err != nil {
		t.Error(err)
		return
	}
	if err := postTx(name, &url.Values{}); err != nil {
		t.Error(err)
		return
	}
}

func TestPartitialEdit(t *testing.T) {
	if err := keyLogin(1); err != nil {
		t.Error(err)
		return
	}
	name := randName(`part`)
	form := url.Values{"Name": {name}, "Value": {"Span(Original text)"},
		"Menu": {"original_menu"}, "Conditions": {"ContractConditions(`MainCondition`)"}}
	err := postTx(`NewPage`, &form)
	if err != nil {
		t.Error(err)
		return
	}
	var retList listResult
	err = sendGet(`list/pages`, nil, &retList)
	if err != nil {
		t.Error(err)
		return
	}
	idItem := retList.Count
	value := `Span(Temp)`
	menu := `temp_menu`
	err = postTx(`EditPage`, &url.Values{"Id": {idItem}, "Value": {value},
		"Menu": {menu}})
	if err != nil {
		t.Error(err)
		return
	}
	var ret rowResult
	err = sendGet(`row/pages/`+idItem, nil, &ret)
	if err != nil {
		t.Error(err)
		return
	}
	if ret.Value["value"] != value || ret.Value["menu"] != menu {
		t.Errorf(`wrong value or menu`)
		return
	}
	value = `Span(Updated)`
	menu = `default_menu`
	conditions := `true`
	err = postTx(`EditPage`, &url.Values{"Id": {idItem}, "Value": {value}})
	if err != nil {
		t.Error(err)
		return
	}
	err = postTx(`EditPage`, &url.Values{"Id": {idItem}, "Menu": {menu}})
	if err != nil {
		t.Error(err)
		return
	}
	err = postTx(`EditPage`, &url.Values{"Id": {idItem}, "Conditions": {conditions}})
	if err != nil {
		t.Error(err)
		return
	}
	err = sendGet(`row/pages/`+idItem, nil, &ret)
	if err != nil {
		t.Error(err)
		return
	}
	if ret.Value["value"] != value || ret.Value["menu"] != menu ||
		ret.Value["conditions"] != conditions {
		t.Errorf(`wrong page parameters`)
		return
	}

	form = url.Values{"Name": {name}, "Value": {`MenuItem(One)`}, "Title": {`My Menu`},
		"Conditions": {"ContractConditions(`MainCondition`)"}}
	err = postTx(`NewMenu`, &form)
	if err != nil {
		t.Error(err)
		return
	}
	err = sendGet(`list/menu`, nil, &retList)
	if err != nil {
		t.Error(err)
		return
	}
	idItem = retList.Count
	value = `MenuItem(Two)`
	err = postTx(`EditMenu`, &url.Values{"Id": {idItem}, "Value": {value}})
	if err != nil {
		t.Error(err)
		return
	}
	err = postTx(`EditMenu`, &url.Values{"Id": {idItem}, "Conditions": {conditions}})
	if err != nil {
		t.Error(err)
		return
	}
	err = sendGet(`row/menu/`+idItem, nil, &ret)
	if err != nil {
		t.Error(err)
		return
	}
	if ret.Value["value"] != value || ret.Value["conditions"] != conditions {
		t.Errorf(`wrong menu parameters`)
		return
	}

	form = url.Values{"Name": {name}, "Value": {`Span(Block)`},
		"Conditions": {"ContractConditions(`MainCondition`)"}}
	err = postTx(`NewBlock`, &form)
	if err != nil {
		t.Error(err)
		return
	}
	err = sendGet(`list/blocks`, nil, &retList)
	if err != nil {
		t.Error(err)
		return
	}
	idItem = retList.Count
	value = `Span(Updated block)`
	err = postTx(`EditBlock`, &url.Values{"Id": {idItem}, "Value": {value}})
	if err != nil {
		t.Error(err)
		return
	}
	err = postTx(`EditBlock`, &url.Values{"Id": {idItem}, "Conditions": {conditions}})
	if err != nil {
		t.Error(err)
		return
	}
	err = sendGet(`row/blocks/`+idItem, nil, &ret)
	if err != nil {
		t.Error(err)
		return
	}
	if ret.Value["value"] != value || ret.Value["conditions"] != conditions {
		t.Errorf(`wrong block parameters`)
		return
	}

}

func TestContractEdit(t *testing.T) {
	if err := keyLogin(1); err != nil {
		t.Error(err)
		return
	}
	name := randName(`part`)
	form := url.Values{"Value": {`contract ` + name + ` {
		    action {
				$result = "before"
			}
		}`},
		"Conditions": {"ContractConditions(`MainCondition`)"}}
	err := postTx(`NewContract`, &form)
	if err != nil {
		t.Error(err)
		return
	}
	var retList listResult
	err = sendGet(`list/contracts`, nil, &retList)
	if err != nil {
		t.Error(err)
		return
	}
	idItem := retList.Count
	value := `contract ` + name + ` {
		action {
			$result = "after"
		}
	}`
	conditions := `true`
	wallet := "1231234123412341230"
	err = postTx(`EditContract`, &url.Values{"Id": {idItem}, "Value": {value}})
	if err != nil {
		t.Error(err)
		return
	}
	err = postTx(`EditContract`, &url.Values{"Id": {idItem}, "Conditions": {conditions},
		"WalletId": {wallet}})
	if err != nil {
		t.Error(err)
		return
	}
	var ret rowResult
	err = sendGet(`row/contracts/`+idItem, nil, &ret)
	if err != nil {
		t.Error(err)
		return
	}
	if ret.Value["value"] != value || ret.Value["conditions"] != conditions ||
		ret.Value["wallet_id"] != wallet {
		t.Errorf(`wrong parameters of contract`)
		return
	}
	_, msg, err := postTxResult(name, &url.Values{})
	if err != nil {
		t.Error(err)
		return
	}
	if msg != "after" {
		t.Errorf(`the wrong result of the contract %s`, msg)
	}
}

func TestDelayedContracts(t *testing.T) {
	if err := keyLogin(1); err != nil {
		t.Error(err)
		return
	}

	form := url.Values{
		"Contract":   {"UnknownContract"},
		"EveryBlock": {"10"},
		"Limit":      {"2"},
		"Conditions": {"true"},
	}
	err := postTx("NewDelayedContract", &form)
	assert.EqualError(t, err, `{"type":"error","error":"Unknown contract @1UnknownContract"}`)

	form.Set("Contract", "MainCondition")
	err = postTx("NewDelayedContract", &form)
	assert.NoError(t, err)

	form.Set("BlockID", "1")
	err = postTx("NewDelayedContract", &form)
	assert.EqualError(t, err, `{"type":"error","error":"The blockID must be greater than the current blockID"}`)

	form = url.Values{
		"Id":         {"1"},
		"Contract":   {"MainCondition"},
		"EveryBlock": {"10"},
		"Conditions": {"true"},
		"Deleted":    {"1"},
	}
	err = postTx("EditDelayedContract", &form)
	assert.NoError(t, err)
}
