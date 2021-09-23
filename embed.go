package scaffolding

import "embed"

// Embeds has the entire ./embeds directory embed in read-only mode
//go:embed embeds/* env/*
var Embeds embed.FS
