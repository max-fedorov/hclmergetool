// Copyright (c) 2022, Max Fedorov <mail@skam.in>

package main

import (
	"fmt"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"os"
	"strings"
)

type Block struct {
	Self       *hclwrite.Block
	Body       *hclwrite.Body
	Type       string
	Name       []string
	Blocks     []*hclwrite.Block
	Attributes map[string]*hclwrite.Attribute
}

func ReadHclFile(path string) *hclwrite.File {
	data, err := os.ReadFile(path)
	if err != nil {
		fmt.Printf("ERROR: Failed to open: %v\n", err)
		os.Exit(1)
	}

	hclData, diag := hclwrite.ParseConfig(data, path, hcl.InitialPos)
	if diag.HasErrors() {
		fmt.Printf("ERROR: Failed to parse: %v\n", diag)
		os.Exit(1)
	}
	return hclData
}

func getChildBlocks(blocks []*hclwrite.Block) []Block {
	inner_blocks := []Block{}
	for _, b := range blocks {
		inner_blocks = append(inner_blocks, *parseBlock(b))
	}
	return inner_blocks
}

func parseBlock(block *hclwrite.Block) *Block {
	newBlock := new(Block)
	newBlock.Self = block
	newBlock.Type = block.Type()
	newBlock.Name = block.Labels()
	newBlock.Attributes = block.Body().Attributes()
	newBlock.Blocks = block.Body().Blocks()
	newBlock.Body = block.Body()
	return newBlock
}

func getBlocksByType(blocks []Block, type_name string) []Block {
	out_blocks := []Block{}
	for _, b := range blocks {
		if strings.Compare(b.Type, type_name) == 0 {
			out_blocks = append(out_blocks, b)
		}
	}
	return out_blocks
}

func getBlocksByLabels(blocks []Block, labels []string) []Block {
	out_blocks := []Block{}
	for _, b := range blocks {
		if len(b.Name) >= len(labels) {
			for k, v := range labels {
				//fmt.Printf("compare block labels: '%s' == '%s' ?", b.Name[k], v)
				if b.Name[k] != v {
					goto out
				}
			}
			out_blocks = append(out_blocks, b)
		out:
		} else {
			if Equal(b.Name, labels) {
				out_blocks = append(out_blocks, b)
			}
		}
	}
	return out_blocks
}

func appendBlock(orig *Block, template *Block, updateArgs *bool) *Block {
	for _, tv := range getChildBlocks(template.Body.Blocks()) {
		//fmt.Printf("Block_template: %d %s %v", tk, tv.Type, tv.Name)
		blocks := getBlocksByType(getChildBlocks(orig.Blocks), tv.Type)
		blocks = getBlocksByLabels(blocks, tv.Name)
		if len(blocks) == 0 {
			orig.Body.AppendNewline()
			orig.Body.AppendBlock((*hclwrite.Block)(tv.Self))
		} else {
			for _, cv := range blocks {
				//fmt.Printf("Block_config: %d %s %v", ck, cv.Type, cv.Name)
				appendBlock(&cv, &tv, updateArgs)
			}
		}
	}

	for ck, cv := range template.Attributes {
		//fmt.Printf("Attr_update: %s %v", ck, cv.Expr())
		_, argExist := orig.Self.Body().Attributes()[ck]
		if *updateArgs || (!argExist && !*updateArgs) {
			orig.Attributes[ck] = cv
			orig.Body.SetAttributeRaw(ck, cv.Expr().BuildTokens(nil))
		}
	}
	return orig
}

func Process(config *hclwrite.File, template *hclwrite.File, updateArgs *bool) *hclwrite.File {
	rootBlocksConfig := getChildBlocks(config.Body().Blocks())
	rootBlocksTemplate := getChildBlocks(template.Body().Blocks())

	rootAttrsTemplate := template.Body().Attributes()

	for _, tv := range rootBlocksTemplate {
		//fmt.Printf("Block_template: %d %s %v", tk, tv.Type, tv.Name)
		blocks := getBlocksByType(rootBlocksConfig, tv.Type)
		blocks = getBlocksByLabels(blocks, tv.Name)
		for _, cv := range blocks {
			//fmt.Printf("Block_config: %d %s %v", ck, cv.Type, cv.Name)
			appendBlock(&cv, &tv, updateArgs)
		}
	}

	for ck, cv := range rootAttrsTemplate {
		_, argExist := config.Body().Attributes()[ck]
		if *updateArgs || (!argExist && !*updateArgs) {
			config.Body().AppendNewline()
			config.Body().AppendNewline()
			config.Body().SetAttributeRaw(ck, cv.Expr().BuildTokens(nil))
		}
	}

	return config
}
