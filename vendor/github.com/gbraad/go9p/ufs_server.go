// Copyright 2009 The go9p Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package go9p

import (
	"log"
)

func StartServer(addrVal string, debugVal int, rootVal string) {
	ufs := new(Ufs)
	ufs.Dotu = true
	ufs.Id = "ufs"
	ufs.Root = rootVal
	ufs.Debuglevel = debugVal
	ufs.Start(ufs)

	// determined by build tags
	// extraFuncs()
	err := ufs.StartNetListener("tcp", addrVal)
	if err != nil {
		log.Println(err)
	}
}
