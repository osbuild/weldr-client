// Copyright 2020 by Red Hat, Inc. All rights reserved.
// Use of this source is goverend by the Apache License
// that can be found in the LICENSE file.

/*
Package weldr contains functions used to interact with a WELDR API Server

For normal usage InitClientUnixSocket() should be called with the api version
and full path of the server's Unix Domain Socket. It will return a weldr.Client
struct that you can then use to interact with the server.

For testing you can initialize a temporary weldr.Client using weldr.NewClient(),
this is used in the weldr test functions.

*/
package weldr
