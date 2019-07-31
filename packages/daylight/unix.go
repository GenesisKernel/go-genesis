// +build linux freebsd darwin
// +build 386 amd64

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

package daylight

import (
	"syscall"

	"github.com/AplaProject/go-apla/packages/converter"

	log "github.com/sirupsen/logrus"
)

// KillPid is killing process by PID
func KillPid(pid string) error {
	err := syscall.Kill(converter.StrToInt(pid), syscall.SIGTERM)
	if err != nil {
		log.WithFields(log.Fields{"pid": pid, "signal": syscall.SIGTERM}).Error("Error killing process with pid")
		return err
	}
	return nil
}
