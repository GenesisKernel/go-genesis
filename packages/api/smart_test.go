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
	"math/rand"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/GenesisKernel/go-genesis/packages/converter"
	"github.com/GenesisKernel/go-genesis/packages/crypto"
	"github.com/shopspring/decimal"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
	form := url.Values{"Name": {"testTable" + rnd}, "ApplicationId": {"1"}, "Columns": {`[{"name":"num","type":"text",   "conditions":"true"},
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
	}`}, "ApplicationId": {"1"}, `Conditions`: {`true`}}
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
	size := 1000000
	big := make([]byte, size)
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < size; i++ {
		big[i] = '0' + byte(rand.Intn(10))
	}
	form = url.Values{`Amount`: {string(big)}, `Recipient`: {`0005-2070-2000-0006-0200`}}
	if err := postTx(`MoneyTransfer`, &form); err.Error() != `400 {"error": "E_LIMITFORSIGN", "msg": "Length of forsign is too big (1000106)" , "params": ["1000106"]}` {
		t.Error(err)
		return
	}
}

func TestPage(t *testing.T) {
	assert.NoError(t, keyLogin(1))

	name := randName(`page`)
	menuname := randName(`menu`)
	menu := `government`
	value := `P(test,test paragraph)`

	form := url.Values{"Name": {name}, "Value": {`Param Value`},
		"Conditions": {`ContractConditions("MainCondition")`}}
	assert.NoError(t, postTx(`NewParameter`, &form))

	err := postTx(`NewParameter`, &form)
	assert.Equal(t, fmt.Sprintf(`{"type":"warning","error":"Parameter %s already exists"}`, name), cutErr(err))

	form = url.Values{"Name": {menuname}, "Value": {`first
			second
			third`}, "Title": {`My Menu`},
		"Conditions": {`true`}}
	assert.NoError(t, postTx(`NewMenu`, &form))

	err = postTx(`NewMenu`, &form)
	assert.Equal(t, fmt.Sprintf(`{"type":"warning","error":"Menu %s already exists"}`, menuname), cutErr(err))

	form = url.Values{"Id": {`7123`}, "Value": {`New Param Value`},
		"Conditions": {`ContractConditions("MainCondition")`}}
	err = postTx(`EditParameter`, &form)
	assert.Equal(t, `{"type":"panic","error":"Item 7123 has not been found"}`, cutErr(err))

	form = url.Values{"Id": {`16`}, "Value": {`Changed Param Value`},
		"Conditions": {`ContractConditions("MainCondition")`}}
	assert.NoError(t, postTx(`EditParameter`, &form))

	name = randName(`page`)
	form = url.Values{"Name": {name}, "Value": {value}, "ApplicationId": {`1`},
		"Menu": {menu}, "Conditions": {"ContractConditions(`MainCondition`)"}}
	assert.NoError(t, postTx(`NewPage`, &form))

	err = postTx(`NewPage`, &form)
	assert.Equal(t, fmt.Sprintf(`{"type":"warning","error":"Page %s already exists"}`, name), cutErr(err))
	err = postTx(`NewPage`, &form)
	if cutErr(err) != fmt.Sprintf(`{"type":"warning","error":"Page %s already exists"}`, name) {
		t.Error(err)
		return
	}
	form = url.Values{"Name": {`app` + name}, "Value": {value}, "ValidateCount": {"2"},
		"ValidateMode": {"1"}, "ApplicationId": {`1`},
		"Menu": {menu}, "Conditions": {"ContractConditions(`MainCondition`)"}}
	err = postTx(`NewPage`, &form)
	if err != nil {
		t.Error(err)
		return
	}

	var ret listResult
	err = sendGet(`list/pages`, nil, &ret)
	if err != nil {
		t.Error(err)
		return
	}
	id := ret.Count
	form = url.Values{"Id": {id}, "ValidateCount": {"2"}, "ValidateMode": {"1"}}
	err = postTx(`EditPage`, &form)
	if err != nil {
		t.Error(err)
		return
	}
	var row rowResult
	err = sendGet(`row/pages/`+id, nil, &row)
	if err != nil {
		t.Error(err)
		return
	}

	if row.Value["validate_mode"] != `1` {
		t.Errorf(`wrong validate value %s`, row.Value["validate_mode"])
		return
	}

	form = url.Values{"Id": {id}, "Value": {value}, "ValidateCount": {"1"},
		"ValidateMode": {"0"}}
	err = postTx(`EditPage`, &form)
	if err != nil {
		t.Error(err)
		return
	}
	err = sendGet(`row/pages/`+id, nil, &row)
	if err != nil {
		t.Error(err)
		return
	}
	if row.Value["validate_mode"] != `0` {
		t.Errorf(`wrong validate value %s`, row.Value["validate_mode"])
		return
	}

	form = url.Values{"Id": {id}, "Value": {value}, "ValidateCount": {"1"},
		"ValidateMode": {"0"}}
	err = postTx(`EditPage`, &form)
	if err != nil {
		t.Error(err)
		return
	}
	err = sendGet(`row/pages/`+id, nil, &row)
	if err != nil {
		t.Error(err)
		return
	}
	if row.Value["validate_mode"] != `0` {
		t.Errorf(`wrong validate value %s`, row.Value["validate_mode"])
		return
	}

	form = url.Values{"Name": {name}, "Value": {value}, "ApplicationId": {`1`},
		"Conditions": {"ContractConditions(`MainCondition`)"}}
	assert.NoError(t, postTx(`NewBlock`, &form))

	err = postTx(`NewBlock`, &form)
	assert.EqualError(t, err, fmt.Sprintf(`{"type":"warning","error":"Block %s already exists"}`, name))

	form = url.Values{"Id": {`1`}, "Name": {name}, "Value": {value},
		"Conditions": {"ContractConditions(`MainCondition`)"}}
	assert.NoError(t, postTx(`EditBlock`, &form))

	form = url.Values{"Id": {`1`}, "Value": {value + `Span(Test)`},
		"Menu": {menu}, "Conditions": {"ContractConditions(`MainCondition`)"}}
	assert.NoError(t, postTx(`EditPage`, &form))

	form = url.Values{"Id": {`1112`}, "Value": {value + `Span(Test)`},
		"Menu": {menu}, "Conditions": {"ContractConditions(`MainCondition`)"}}
	err = postTx(`EditPage`, &form)
	assert.Equal(t, `{"type":"panic","error":"Item 1112 has not been found"}`, cutErr(err))

	form = url.Values{"Id": {`1`}, "Value": {`Span(Append)`}}
	assert.NoError(t, postTx(`AppendPage`, &form))
}

func TestNewTableOnly(t *testing.T) {
	assert.NoError(t, keyLogin(1))

	name := "MMy_s_test_table"
	form := url.Values{"Name": {name}, "ApplicationId": {"1"}, "Columns": {`[{"name":"MyName","type":"varchar", 
		"conditions":"true"},
	  {"name":"Name", "type":"varchar","index": "0", "conditions":"{\"read\":\"true\",\"update\":\"true\"}"}]`},
		"Permissions": {`{"insert": "true", "update" : "true", "new_column": "true"}`}}
	require.NoError(t, postTx(`NewTable`, &form))

	var ret tableResult
	require.NoError(t, sendGet(`table/`+name, nil, &ret))
	fmt.Printf("%+v\n", ret)
}

func TestDBFind(t *testing.T) {
	assert.NoError(t, keyLogin(1))

	name := randName(`tbl`)
	form := url.Values{"Name": {name}, "ApplicationId": {"1"}, "Columns": {`[{"name":"txt","type":"varchar", 
		"conditions":"true"},
	  {"name":"Name", "type":"varchar","index": "0", "conditions":"{\"read\":\"true\",\"update\":\"true\"}"}]`},
		"Permissions": {`{"insert": "true", "update" : "true", "new_column": "true"}`}}
	assert.NoError(t, postTx(`NewTable`, &form))

	form = url.Values{`Value`: {`contract sub` + name + ` {
		action {
			DBInsert("` + name + `", "txt,name", "ok", "thisis")
			DBInsert("` + name + `", "txt,name", "текст", "заголовок")
			$result = DBFind("` + name + `").Columns("name").Where("txt=?", "текст").One("name")
		}
	}`}, `Conditions`: {`true`}, "ApplicationId": {"1"}}
	assert.NoError(t, postTx(`NewContract`, &form))

	_, ret, err := postTxResult(`sub`+name, &url.Values{})
	assert.Equal(t, `заголовок`, ret)

	var retPage contentResult
	value := `DBFind(` + name + `, src).Columns(name).Where(txt='текст')`
	form = url.Values{"Name": {name}, "Value": {value}, "ApplicationId": {`1`},
		"Menu": {`default_menu`}, "Conditions": {"ContractConditions(`MainCondition`)"}}
	assert.NoError(t, postTx(`NewPage`, &form))

	assert.NoError(t, sendPost(`content/page/`+name, &url.Values{}, &retPage))
	if err != nil {
		t.Error(err)
		return
	}
	if RawToString(retPage.Tree) != `[{"tag":"dbfind","attr":{"columns":["name","id"],"data":[["заголовок","2"]],"name":"`+name+`","source":"src","types":["text","text"],"where":"txt='текст'"}}]` {
		t.Error(fmt.Errorf(`wrong tree %s`, RawToString(retPage.Tree)))
		return
	}
}

func TestNewTable(t *testing.T) {
	assert.NoError(t, keyLogin(1))

	name := randName(`tbl`)
	form := url.Values{"Name": {`1_` + name}, "ApplicationId": {"1"}, "Columns": {`[{"name":"MyName","type":"varchar", 
		"conditions":"true"},
	  {"name":"Name", "type":"varchar","index": "0", "conditions":"{\"read\":\"true\",\"update\":\"true\"}"}]`},
		"Permissions": {`{"insert": "true", "update" : "true", "new_column": "true"}`}}
	assert.NoError(t, postTx(`NewTable`, &form))

	form = url.Values{"TableName": {`1_` + name}, "Name": {`newCol`},
		"Type": {"varchar"}, "Index": {"0"}, "Permissions": {"true"}}
	assert.NoError(t, postTx(`NewColumn`, &form))

	form = url.Values{`Value`: {`contract sub` + name + ` {
		action {
			DBInsert("1_` + name + `", "name", "ok")
			DBUpdate("1_` + name + `", 1, "name", "test value" )
			$result = DBFind("1_` + name + `").Columns("name").WhereId(1).One("name")
		}
	}`}, `Conditions`: {`true`}, "ApplicationId": {"1"}}
	assert.NoError(t, postTx(`NewContract`, &form))

	_, msg, err := postTxResult(`sub`+name, &url.Values{})
	assert.NoError(t, err)
	assert.Equal(t, msg, "test value")

	form = url.Values{"Name": {name}, "ApplicationId": {"1"}, "Columns": {`[{"name":"MyName","type":"varchar", "index": "1", 
	  "conditions":"true"},
	{"name":"Amount", "type":"number","index": "0", "conditions":"true"},
	{"name":"Doc", "type":"json","index": "0", "conditions":"true"},	
	{"name":"Active", "type":"character","index": "0", "conditions":"true"}]`},
		"Permissions": {`{"insert": "true", "update" : "true", "new_column": "true"}`}}
	assert.NoError(t, postTx(`NewTable`, &form))

	assert.EqualError(t, postTx(`NewTable`, &form), fmt.Sprintf(`{"type":"panic","error":"table %s exists"}`, name))

	form = url.Values{"Name": {name},
		"Permissions": {`{"insert": "ContractConditions(\"MainCondition\")",
				"update" : "true", "new_column": "ContractConditions(\"MainCondition\")"}`}}
	assert.NoError(t, postTx(`EditTable`, &form))

	form = url.Values{"TableName": {name}, "Name": {`newDoc`},
		"Type": {"json"}, "Index": {"0"}, "Permissions": {"true"}}
	assert.NoError(t, postTx(`NewColumn`, &form))

	form = url.Values{"TableName": {name}, "Name": {`newCol`},
		"Type": {"varchar"}, "Index": {"0"}, "Permissions": {"true"}}
	assert.NoError(t, postTx(`NewColumn`, &form))

	err = postTx(`NewColumn`, &form)
	if err.Error() != `{"type":"panic","error":"column newcol exists"}` {
		t.Error(err)
		return
	}
	form = url.Values{"TableName": {name}, "Name": {`newCol`},
		"Permissions": {"ContractConditions(\"MainCondition\")"}}
	assert.NoError(t, postTx(`EditColumn`, &form))

	upname := strings.ToUpper(name)
	form = url.Values{"TableName": {upname}, "Name": {`UPCol`},
		"Type": {"varchar"}, "Index": {"0"}, "Permissions": {"true"}}
	assert.NoError(t, postTx(`NewColumn`, &form))

	form = url.Values{"TableName": {upname}, "Name": {`upCOL`},
		"Permissions": {"ContractConditions(\"MainCondition\")"}}
	assert.NoError(t, postTx(`EditColumn`, &form))

	form = url.Values{"Name": {upname},
		"Permissions": {`{"insert": "ContractConditions(\"MainCondition\")", 
			"update" : "true", "new_column": "ContractConditions(\"MainCondition\")"}`}}
	assert.NoError(t, postTx(`EditTable`, &form))

	var ret tablesResult
	assert.NoError(t, sendGet(`tables`, nil, &ret))
}

type invalidPar struct {
	Name  string
	Value string
}

func TestUpdateSysParam(t *testing.T) {
	assert.NoError(t, keyLogin(1))

	form := url.Values{"Name": {`max_columns`}, "Value": {`49`}}
	assert.NoError(t, postTx(`UpdateSysParam`, &form))

	var sysList ecosystemParamsResult
	assert.NoError(t, sendGet(`systemparams?names=max_columns`, nil, &sysList))
	assert.Len(t, sysList.List, 1)
	assert.Equal(t, "49", sysList.List[0].Value)

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
		}`}, "ApplicationId": {"1"},
		"Conditions": {`ContractConditions("MainCondition")`}}
	assert.NoError(t, postTx("NewContract", &form))

	err := postTx(name, &form)
	if err != nil {
		assert.EqualError(t, err, `{"type":"panic","error":"Access denied"}`)
	}

	assert.NoError(t, sendGet(`systemparams?names=max_columns,max_indexes`, nil, &sysList))
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
		{`full_nodes`, `[["", "http://127.0.0.1", "100", "c1a9e7b2fb8cea2a272e183c3e27e2d59a3ebe613f51873a46885c9201160bd263ef43b583b631edd1284ab42483712fd2ccc40864fe9368115ceeee47a7"]]`},
		{`full_nodes`, `[["127.0.0.1", "", "100", "c1a9e7b2fb8cea2a272e183c3e27e2d59a3ebe613f51873a46885c9201160bd263ef43b583b631edd1284ab42483712fd2ccc40864fe9368115ceeee47a7c7d0"]]`},
		{`full_nodes`, `[["127.0.0.1", "http://127.0.0.1", "0", "c1a9e7b2fb8cea2a272e183c3e27e2d59a3ebe613f51873a46885c9201160bd263ef43b583b631edd1284ab42483712fd2ccc40864fe9368115ceeee47a7c7d0"]]`},
		{"full_nodes", "[]"},
	}
	for _, item := range notvalid {
		assert.Error(t, postTx(`UpdateSysParam`, &url.Values{`Name`: {item.Name}, `Value`: {item.Value}}))
		assert.NoError(t, sendGet(`systemparams?names=`+item.Name, nil, &sysList))
		assert.Len(t, sysList.List, 1, `have got wrong parameter `+item.Name)

		if len(sysList.List[0].Value) == 0 {
			continue
		}

		err = postTx(`UpdateSysParam`, &url.Values{`Name`: {item.Name}, `Value`: {sysList.List[0].Value}})
		assert.NoError(t, err, item.Name, sysList.List[0].Value, sysList.List[0])
	}
}

func TestUpdateFullNodesWithEmptyArray(t *testing.T) {
	require.NoErrorf(t, keyLogin(1), "on login")

	byteNodes := `[`
	byteNodes += `{"tcp_address":"127.0.0.1:7078", "api_address":"https://127.0.0.1:7079", "key_id":"-4466900793776865315", "public_key":"ca901a97e84d76f8d46e2053028f709074b3e60d3e2e33495840586567a0c961820d789592666b67b05c6ae120d5bd83d4388b2f1218638d8226d40ced0bb208"},`
	byteNodes += `{"tcp_address":"127.0.0.1:7080", "api_address":"https://127.0.0.1:7081", "key_id":"542353610328569127", "public_key":"a8ada71764fd2f0c9fa1d2986455288f11f0f3931492d27dc62862fdff9c97c38923ef46679488ad1cd525342d4d974621db58f809be6f8d1c19fdab50abc06b"},`
	byteNodes += `{"tcp_address":"127.0.0.1:7082", "api_address":"https://127.0.0.1:7083", "key_id":"5972241339967729614", "public_key":"de1b74d36ae39422f2478cba591f4d14eb017306f6ffdc3b577cc52ee50edb8fe7c7b2eb191a24c8ddfc567cef32152bab17de698ed7b3f2ab75f3bcc8b9b372"}`
	byteNodes += `]`
	form := &url.Values{
		"Name":  {"full_nodes"},
		"Value": {string(byteNodes)},
	}

	require.NoError(t, postTx(`UpdateSysParam`, form))
}

/*
func TestHelper_InsertNodeKey(t *testing.T) {

	if err := keyLogin(1); err != nil {
		t.Error(err)
		return
	}

	form := url.Values{
		`Value`: {`contract InsertNodeKey {
			data {
				KeyID string
				PubKey string
			}
			conditions {}
			action {
				DBInsert("keys", "id,pub,amount", $KeyID, $PubKey, "100000000000000000000")
			}
		}`},
		`ApplicationId`: {`1`},
		`Conditions`:    {`true`},
	}

	require.NoError(t, postTx(`NewContract`, &form))

	forms := []url.Values{
		url.Values{
			`KeyID`:  {"542353610328569127"},
			`PubKey`: {"be78f54bcf6bb7b49b7ea00790b18b40dd3f5e231ffc764f1c32d3f5a82ab322aee157931bbfca733bac83255002f5ded418f911b959b77a937f0d5d07de74f8"},
		},
		url.Values{
			`KeyID`:  {"5972241339967729614"},
			`PubKey`: {"7b11a9ee4f509903118d5b965a819b778c83a21a52a033e5768d697a70a61a1bad270465f25d7f70683e977be93a9252e762488fc53808a90220d363d0a38eb6"},
		},
	}

	for _, frm := range forms {
		require.NoError(t, postTx(`InsertNodeKey`, &frm))
	}
}
*/
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

func TestDBMetrics(t *testing.T) {
	assert.NoError(t, keyLogin(1))

	contract := randName("Metric")
	form := url.Values{
		"Value": {`
			contract ` + contract + ` {
				data {
					Metric string
				}
				conditions {}
				action {
					UpdateMetrics()
					$result = One(DBSelectMetrics($Metric, "1 days", "max"), "value")
				}
			}`},
		"Conditions": {"true"},
	}
	assert.NoError(t, postTx("NewContract", &form))

	metricValue := func(metric string) int {
		assert.NoError(t, postTx("UpdateMetrics", &url.Values{}))

		_, result, err := postTxResult(contract, &url.Values{"Metric": {metric}})
		assert.NoError(t, err)
		return converter.StrToInt(result)
	}

	ecosystemPages := metricValue("ecosystem_pages")
	ecosystemTx := metricValue("ecosystem_tx")

	form = url.Values{
		"Name":       {randName("page")},
		"Value":      {"P()"},
		"Menu":       {"default_menu"},
		"Conditions": {"true"},
	}
	assert.NoError(t, postTx("NewPage", &form))

	assert.Equal(t, 1, metricValue("ecosystem_pages")-ecosystemPages)
	assert.True(t, metricValue("ecosystem_tx") > ecosystemTx)

}

func TestPartitialEdit(t *testing.T) {
	assert.NoError(t, keyLogin(1))

	name := randName(`part`)
	form := url.Values{"Name": {name}, "Value": {"Span(Original text)"},
		"Menu": {"original_menu"}, "ApplicationId": {"1"}, "Conditions": {"ContractConditions(`MainCondition`)"}}
	assert.NoError(t, postTx(`NewPage`, &form))

	var retList listResult
	assert.NoError(t, sendGet(`list/pages`, nil, &retList))

	idItem := retList.Count
	value := `Span(Temp)`
	menu := `temp_menu`
	assert.NoError(t, postTx(`EditPage`, &url.Values{
		"Id":    {idItem},
		"Value": {value},
		"Menu":  {menu},
	}))

	var ret rowResult
	assert.NoError(t, sendGet(`row/pages/`+idItem, nil, &ret))
	assert.Equal(t, value, ret.Value["value"])
	assert.Equal(t, menu, ret.Value["menu"])

	value = `Span(Updated)`
	menu = `default_menu`
	conditions := `true`
	assert.NoError(t, postTx(`EditPage`, &url.Values{"Id": {idItem}, "Value": {value}}))
	assert.NoError(t, postTx(`EditPage`, &url.Values{"Id": {idItem}, "Menu": {menu}}))
	assert.NoError(t, postTx(`EditPage`, &url.Values{"Id": {idItem}, "Conditions": {conditions}}))
	assert.NoError(t, sendGet(`row/pages/`+idItem, nil, &ret))
	assert.Equal(t, value, ret.Value["value"])
	assert.Equal(t, menu, ret.Value["menu"])

	form = url.Values{"Name": {name}, "Value": {`MenuItem(One)`}, "Title": {`My Menu`},
		"ApplicationId": {"1"}, "Conditions": {"ContractConditions(`MainCondition`)"}}
	assert.NoError(t, postTx(`NewMenu`, &form))
	assert.NoError(t, sendGet(`list/menu`, nil, &retList))
	idItem = retList.Count
	value = `MenuItem(Two)`
	assert.NoError(t, postTx(`EditMenu`, &url.Values{"Id": {idItem}, "Value": {value}}))
	assert.NoError(t, postTx(`EditMenu`, &url.Values{"Id": {idItem}, "Conditions": {conditions}}))
	assert.NoError(t, sendGet(`row/menu/`+idItem, nil, &ret))
	assert.Equal(t, value, ret.Value["value"])
	assert.Equal(t, conditions, ret.Value["conditions"])

	form = url.Values{"Name": {name}, "Value": {`Span(Block)`},
		"ApplicationId": {"1"}, "Conditions": {"ContractConditions(`MainCondition`)"}}
	assert.NoError(t, postTx(`NewBlock`, &form))
	assert.NoError(t, sendGet(`list/blocks`, nil, &retList))
	idItem = retList.Count
	value = `Span(Updated block)`
	assert.NoError(t, postTx(`EditBlock`, &url.Values{"Id": {idItem}, "Value": {value}}))
	assert.NoError(t, postTx(`EditBlock`, &url.Values{"Id": {idItem}, "Conditions": {conditions}}))
	assert.NoError(t, sendGet(`row/blocks/`+idItem, nil, &ret))
	assert.Equal(t, value, ret.Value["value"])
	assert.Equal(t, conditions, ret.Value["conditions"])
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
		}`}, "ApplicationId": {"1"},
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

func TestJSON(t *testing.T) {
	assert.NoError(t, keyLogin(1))

	contract := randName("JSONEncode")
	assert.NoError(t, postTx("NewContract", &url.Values{
		"Value": {`contract ` + contract + ` {
			action {
				var a array, m map
				m["k1"] = 1
				m["k2"] = 2
				a[0] = m
				a[1] = m

				info JSONEncode(a)
			}
		}`}, "ApplicationId": {"1"},
		"Conditions": {"true"},
	}))
	assert.EqualError(t, postTx(contract, &url.Values{}), `{"type":"info","error":"[{\"k1\":1,\"k2\":2},{\"k1\":1,\"k2\":2}]"}`)

	contract = randName("JSONDecode")
	assert.NoError(t, postTx("NewContract", &url.Values{
		"Value": {`contract ` + contract + ` {
			data {
				Input string
			}
			action {
				info Sprintf("%#v", JSONDecode($Input))
			}
		}`}, "ApplicationId": {"1"},
		"Conditions": {"true"},
	}))

	cases := []struct {
		source string
		result string
	}{
		{`"test"`, `{"type":"info","error":"\"test\""}`},
		{`["test"]`, `{"type":"info","error":"[]interface {}{\"test\"}"}`},
		{`{"test":1}`, `{"type":"info","error":"map[string]interface {}{\"test\":1}"}`},
		{`[{"test":1}]`, `{"type":"info","error":"[]interface {}{map[string]interface {}{\"test\":1}}"}`},
		{`{"test":1`, `{"type":"panic","error":"unexpected end of JSON input"}`},
	}

	for _, v := range cases {
		assert.EqualError(t, postTx(contract, &url.Values{"Input": {v.source}}), v.result)
	}
}

func TestBytesToString(t *testing.T) {
	assert.NoError(t, keyLogin(1))

	contract := randName("BytesToString")
	assert.NoError(t, postTx("NewContract", &url.Values{
		"Value": {`contract ` + contract + ` {
			data {
				File bytes "file"
			}
			action {
				$result = BytesToString($File)
			}
		}`},
		"Conditions":    {"true"},
		"ApplicationId": {"1"},
	}))

	content := crypto.RandSeq(100)
	_, res, err := postTxMultipart(contract, nil, map[string][]byte{"File": []byte(content)})
	assert.NoError(t, err)
	assert.Equal(t, content, res)
}

func TestMoneyDigits(t *testing.T) {
	assert.NoError(t, keyLogin(1))

	var v paramValue
	assert.NoError(t, sendGet("/ecosystemparam/money_digit", &url.Values{}, &v))

	contract := randName("MoneyDigits")
	assert.NoError(t, postTx("NewContract", &url.Values{
		"Value": {`contract ` + contract + ` {
			data {
				Value money
			}
			action {
				$result = $Value
			}
		}`},
		"ApplicationId": {"1"},
		"Conditions":    {"true"},
	}))

	_, result, err := postTxResult(contract, &url.Values{
		"Value": {"1"},
	})
	assert.NoError(t, err)

	d := decimal.New(1, int32(converter.StrToInt(v.Value)))
	assert.Equal(t, d.StringFixed(0), result)
}

func TestMemoryLimit(t *testing.T) {
	assert.NoError(t, keyLogin(1))

	contract := randName("Contract")
	assert.NoError(t, postTx("NewContract", &url.Values{
		"Value": {`contract ` + contract + ` {
			data {
				Count int "optional"
			}
			action {
				var a array
				while (true) {
					$Count = $Count + 1
					a[Len(a)] = JSONEncode(a)
				}
			}
		}`},
		"ApplicationId": {"1"},
		"Conditions":    {"true"},
	}))

	assert.EqualError(t, postTx(contract, &url.Values{}), `{"type":"panic","error":"Memory limit exceeded"}`)
}

func TestStack(t *testing.T) {
	assert.NoError(t, keyLogin(1))

	parent := randName("Parent")
	child := randName("Child")

	assert.NoError(t, postTx("NewContract", &url.Values{
		"Value": {`contract ` + child + ` {
			action {
				$result = $stack
			}
		}`},
		"ApplicationId": {"1"},
		"Conditions":    {"true"},
	}))

	assert.NoError(t, postTx("NewContract", &url.Values{
		"Value": {`contract ` + parent + ` {
			action {
				var arr array
				arr[0] = $stack
				arr[1] = ` + child + `()
				$result = arr
			}
		}`},
		"ApplicationId": {"1"},
		"Conditions":    {"true"},
	}))

	_, res, err := postTxResult(parent, &url.Values{})
	assert.NoError(t, err)
	assert.Equal(t, fmt.Sprintf("[[@1%s] [@1%[1]s @1%s]]", parent, child), res)
}

func TestPageHistory(t *testing.T) {
	assert.NoError(t, keyLogin(1))

	name := randName(`page`)
	value := `P(test,test paragraph)`

	form := url.Values{"Name": {name}, "Value": {value}, "ApplicationId": {`1`},
		"Menu": {"default_menu"}, "Conditions": {"ContractConditions(`MainCondition`)"}}
	assert.NoError(t, postTx(`NewPage`, &form))

	var ret listResult
	assert.NoError(t, sendGet(`list/pages`, nil, &ret))
	id := ret.Count
	assert.NoError(t, postTx(`EditPage`, &url.Values{"Id": {id}, "Value": {"Div(style){ok}"}}))
	assert.NoError(t, postTx(`EditPage`, &url.Values{"Id": {id}, "Conditions": {"true"}}))

	form = url.Values{"Name": {randName(`menu`)}, "Value": {`MenuItem(First)MenuItem(Second)`},
		"ApplicationId": {`1`}, "Conditions": {"ContractConditions(`MainCondition`)"}}
	assert.NoError(t, postTx(`NewMenu`, &form))

	assert.NoError(t, sendGet(`list/menu`, nil, &ret))
	idmenu := ret.Count
	assert.NoError(t, postTx(`EditMenu`, &url.Values{"Id": {idmenu}, "Conditions": {"true"}}))
	assert.NoError(t, postTx(`EditMenu`, &url.Values{"Id": {idmenu}, "Value": {"MenuItem(Third)"}}))
	assert.NoError(t, postTx(`EditMenu`, &url.Values{"Id": {idmenu},
		"Value": {"MenuItem(Third)"}, "Conditions": {"false"}}))

	form = url.Values{"Value": {`contract C` + name + `{ action {}}`},
		"ApplicationId": {`1`}, "Conditions": {"ContractConditions(`MainCondition`)"}}
	_, idCont, err := postTxResult(`NewContract`, &form)
	assert.NoError(t, err)
	assert.NoError(t, postTx(`EditContract`, &url.Values{"Id": {idCont},
		"Value": {`contract C` + name + `{ action {Println("OK")}}`}, "Conditions": {"true"}}))

	form = url.Values{`Value`: {`contract Get` + name + ` {
		data {
			IdPage int
			IdMenu int
			IdCont int
		}
		action {
			var ret array
			ret = GetPageHistory($IdPage)
			$result = Str(Len(ret))
			ret = GetMenuHistory($IdMenu)
			$result = $result + Str(Len(ret))
			ret = GetContractHistory($IdCont)
			$result = $result + Str(Len(ret))
		}
	}`}, "ApplicationId": {`1`}, `Conditions`: {`true`}}
	assert.NoError(t, postTx(`NewContract`, &form))

	form = url.Values{`Value`: {`contract GetRow` + name + ` {
		data {
			IdPage int
		}
		action {
			var ret array
			var row got map
			ret = GetPageHistory($IdPage)
			row = ret[1]
			got = GetPageHistoryRow($IdPage, Int(row["id"]))
			if got["block_id"] != row["block_id"] {
				error "GetPageHistory"
			}
		}
	}`}, "ApplicationId": {`1`}, `Conditions`: {`true`}}
	assert.NoError(t, postTx(`NewContract`, &form))

	_, msg, err := postTxResult(`Get`+name, &url.Values{"IdPage": {id}, "IdMenu": {idmenu},
		"IdCont": {idCont}})
	assert.NoError(t, err)
	assert.Equal(t, `231`, msg)

	form = url.Values{"Name": {name + `1`}, "Value": {value}, "ApplicationId": {`1`},
		"Menu": {"default_menu"}, "Conditions": {"ContractConditions(`MainCondition`)"}}
	assert.NoError(t, postTx(`NewPage`, &form))

	assert.NoError(t, postTx(`Get`+name, &url.Values{"IdPage": {converter.Int64ToStr(
		converter.StrToInt64(id) + 1)}, "IdMenu": {idmenu}, "IdCont": {idCont}}))

	assert.EqualError(t, postTx(`Get`+name, &url.Values{"IdPage": {`1000000`}, "IdMenu": {idmenu},
		"IdCont": {idCont}}), `{"type":"panic","error":"Record has not been found"}`)

	assert.NoError(t, postTx(`GetRow`+name, &url.Values{"IdPage": {id}}))

	var retTemp contentResult
	assert.NoError(t, sendPost(`content`, &url.Values{`template`: {fmt.Sprintf(`GetPageHistory(MySrc,%s)`,
		id)}}, &retTemp))

	if len(RawToString(retTemp.Tree)) < 400 {
		t.Error(fmt.Errorf(`wrong tree %s`, RawToString(retTemp.Tree)))
	}
}
