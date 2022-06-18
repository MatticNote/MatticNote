package main

import "embed"

//go:embed template/*
//go:embed template/**/*
var template embed.FS

//go:embed client/dist/client/*
var webCli embed.FS

//go:embed client/dist/ui/*
var webUi embed.FS
