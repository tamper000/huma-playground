package handler

import "github.com/danielgtaylor/huma/v2"

var unknownError = huma.Error500InternalServerError("Unknown error")
var IDNotExists = huma.Error404NotFound("no such ID exists")
