# About `hclmergetool`

Works with HashiCorp Configuration Language (HCL) files. Allows you to supplement the _input_ file with blocks and attributes from the _template_ file. You can use any file name (file extension doesn't matter) for the _input_ and _template_ files. The only thing required is valid HCL syntax.

# Installation

## Binary Release

Latest releases are available here on Github -> [Releases](https://github.com/max-fedorov/hclmergetool/releases) . `hclmergetool` is a single binary, installation is as simple as placing the binary in your PATH.


## Install with Go

    go install github.com/max-fedorov/hclmergetool@latest

## Update version with Go

    go get -u github.com/max-fedorov/hclmergetool@latest

# Usage

    hclmergetool -i [FILE] -t [FILE]

Arguments:

    -i  --input         path to HCL input file
    -t  --template      path to HCL template file
    -o  --output        path to HCL output file. If not specified, print to stdout

    -u  --update-args   if set, the values for existing arguments will be updated from
                        the template values, otherwise, the default value remains
                  
    -v  --version       show version
    -h  --help          print help information
  
# Example

    hclmergetool --update-args -i current.tf -t current_templ.tf -o current_override.tf

- current.tf

```hcl
resource "aws_subnet" "public-1" {
    var_a = asd
    var_b = old_value
    old_block {
        var_block_a = qwe
    }
}

resource "aws_subnet" "public-2" {
    var_a = asd
    var_b = old_value
    old_block {
        var_block_a = qwe
    }
    new_block {

    }
}

resource "aws_subnet" "public-3" {
    var_a = asd
    var_b = old_value
    var_c = qwe
}

resource "aws_subnet_local" "public-4" {
}
```

- current_templ.tf

```hcl
resource "aws_subnet" {
    var_a = asd
    var_b = new_value
    var_c = qwe_new

    old_block {
        var_block_a = qwe
        var_block_b = qwe_new
    }

    new_block {
        var_block_c = qwe
    }
}

```

- current_override.tf

```hcl
resource "aws_subnet" "public-1" {
    var_a = asd
    var_b = new_value
    old_block {
        var_block_a = qwe
        var_block_b = qwe_new
    }

    new_block {
        var_block_c = qwe
    }
    var_c = qwe_new
}

resource "aws_subnet" "public-2" {
    var_a = asd
    var_b = new_value
    old_block {
        var_block_a = qwe
        var_block_b = qwe_new
    }
    new_block {

        var_block_c = qwe
    }
    var_c = qwe_new
}

resource "aws_subnet" "public-3" {
    var_a = asd
    var_b = new_value
    var_c = qwe_new

    old_block {
        var_block_a = qwe
        var_block_b = qwe_new
    }

    new_block {
        var_block_c = qwe
    }
}

resource "aws_subnet_local" "public-4" {
}
```
