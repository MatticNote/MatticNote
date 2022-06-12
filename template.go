package main

import "embed"

//go:embed template/*
//go:embed template/**/*
var template embed.FS
