package cache

import "errors"

var IDExists = errors.New("ID already exists")
var IDNotExists = errors.New("no such ID exists")
