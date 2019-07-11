// Apla Software includes an integrated development
// environment with a multi-level system for the management
// of access rights to data, interfaces, and Smart contracts. The
// technical characteristics of the Apla Software are indicated in
// Apla Technical Paper.
//
// Apla Users are granted a permission to deal in the Apla
// Software without restrictions, including without limitation the
// rights to use, copy, modify, merge, publish, distribute, sublicense,
// and/or sell copies of Apla Software, and to permit persons
// to whom Apla Software is furnished to do so, subject to the
// following conditions:
// * the copyright notice of GenesisKernel and EGAAS S.A.
// and this permission notice shall be included in all copies or
// substantial portions of the software;
// * a result of the dealing in Apla Software cannot be
// implemented outside of the Apla Platform environment.
//
// THE APLA SOFTWARE IS PROVIDED “AS IS”, WITHOUT WARRANTY
// OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED
// TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A
// PARTICULAR PURPOSE, ERROR FREE AND NONINFRINGEMENT. IN
// NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE
// LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
// IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
// FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR
// THE USE OR OTHER DEALINGS IN THE APLA SOFTWARE.

package obs

var tablesDataSQL = `INSERT INTO "1_tables" ("id", "name", "permissions","columns", "conditions") VALUES 
(next_id('1_tables'), 'contracts', '{"insert": "ContractConditions(\"MainCondition\")", "update": "ContractConditions(\"MainCondition\")", "new_column": "ContractConditions(\"MainCondition\")"}', 
'{"name": "false", 
	"value": "ContractConditions(\"MainCondition\")",
	  "wallet_id": "ContractConditions(\"MainCondition\")",
	  "token_id": "ContractConditions(\"MainCondition\")",
	  "conditions": "ContractConditions(\"MainCondition\")"}', 'ContractAccess("@1EditTable")'),
	(next_id('1_tables'), 'keys', 
	'{"insert": "true", "update": "true", 
	  "new_column": "ContractConditions(\"MainCondition\")"}',
	'{
		"pub": "ContractConditions(\"MainCondition\")",
		"amount": "ContractConditions(\"MainCondition\")",
		"deleted": "ContractConditions(\"MainCondition\")",
		"blocked": "ContractConditions(\"MainCondition\")",
		"multi": "ContractConditions(\"MainCondition\")",
		"account": "false",
		"ecosystem": "false",
		"multi": "ContractConditions(\"@1AdminCondition\")"
	}', 
	'ContractAccess("@1EditTable")'),
	(next_id('1_tables'), 'history', 
	'{"insert": "ContractConditions(\"NodeOwnerCondition\")", "update": "ContractConditions(\"MainCondition\")", 
	  "new_column": "ContractConditions(\"MainCondition\")"}',
	'{"sender_id": "ContractConditions(\"MainCondition\")",
	  "recipient_id": "ContractConditions(\"MainCondition\")",
	  "amount":  "ContractConditions(\"MainCondition\")",
	  "comment": "ContractConditions(\"MainCondition\")",
	  "block_id":  "ContractConditions(\"MainCondition\")",
		"txhash": "ContractConditions(\"MainCondition\")",
		"created_at": "false"}', 'ContractAccess("@1EditTable")'),        
	(next_id('1_tables'), 'languages', 
	'{"insert": "ContractConditions(\"MainCondition\")", "update": "ContractConditions(\"MainCondition\")", 
	  "new_column": "ContractConditions(\"MainCondition\")"}',
	'{"app_id": "ContractConditions(\"MainCondition\")",
	  "name": "ContractConditions(\"MainCondition\")",
	  "res": "ContractConditions(\"MainCondition\")",
	  "conditions": "ContractConditions(\"MainCondition\")",
	  "app_id": "ContractConditions(\"MainCondition\")"}', 'ContractAccess("@1EditTable")'),
	(next_id('1_tables'), 'menu', 
		'{"insert": "ContractConditions(\"MainCondition\")", "update": "ContractConditions(\"MainCondition\")", 
	  "new_column": "ContractConditions(\"MainCondition\")"}',
	'{"name": "ContractConditions(\"MainCondition\")",
"value": "ContractConditions(\"MainCondition\")",
"conditions": "ContractConditions(\"MainCondition\")"
	}', 'ContractAccess("@1EditTable")'),
	(next_id('1_tables'), 'pages', 
		'{"insert": "ContractConditions(\"MainCondition\")", "update": "ContractConditions(\"MainCondition\")", 
	  "new_column": "ContractConditions(\"MainCondition\")"}',
	'{"name": "ContractConditions(\"MainCondition\")",
"value": "ContractConditions(\"MainCondition\")",
"menu": "ContractConditions(\"MainCondition\")",
"validate_count": "ContractConditions(\"MainCondition\")",
"validate_mode": "ContractConditions(\"MainCondition\")",
"app_id": "ContractConditions(\"MainCondition\")",
"conditions": "ContractConditions(\"MainCondition\")"
	}', 'ContractAccess("@1EditTable")'),
	(next_id('1_tables'), 'blocks', 
	'{"insert": "ContractConditions(\"MainCondition\")", "update": "ContractConditions(\"MainCondition\")", 
	  "new_column": "ContractConditions(\"MainCondition\")"}',
	'{"name": "ContractConditions(\"MainCondition\")",
"value": "ContractConditions(\"MainCondition\")",
"conditions": "ContractConditions(\"MainCondition\")"
	}', 'ContractAccess("@1EditTable")'),
	('8', 'signatures', 
	'{"insert": "ContractConditions(\"MainCondition\")", "update": "ContractConditions(\"MainCondition\")", 
	  "new_column": "ContractConditions(\"MainCondition\")"}',
	'{"name": "ContractConditions(\"MainCondition\")",
"value": "ContractConditions(\"MainCondition\")",
"conditions": "ContractConditions(\"MainCondition\")"
	}', 'ContractAccess("@1EditTable")'),
	('9', 'members', 
		'{"insert":"ContractAccess(\"Profile_Edit\")","update":"ContractConditions(\"MainCondition\")","new_column":"ContractConditions(\"MainCondition\")"}',
		'{
			"image_id":"ContractAccess(\"Profile_Edit\")",
			"member_info":"ContractAccess(\"Profile_Edit\")",
			"member_name":"false",
			"account":"false"
		}', 
		'ContractConditions("MainCondition")'),
	('10', 'roles',
		'{"insert":"ContractAccess(\"Roles_Create\")",
			"update":"ContractConditions(\"MainCondition\")",
			"new_column":"ContractConditions(\"MainCondition\")"}', 
		'{"default_page":"false",
			"creator":"false",
			"deleted":"ContractAccess(\"Roles_Del\")",
			"company_id":"false",
			"date_deleted":"ContractAccess(\"Roles_Del\")",
			"image_id":"ContractAccess(\"Roles_Create\")",
			"role_name":"false",
			"date_created":"false",
			"roles_access":"ContractAccess(\"Roles_AccessManager\")",
			"role_type":"false"}',
		'ContractConditions("MainCondition")'),
	('11', 'roles_participants',
		'{"insert":"ContractAccess(\"Roles_Assign\",\"voting_CheckDecision\")",
			"update":"ContractConditions(\"MainCondition\")",
			"new_column":"ContractConditions(\"MainCondition\")"}',
		'{"deleted":"ContractAccess(\"Roles_Unassign\")",
			"date_deleted":"ContractAccess(\"Roles_Unassign\")",
			"member":"false",
			"role":"false",
			"date_created":"false",
			"appointed":"false"}', 
		'ContractConditions("MainCondition")'),
	('12', 'notifications',
		'{"insert":"ContractAccess(\"notifications_Send\", \"CheckNodesBan\")",
			"update":"ContractAccess(\"notifications_Send\", \"notifications_Close\", \"notifications_Process\")",
			"new_column":"ContractConditions(\"MainCondition\")"}',
		'{"date_closed":"ContractAccess(\"notifications_Close\")",
			"sender":"false",
			"processing_info":"ContractAccess(\"notifications_Close\",\"notifications_Process\")",
			"date_start_processing":"ContractAccess(\"notifications_Close\",\"notifications_Process\")",
			"notification":"false",
			"page_name":"false",
			"page_params":"false",
			"closed":"ContractAccess(\"notifications_Close\")",
			"date_created":"false",
			"recipient":"false"}',
		'ContractAccess("@1EditTable")'),
	('13', 'sections', 
		'{"insert": "ContractConditions(\"MainCondition\")", "update": "ContractConditions(\"MainCondition\")", 
		"new_column": "ContractConditions(\"MainCondition\")"}',
		'{"title": "ContractConditions(\"MainCondition\")",
			"urlname": "ContractConditions(\"MainCondition\")",
			"page": "ContractConditions(\"MainCondition\")",
			"roles_access": "ContractConditions(\"MainCondition\")",
			"delete": "ContractConditions(\"MainCondition\")"}', 
			'ContractConditions("MainCondition")'),
	('14', 'applications',
		'{"insert": "ContractConditions(\"MainCondition\")",
			 "update": "ContractConditions(\"MainCondition\")", 
			 "new_column": "ContractConditions(\"MainCondition\")"}',
		'{"name": "ContractConditions(\"MainCondition\")",
		  "uuid": "false",
		  "conditions": "ContractConditions(\"MainCondition\")",
		  "deleted": "ContractConditions(\"MainCondition\")"}',
		'ContractConditions("MainCondition")'),
	('15', 'binaries',
		'{"insert":"ContractAccess(\"@1UploadBinary\")",
			"update":"ContractAccess(\"@1UploadBinary\")",
			"new_column":"ContractConditions(\"MainCondition\")"}',
		'{
			"hash":"ContractAccess(\"@1UploadBinary\")",
			"account": "false",
			"data":"ContractAccess(\"@1UploadBinary\")",
			"name":"false",
			"app_id":"false"
		}',
		'ContractConditions(\"MainCondition\")'),
	('16', 'parameters',
		'{"insert": "ContractConditions(\"MainCondition\")", "update": "ContractConditions(\"MainCondition\")",
			"new_column": "ContractConditions(\"MainCondition\")"}',
		'{"name": "ContractConditions(\"MainCondition\")",
			"value": "ContractConditions(\"MainCondition\")",
			"conditions": "ContractConditions(\"MainCondition\")"}',
		'ContractAccess("@1EditTable")'),
	('17', 'app_params',
		'{"insert": "ContractConditions(\"MainCondition\")", "update": "ContractConditions(\"MainCondition\")",
			"new_column": "ContractConditions(\"MainCondition\")"}',
		'{"app_id": "ContractConditions(\"MainCondition\")",
			"name": "ContractConditions(\"MainCondition\")",
			"value": "ContractConditions(\"MainCondition\")",
			"conditions": "ContractConditions(\"MainCondition\")"}',
		'ContractAccess("@1EditTable")'),
		('18', 'cron',
	  '{"insert": "ContractConditions(\"MainCondition\")", "update": "ContractConditions(\"MainCondition\")",
		"new_column": "ContractConditions(\"MainCondition\")"}',
	  '{"owner": "ContractConditions(\"MainCondition\")",
	  "cron": "ContractConditions(\"MainCondition\")",
	  "contract": "ContractConditions(\"MainCondition\")",
	  "counter": "ContractConditions(\"MainCondition\")",
	  "till": "ContractConditions(\"MainCondition\")",
		"conditions": "ContractConditions(\"MainCondition\")"
	  }', 'ContractConditions("MainCondition")'),
	('19', 'buffer_data',
		'{"insert":"true","update":"ContractConditions(\"MainCondition\")",
			"new_column":"ContractConditions(\"MainCondition\")"}',
		'{
			"key": "false",
			"value": "true",
			"account": "false"
		}',
		'ContractConditions("MainCondition")');
`
