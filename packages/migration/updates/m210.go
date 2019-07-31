// Copyright (C) 2017, 2018, 2019 EGAAS S.A.
//
// This program is free software; you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation; either version 2 of the License, or (at
// your option) any later version.
//
// This program is distributed in the hope that it will be useful, but
// WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the GNU
// General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program; if not, write to the Free Software
// Foundation, Inc., 51 Franklin Street, Fifth Floor, Boston, MA 02110-1301, USA.

package updates

var M210 = `
    insert into "1_system_parameters" (id, name, value, conditions) values
 (next_id('1_system_parameters'), 'external_blockchain', '', 'ContractAccess("@1UpdateSysParam")');

    DROP TABLE IF EXISTS "external_blockchain";
	CREATE TABLE "external_blockchain" (
	"id" bigint NOT NULL DEFAULT '0',
	"netname" varchar(255)  NOT NULL DEFAULT '',
	"value" text NOT NULL DEFAULT ''
	);
	ALTER TABLE ONLY "external_blockchain" ADD CONSTRAINT "external_blockchain_pkey" PRIMARY KEY (id);
	CREATE INDEX "external_blockchain_index_name" ON "external_blockchain" (netname);
	
`
