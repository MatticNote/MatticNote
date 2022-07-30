package main

import "embed"

//go:embed template/* template/**/*
var template embed.FS

//go:embed client/dist/client/*
var webCli embed.FS

//go:embed client/dist/ui/*.css client/dist/ui/*.js
var webUi embed.FS

//go:embed client/dist/ui/manifest.json
var webUiManifest []byte
