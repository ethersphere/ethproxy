// Copyright 2021 The Swarm Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ethproxy

var CommitTime string

var (
	version = "0.0.4" // manually set semantic version number
	commit  string    // automatically set git commit hash

	Version = func() string {
		if commit != "" {
			return version + "-" + commit
		}
		return version + "-dev"
	}()
)
