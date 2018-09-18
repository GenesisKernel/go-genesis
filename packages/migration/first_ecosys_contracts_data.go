// Code generated by go generate; DO NOT EDIT.

package migration

var firstEcosystemContractsSQL = `
INSERT INTO "1_contracts" (id, name, value, conditions, app_id, wallet_id)
VALUES
	(next_id('1_contracts'), 'ActivateContract', 'contract ActivateContract {
	data {
		Id  int
	}
	conditions {
		$cur = DBRow("contracts").Columns("id,conditions,active,wallet_id").WhereId($Id)
		if !$cur {
			error Sprintf("Contract %%d does not exist", $Id)
		}
		if Int($cur["active"]) == 1 {
			error Sprintf("The contract %%d has been already activated", $Id)
		}
		Eval($cur["conditions"])
		if $key_id != Int($cur["wallet_id"]) {
			error Sprintf("Wallet %%d cannot activate the contract", $key_id)
		}
	}
	action {
		DBUpdate("contracts", $Id, {"active": 1})
		Activate($Id, $ecosystem_id)
	}
}
', 'ContractConditions("MainCondition")', 1, %[1]d),
	(next_id('1_contracts'), 'CallDelayedContract', 'contract CallDelayedContract {
	data {
		Id int
	}
	conditions {
		var rows array
		rows = DBFind("delayed_contracts").Where({id: $Id, deleted: "false"} )

		if !Len(rows) {
			error Sprintf("Delayed contract %%d does not exist", $Id)
		}
		$cur = rows[0]

		if $key_id != Int($cur["key_id"]) {
			error "Access denied"
		}

		if $block < Int($cur["block_id"]) {
			error Sprintf("Delayed contract %%d must run on block %%s, current block %%d", $Id, $cur["block_id"], $block)
		}
	}
	action {
		var limit, counter, block_id int

		limit = Int($cur["limit"])
		counter = Int($cur["counter"])+1
		block_id = $block

		if limit == 0 || limit > counter {
			block_id = block_id + Int($cur["every_block"])
		}
		DBUpdate("delayed_contracts", $Id, {"counter": counter, "block_id": block_id})

		var params map
		CallContract($cur["contract"], params)
	}
}
', 'ContractConditions("MainCondition")', 1, %[1]d),
	(next_id('1_contracts'), 'CheckNodesBan', 'contract CheckNodesBan {
	action {
		UpdateNodesBan($block_time)
	}
}
', 'ContractConditions("MainCondition")', 1, %[1]d),
	(next_id('1_contracts'), 'DeactivateContract', 'contract DeactivateContract {
	data {
		Id         int
	}
	conditions {
		$cur = DBRow("contracts").Columns("id,conditions,active,wallet_id").WhereId($Id)
		if !$cur {
			error Sprintf("Contract %%d does not exist", $Id)
		}
		if Int($cur["active"]) == 0 {
			error Sprintf("The contract %%d has been already deactivated", $Id)
		}
		Eval($cur["conditions"])
		if $key_id != Int($cur["wallet_id"]) {
			error Sprintf("Wallet %%d cannot deactivate the contract", $key_id)
		}
	}
	action {
		DBUpdate("contracts", $Id, {"active": 0})
		Deactivate($Id, $ecosystem_id)
	}
}
', 'ContractConditions("MainCondition")', 1, %[1]d),
	(next_id('1_contracts'), 'EditAppParam', 'contract EditAppParam {
    data {
        Id int
        Value string "optional"
        Conditions string "optional"
    }
    func onlyConditions() bool {
        return $Conditions && !$Value
    }

    conditions {
        RowConditions("app_params", $Id, onlyConditions())
        if $Conditions {
            ValidateCondition($Conditions, $ecosystem_id)
        }
    }

    action {
        var pars map
        if $Value {
            pars["value"] = $Value
        }
        if $Conditions {
            pars["conditions"] = $Conditions
        }
        if pars {
            DBUpdate("app_params", $Id, pars)
        }
    }
}
', 'ContractConditions("MainCondition")', 1, %[1]d),
	(next_id('1_contracts'), 'EditApplication', 'contract EditApplication {
    data {
        ApplicationId int
        Conditions string "optional"
    }
    func onlyConditions() bool {
        return $Conditions && false
    }

    conditions {
        RowConditions("applications", $ApplicationId, onlyConditions())
        if $Conditions {
            ValidateCondition($Conditions, $ecosystem_id)
        }
    }

    action {
        var pars map
        if $Conditions {
            pars["conditions"] = $Conditions
        }
        if pars {
            DBUpdate("applications", $ApplicationId, pars)
        }
    }
}
', 'ContractConditions("MainCondition")', 1, %[1]d),
	(next_id('1_contracts'), 'EditBlock', 'contract EditBlock {
    data {
        Id int
        Value string "optional"
        Conditions string "optional"
    }
    func onlyConditions() bool {
        return $Conditions && !$Value
    }

    conditions {
        RowConditions("blocks", $Id, onlyConditions())
        if $Conditions {
            ValidateCondition($Conditions, $ecosystem_id)
        }
    }

    action {
        var pars map
        if $Value {
            pars["value"] = $Value
        }
        if $Conditions {
            pars["conditions"] = $Conditions
        }
        if pars {
            DBUpdate("blocks", $Id, pars)
        }
    }
}
', 'ContractConditions("MainCondition")', 1, %[1]d),
	(next_id('1_contracts'), 'EditColumn', 'contract EditColumn {
    data {
        TableName string
        Name string
        Permissions string
    }

    conditions {
        ColumnCondition($TableName, $Name, "", $Permissions)
    }

    action {
        PermColumn($TableName, $Name, $Permissions)
    }
}
', 'ContractConditions("MainCondition")', 1, %[1]d),
	(next_id('1_contracts'), 'EditContract', 'contract EditContract {
    data {
        Id int
        Value string "optional"
        Conditions string "optional"
        WalletId string "optional"
    }
    func onlyConditions() bool {
        return $Conditions && !$Value && !$WalletId
    }

    conditions {
        RowConditions("contracts", $Id, onlyConditions())
        if $Conditions {
            ValidateCondition($Conditions, $ecosystem_id)
        }
        $cur = DBFind("contracts").Columns("id,value,conditions,active,wallet_id,token_id").WhereId($Id).Row()
        if !$cur {
            error Sprintf("Contract %%d does not exist", $Id)
        }
        if $Value {
            ValidateEditContractNewValue($Value, $cur["value"])
        }
        if $WalletId != "" {
            $recipient = AddressToId($WalletId)
            if $recipient == 0 {
                error Sprintf("New contract owner %%s is invalid", $WalletId)
            }
            if Int($cur["active"]) == 1 {
                error "Contract must be deactivated before wallet changing"
            }
        } else {
            $recipient = Int($cur["wallet_id"])
        }
    }

    action {
        UpdateContract($Id, $Value, $Conditions, $WalletId, $recipient, $cur["active"], $cur["token_id"])
    }
}
', 'ContractConditions("MainCondition")', 1, %[1]d),
	(next_id('1_contracts'), 'EditLang', 'contract EditLang {
    data {
        Id int
        Trans string
    }

    conditions {
        EvalCondition("parameters", "changing_language", "value")
        $lang = DBFind("languages").Where({id: $Id}).Row()
    }

    action {
        EditLanguage($Id, $lang["name"], $Trans, Int($lang["app_id"]))
    }
}
', 'ContractConditions("MainCondition")', 1, %[1]d),
	(next_id('1_contracts'), 'EditMenu', 'contract EditMenu {
    data {
        Id int
        Value string "optional"
        Title string "optional"
        Conditions string "optional"
    }
    func onlyConditions() bool {
        return $Conditions && !$Value && !$Title
    }

    conditions {
        RowConditions("menu", $Id, onlyConditions())
        if $Conditions {
            ValidateCondition($Conditions, $ecosystem_id)
        }
    }

    action {
        var pars map
        if $Value {
            pars["value"] = $Value
        }
        if $Title {
            pars["title"] = $Title
        }
        if $Conditions {
            pars["conditions"] = $Conditions
        }
        if pars {
            DBUpdate("menu", $Id, pars)
        }            
    }
}
', 'ContractConditions("MainCondition")', 1, %[1]d),
	(next_id('1_contracts'), 'EditPage', 'contract EditPage {
    data {
        Id int
        Value string "optional"
        Menu string "optional"
        Conditions string "optional"
        ValidateCount int "optional"
        ValidateMode string "optional"
    }
    func onlyConditions() bool {
        return $Conditions && !$Value && !$Menu && !$ValidateCount 
    }
    func preparePageValidateCount(count int) int {
        var min, max int
        min = Int(EcosysParam("min_page_validate_count"))
        max = Int(EcosysParam("max_page_validate_count"))
        if count < min {
            count = min
        } else {
            if count > max {
                count = max
            }
        }
        return count
    }

    conditions {
        RowConditions("pages", $Id, onlyConditions())
        if $Conditions {
            ValidateCondition($Conditions, $ecosystem_id)
        }
        $ValidateCount = preparePageValidateCount($ValidateCount)
    }

    action {
        var pars map
        if $Value {
            pars["value"] = $Value
        }
        if $Menu {
            pars["menu"] = $Menu
        }
        if $Conditions {
            pars["conditions"] = $Conditions
        }
        if $ValidateCount {
            pars["validate_count"] = $ValidateCount
        }
        if $ValidateMode {
            if $ValidateMode != "1" {
                $ValidateMode = "0"
            }
            pars["validate_mode"] = $ValidateMode
        }
        if pars {
            DBUpdate("pages", $Id, pars)
        }
    }
}
', 'ContractConditions("MainCondition")', 1, %[1]d),
	(next_id('1_contracts'), 'EditTable', 'contract EditTable {
    data {
        Name string
        InsertPerm string
        UpdatePerm string
        NewColumnPerm string
        ReadPerm string "optional"
    }

    conditions {
        if !$InsertPerm {
            info("Insert condition is empty")
        }
        if !$UpdatePerm {
            info("Update condition is empty")
        }
        if !$NewColumnPerm {
            info("New column condition is empty")
        }

        var permissions map
        permissions["insert"] = $InsertPerm
        permissions["update"] = $UpdatePerm
        permissions["new_column"] = $NewColumnPerm
        if $ReadPerm {
            permissions["read"] = $ReadPerm
        }
        $Permissions = permissions
        TableConditions($Name, "", JSONEncode($Permissions))
    }

    action {
        PermTable($Name, JSONEncode($Permissions))
    }
}
', 'ContractConditions("MainCondition")', 1, %[1]d),
	(next_id('1_contracts'), 'Import', 'contract Import {
    data {
        Data string
    }
    func ReplaceValue(s string) string {
        s = Replace(s, "#IMPORT_ECOSYSTEM_ID#", "#ecosystem_id#")
        s = Replace(s, "#IMPORT_KEY_ID#", "#key_id#")
        s = Replace(s, "#IMPORT_ISMOBILE#", "#isMobile#")
        s = Replace(s, "#IMPORT_ROLE_ID#", "#role_id#")
        s = Replace(s, "#IMPORT_ECOSYSTEM_NAME#", "#ecosystem_name#")
        s = Replace(s, "#IMPORT_APP_ID#", "#app_id#")
        return s
    }

    conditions {
        $Data = ReplaceValue($Data)

        $ApplicationId = 0
        var app_map map
        app_map = DBFind("buffer_data").Columns("value->app_name").Where({key: "import_info",
          member_id: $key_id}).Row()

        if app_map{
            var app_id int
            var ival string
            ival = Str(app_map["value.app_name"])
            app_id = DBFind("applications").Columns("id").Where({name: ival}).One("id")
            if app_id {
                $ApplicationId = Int(app_id)
            }
        }
    }

    action {
        var editors, creators map
        editors["pages"] = "EditPage"
        editors["blocks"] = "EditBlock"
        editors["menu"] = "EditMenu"
        editors["app_params"] = "EditAppParam"
        editors["languages"] = "EditLang"
        editors["contracts"] = "EditContract"
        editors["tables"] = "" // nothing

        creators["pages"] = "NewPage"
        creators["blocks"] = "NewBlock"
        creators["menu"] = "NewMenu"
        creators["app_params"] = "NewAppParam"
        creators["languages"] = "NewLang"
        creators["contracts"] = "NewContract"
        creators["tables"] = "NewTable"

        var dataImport array
        dataImport = JSONDecode($Data)
        var i int
        while i<Len(dataImport){
            var item, cdata map
            cdata = dataImport[i]
            if cdata {
                cdata["ApplicationId"] = $ApplicationId
                $Type = cdata["Type"]
                $Name = cdata["Name"]

                // Println(Sprintf("import %%v: %%v", $Type, cdata["Name"]))

                item = DBFind($Type).Where({name: $Name}).Row()
                var contractName string
                if item {
                    contractName = editors[$Type]
                    cdata["Id"] = Int(item["id"])
                    if $Type == "menu"{
                        var menu menuItem string
                        menu = Replace(item["value"], " ", "")
                        menu = Replace(menu, "\n", "")
                        menu = Replace(menu, "\r", "")
                        menuItem = Replace(cdata["Value"], " ", "")
                        menuItem = Replace(menuItem, "\n", "")
                        menuItem = Replace(menuItem, "\r", "")
                        if Contains(menu, menuItem) {
                            // ignore repeated
                            contractName = ""
                        }else{
                            cdata["Value"] = item["value"] + "\n" + cdata["Value"]
                        }
                    }
                } else {
                    contractName = creators[$Type]
                }

                if contractName != ""{
                    CallContract(contractName, cdata)
                }
            }
            i=i+1
        }
        // Println(Sprintf("> time: %%v", $time))
    }
}
', 'ContractConditions("MainCondition")', 1, %[1]d),
	(next_id('1_contracts'), 'ImportUpload', 'contract ImportUpload {
    data {
        input_file string "file"
    }
    func ReplaceValue(s string) string {
        s = Replace(s, "#ecosystem_id#", "#IMPORT_ECOSYSTEM_ID#")
        s = Replace(s, "#key_id#", "#IMPORT_KEY_ID#")
        s = Replace(s, "#isMobile#", "#IMPORT_ISMOBILE#")
        s = Replace(s, "#role_id#", "#IMPORT_ROLE_ID#")
        s = Replace(s, "#ecosystem_name#", "#IMPORT_ECOSYSTEM_NAME#")
        s = Replace(s, "#app_id#", "#IMPORT_APP_ID#")
        return s
    }

    conditions {
        $input_file = BytesToString($input_file)
        $input_file = ReplaceValue($input_file)
        $limit = 5 // data piece size of import

        // init buffer_data, cleaning old buffer
        var initJson map
        $import_id = DBFind("buffer_data").Where({member_id:$key_id, key: "import"}).One("id")
        if $import_id {
            $import_id = Int($import_id)
            DBUpdate("buffer_data", $import_id, {"value": initJson})
        } else {
            $import_id = DBInsert("buffer_data", {"member_id":$key_id,"key": "import",
                 "value": initJson})
        }

        $info_id = DBFind("buffer_data").Where({member_id:$key_id, key: "import_info"}).One("id")
        if $info_id {
            $info_id = Int($info_id)
            DBUpdate("buffer_data", $info_id, {"value": initJson})
        } else {
            $info_id = DBInsert("buffer_data", {"member_id":$key_id,"key": "import_info",
            "value": initJson})
        }
    }

    action {
        var input map
        input = JSONDecode($input_file)
        var arr_data array
        arr_data = input["data"]

        var pages_arr, blocks_arr, menu_arr, parameters_arr, languages_arr, contracts_arr, tables_arr array

        // import info
        var i int
        while i<Len(arr_data){
            var tmp_object map
            tmp_object = arr_data[i]

            if tmp_object["Type"] == "pages" {
                pages_arr = Append(pages_arr, Str(tmp_object["Name"]))
            }
            if tmp_object["Type"] == "blocks" {
                blocks_arr = Append(blocks_arr, Str(tmp_object["Name"]))
            }
            if tmp_object["Type"] == "menu" {
                menu_arr = Append(menu_arr, Str(tmp_object["Name"]))
            }
            if tmp_object["Type"] == "app_params" {
                parameters_arr = Append(parameters_arr, Str(tmp_object["Name"]))
            }
            if tmp_object["Type"] == "languages" {
                languages_arr = Append(languages_arr, Str(tmp_object["Name"]))
            }
            if tmp_object["Type"] == "contracts" {
                contracts_arr = Append(contracts_arr, Str(tmp_object["Name"]))
            }
            if tmp_object["Type"] == "tables" {
                tables_arr = Append(tables_arr, Str(tmp_object["Name"]))
            }

            i = i + 1
        }

        var info_map map
        info_map["app_name"] = input["name"]
        info_map["pages"] = Join(pages_arr, ", ")
        info_map["pages_count"] = Len(pages_arr)
        info_map["blocks"] = Join(blocks_arr, ", ")
        info_map["blocks_count"] = Len(blocks_arr)
        info_map["menu"] = Join(menu_arr, ", ")
        info_map["menu_count"] = Len(menu_arr)
        info_map["parameters"] = Join(parameters_arr, ", ")
        info_map["parameters_count"] = Len(parameters_arr)
        info_map["languages"] = Join(languages_arr, ", ")
        info_map["languages_count"] = Len(languages_arr)
        info_map["contracts"] = Join(contracts_arr, ", ")
        info_map["contracts_count"] = Len(contracts_arr)
        info_map["tables"] = Join(tables_arr, ", ")
        info_map["tables_count"] = Len(tables_arr)

        if 0 == Len(pages_arr) + Len(blocks_arr) + Len(menu_arr) + Len(parameters_arr) + Len(languages_arr) + Len(contracts_arr) + Len(tables_arr) {
            warning "Invalid or empty import file"
        }

        // import data
        // the contracts is imported in one piece, the rest is cut under the $limit, a crutch to bypass the error when you import dependent contracts in different pieces
        i=0
        var sliced contracts array, arr_data_len int
        arr_data_len = Len(arr_data)
        while i <arr_data_len{
            var part array, l int, tmp map
            while l < $limit && (i+l < arr_data_len) {
                tmp = arr_data[i+l]
                if tmp["Type"] == "contracts" {
                    contracts = Append(contracts, tmp)
                }else{
                    part = Append(part, tmp)
                }
                l=l+1
            }
            var batch map
            batch["Data"] = JSONEncode(part)
            sliced = Append(sliced, batch)
            i=i+$limit
        }
        if Len(contracts) > 0{
            var batch map
            batch["Data"] = JSONEncode(contracts)
            sliced = Append(sliced, batch)
        }
        input["data"] = sliced

        // storing
        DBUpdate("buffer_data", $import_id, {"value": input})
        DBUpdate("buffer_data", $info_id, {"value": info_map})

        var app_id int
        var ival string
        ival =  Str(input["name"])
        app_id = DBFind("applications").Columns("id").Where({name:ival}).One("id")

        if !app_id {
            var val string
            val = Str(input["name"])
            DBInsert("applications", {"name": val, "conditions": "true"})
        }
    }
}
', 'ContractConditions("MainCondition")', 1, %[1]d),
	(next_id('1_contracts'), 'MoneyTransfer', 'contract MoneyTransfer {
	data {
		Recipient string
		Amount    string
		Comment     string "optional"
	}
	conditions {
		$recipient = AddressToId($Recipient)
		if $recipient == 0 {
			error Sprintf("Recipient %%s is invalid", $Recipient)
		}
		var total money
		$amount = Money($Amount) 
		if $amount <= 0 {
			error "Amount must be greater then zero"
		}

        var key map
        var req money
		key = GetKey($key_id)
        total = Money($key["amount"])
        req = $amount + Money(100000000000000000) 
        if req > total {
			error Sprintf("Money is not enough. You have got %%v but you should reserve %%v", total, req)
		}
	}
	action {
		EditKey($key_id, $key["amount"]-$amount)

		var recipientKey map
		$recipientKey = GetKey($recipient)
		if $recipientKey == nil {
			CreateKey($recipient, $amount, "")
		} else {
			UpdateKey($recipient, $recipientKey["amount"]+$amount)
		}

        DBInsert("history", {sender_id: $key_id,recipient_id: $recipient,
             amount:$amount,comment: $Comment,block_id: $block,txhash: $txhash})
	}
}
', 'ContractConditions("MainCondition")', 1, %[1]d),
	(next_id('1_contracts'), 'NewAppParam', 'contract NewAppParam {
    data {
        ApplicationId int
        Name string
        Value string
        Conditions string
    }

    conditions {
        ValidateCondition($Conditions, $ecosystem_id)

        if $ApplicationId == 0 {
            warning "Application id cannot equal 0"
        }

        if DBFind("app_params").Columns("id").Where({"name":$Name}).One("id") {
            warning Sprintf( "Application parameter %%s already exists", $Name)
        }
    }

    action {
        DBInsert("app_params", {app_id: $ApplicationId, name: $Name, value: $Value,
              conditions: $Conditions})
    }
}
', 'ContractConditions("MainCondition")', 1, %[1]d),
	(next_id('1_contracts'), 'NewApplication', 'contract NewApplication {
    data {
        Name string
        Conditions string
    }

    conditions {
        ValidateCondition($Conditions, $ecosystem_id)

        if Size($Name) == 0 {
            warning "Application name missing"
        }

        if DBFind("applications").Columns("id").Where({name:$Name}).One("id") {
            warning Sprintf( "Application %%s already exists", $Name)
        }
    }

    action {
        $result = DBInsert("applications", {name: $Name,conditions: $Conditions})
    }
}
', 'ContractConditions("MainCondition")', 1, %[1]d),
	(next_id('1_contracts'), 'NewBadBlock', 'contract NewBadBlock {
	data {
		ProducerNodeID int
		ConsumerNodeID int
		BlockID int
		Timestamp int
		Reason string
	}
	action {
        DBInsert("@1_bad_blocks", {producer_node_id: $ProducerNodeID,consumer_node_id: $ConsumerNodeID,
            block_id: $BlockID, "timestamp block_time": $Timestamp, reason: $Reason})
	}
}
', 'ContractConditions("MainCondition")', 1, %[1]d),
	(next_id('1_contracts'), 'NewBlock', 'contract NewBlock {
    data {
        ApplicationId int
        Name string
        Value string
        Conditions string
    }

    conditions {
        ValidateCondition($Conditions, $ecosystem_id)

        if $ApplicationId == 0 {
            warning "Application id cannot equal 0"
        }

        if DBFind("blocks").Columns("id").Where({name:$Name}).One("id") {
            warning Sprintf( "Block %%s already exists", $Name)
        }
    }

    action {
        DBInsert("blocks", {name: $Name, value: $Value, conditions: $Conditions,
              app_id: $ApplicationId})
    }
}
', 'ContractConditions("MainCondition")', 1, %[1]d),
	(next_id('1_contracts'), 'NewContract', 'contract NewContract {
    data {
        ApplicationId int
        Value string
        Conditions string
        Wallet string "optional"
        TokenEcosystem int "optional"
    }

    conditions {
        ValidateCondition($Conditions,$ecosystem_id)

        if $ApplicationId == 0 {
            warning "Application id cannot equal 0"
        }

        $walletContract = $key_id
        if $Wallet {
            $walletContract = AddressToId($Wallet)
            if $walletContract == 0 {
                error Sprintf("wrong wallet %%s", $Wallet)
            }
        }

        $contract_name = ContractName($Value)

        if !$contract_name {
            error "must be the name"
        }

        if !$TokenEcosystem {
            $TokenEcosystem = 1
        } else {
            if !SysFuel($TokenEcosystem) {
                warning Sprintf("Ecosystem %%d is not system", $TokenEcosystem)
            }
        }
    }

    action {
        $result = CreateContract($contract_name, $Value, $Conditions, $walletContract, $TokenEcosystem, $ApplicationId)
    }
    func price() int {
        return SysParamInt("contract_price")
    }
}
', 'ContractConditions("MainCondition")', 1, %[1]d),
	(next_id('1_contracts'), 'NewLang', 'contract NewLang {
    data {
        ApplicationId int
        Name string
        Trans string
    }

    conditions {
        if $ApplicationId == 0 {
            warning "Application id cannot equal 0"
        }

        if DBFind("languages").Columns("id").Where({name: $Name}).One("id") {
            warning Sprintf( "Language resource %%s already exists", $Name)
        }

        EvalCondition("parameters", "changing_language", "value")
    }

    action {
        CreateLanguage($Name, $Trans, $ApplicationId)
    }
}
', 'ContractConditions("MainCondition")', 1, %[1]d),
	(next_id('1_contracts'), 'NewMenu', 'contract NewMenu {
    data {
        Name string
        Value string
        Title string "optional"
        Conditions string
    }

    conditions {
        ValidateCondition($Conditions,$ecosystem_id)

        if DBFind("menu").Columns("id").Where({name: $Name}).One("id") {
            warning Sprintf( "Menu %%s already exists", $Name)
        }
    }

    action {
        DBInsert("menu", {name:$Name,value: $Value, title: $Title, conditions: $Conditions})
    }
    func price() int {
        return SysParamInt("menu_price")
    }
}
', 'ContractConditions("MainCondition")', 1, %[1]d),
	(next_id('1_contracts'), 'NewPage', 'contract NewPage {
    data {
        ApplicationId int
        Name string
        Value string
        Menu string
        Conditions string
        ValidateCount int "optional"
        ValidateMode string "optional"
    }
    func preparePageValidateCount(count int) int {
        var min, max int
        min = Int(EcosysParam("min_page_validate_count"))
        max = Int(EcosysParam("max_page_validate_count"))

        if count < min {
            count = min
        } else {
            if count > max {
                count = max
            }
        }
        return count
    }

    conditions {
        ValidateCondition($Conditions,$ecosystem_id)

        if $ApplicationId == 0 {
            warning "Application id cannot equal 0"
        }

        if DBFind("pages").Columns("id").Where({name: $Name}).One("id") {
            warning Sprintf( "Page %%s already exists", $Name)
        }

        $ValidateCount = preparePageValidateCount($ValidateCount)

        if $ValidateMode {
            if $ValidateMode != "1" {
                $ValidateMode = "0"
            }
        }
    }

    action {
        DBInsert("pages", {name: $Name,value: $Value, menu: $Menu,
             validate_count:$ValidateCount,validate_mode: $ValidateMode,
             conditions: $Conditions,app_id: $ApplicationId})
    }
    func price() int {
        return SysParamInt("page_price")
    }
}
', 'ContractConditions("MainCondition")', 1, %[1]d),
	(next_id('1_contracts'), 'NewTable', 'contract NewTable {
    data {
        ApplicationId int
        Name string
        Columns string
        Permissions string
    }
    conditions {
        if $ApplicationId == 0 {
            warning "Application id cannot equal 0"
        }
        TableConditions($Name, $Columns, $Permissions)
    }
    
    action {
        CreateTable($Name, $Columns, $Permissions, $ApplicationId)
    }
    func price() int {
        return SysParamInt("table_price")
    }
}
', 'ContractConditions("MainCondition")', 1, %[1]d),
	(next_id('1_contracts'), 'NewUser', 'contract NewUser {
	data {
		NewPubkey string
	}
	conditions {
		$newId = PubToID($NewPubkey)
		if $newId == 0 {
			error "Wrong pubkey"
		}
		if DBFind("keys").Columns("id").WhereId($newId).One("id") != nil {
			error "User already exists"
		}

        $amount = Money(1000) * Money(1000000000000000000)
	}
	action {
        NewMoney($newId, Str($amount), "New user deposit")
        SetPubKey($newId, StringToBytes($NewPubkey))
	}
}
', 'ContractConditions("NodeOwnerCondition")', 1, %[1]d),
	(next_id('1_contracts'), 'NodeOwnerCondition', 'contract NodeOwnerCondition {
	conditions {
        $raw_full_nodes = SysParamString("full_nodes")
        if Size($raw_full_nodes) == 0 {
            ContractConditions("MainCondition")
        } else {
            $full_nodes = JSONDecode($raw_full_nodes)
            var i int
            while i < Len($full_nodes) {
                $fn = $full_nodes[i]
                if $fn["key_id"] == $key_id {
                    return true
                }
                i = i + 1
            }
            warning "Sorry, you do not have access to this action."
        }
	}
}
', 'ContractConditions("MainCondition")', 1, %[1]d),
	(next_id('1_contracts'), 'UpdateMetrics', 'contract UpdateMetrics {
	conditions {
		ContractConditions("MainCondition")
	}
	action {
		var values array
		values = DBCollectMetrics()

		var i, id int
		var v map
		while (i < Len(values)) {
            var inmap map

            v = values[i]
            inmap["time"] = v["time"]
            inmap["key"] = v["key"]
            inmap["metric"] = v["metric"]
            
            id = Int(DBFind("metrics").Columns("id").Where(inmap).One("id"))
            var ival int
			if id != 0 {
                ival = Int(v["value"])
				DBUpdate("metrics", id, {"value": ival})
			} else {
                inmap["value"] = Int(v["value"])
				DBInsert("metrics", inmap )
			}
			i = i + 1
		}
	}
}
', 'ContractConditions("MainCondition")', 1, %[1]d),
	(next_id('1_contracts'), 'UploadBinary', 'contract UploadBinary {
    data {
        ApplicationId int
        Name string
        Data bytes "file"
        DataMimeType string "optional"
    }

    conditions {
        $Id = Int(DBFind("binaries").Columns("id").Where({app_id: $ApplicationId,
            member_id: $key_id, name: $Name}).One("id"))

        if $Id == 0 {
            if $ApplicationId == 0 {
                warning "Application id cannot equal 0"
            }
        }
    }
    action {
        var hash string
        hash = Hash($Data)

        if $DataMimeType == "" {
            $DataMimeType = "application/octet-stream"
        }

        if $Id != 0 {
            DBUpdate("binaries", $Id, {"data": $Data,"hash": hash,"mime_type": $DataMimeType})
        } else {
            $Id = DBInsert("binaries", {"app_id": $ApplicationId,"member_id": $key_id,
               "name": $Name,"data": $Data,"hash": hash, "mime_type": $DataMimeType})
        }

        $result = $Id
    }
}
', 'ContractConditions("MainCondition")', 1, %[1]d);
`
