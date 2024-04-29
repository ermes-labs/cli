<div align="center">
  
<br>

<img src="https://raw.githubusercontent.com/ermes-labs/docs/main/docs/public/icon.png" width="30%">

<h1>@ermes-labs/cli</h1>

CLI for the [`Ermes`](https://ermes-labs.github.io/docs) framework

[![language: Go](https://img.shields.io/badge/go-language-50b7e0?style=flat-square&logo=go)](https://go.dev/)
[![Github CI](https://img.shields.io/github/actions/workflow/status/ermes-labs/cli/ci.yml?style=flat-square&branch=main)](https://github.com/ermes-labs/cli/actions/workflows/ci.yml)
[![Codecov](https://img.shields.io/codecov/c/github/ermes-labs/cli?color=44cc11&logo=codecov&style=flat-square)](https://codecov.io/gh/ermes-labs/cli)
![GitHub Latest Release)](https://img.shields.io/github/v/release/ermes-labs/cli?logo=github)

</div>

# Introduction 📖

To uild the cli

##

```
CLI for Ermes

Usage:
  ermes-cli [command]

Available Commands:
  build       Builds OpenFaaS function containers
  check       Check the infrastructure
  completion  Generate the autocompletion script for the specified shell
  deploy      Deploy a function to specified infrastructure.
  help        Help about any command
  new         Create a new template in the current folder with the name given as name
  print       Print the infrastructure
  pull        Downloads templates from the specified git repo
  push        Push OpenFaaS functions to remote registry (Docker Hub)

Flags:
  -h, --help      help for ermes-cli
  -v, --version   version for ermes-cli

Use "ermes-cli [command] --help" for more information about a command.
```

# Query Language

Queries can target specific areas using their identifiers and filter them:

- `#AreaId`: Targets a specific area by its unique identifier.
- `AreaId(filters...)`: Targets a specific area by its unique identifier.
- `*(filters...)`: A wildcard character representing all areas.

## Filters

Filters allow for refined queries based on specific criteria, multiple filters can be combine

- **Tag Filters**: Specify conditions based on tags assigned to nodes. The syntax for a tag filter is `tag:<key>=<value>`, where Key is the tag name, and Value is the tag value. <br> e.g. `tag:key=value`
- **Level Filters**: Define conditions based on the hierarchical level of nodes. The syntax for a level filter is `level<op><number>` is a comparison operator `(<, <=, >, >=, =, !=)`, and Numb is an integer representing the level.
  Filters can be combined within a single query to apply multiple criteria:

mathematica
Copy code
AreaId(tag:Key=Value, level<Comp>Numb)
Sets and Operations
Queries can be grouped into sets, allowing for union and difference operations to combine or exclude specific areas:

less
Copy code
{ #AreaId1 + #AreaId2 - #AreaId3 }
{}: Encloses a set of area queries.
+: Adds areas to the set (union operation).
-: Removes areas from the set (difference operation).
Formal Language Definition
The formal grammar of the query language is defined as follows:

```ts
Query       ::= Set
Set         ::= "#" Ident | "*" | NodesInArea | "{" Sets "}"
NodesInArea ::= Ident { "(" Filters ")" }*
Filters     ::= Filter { "," Filter }*
Filter      ::= TagFilter | LevelFilter
TagFilter   ::= "tag:" Ident "=" String
LevelFilter ::= "level" Comp Number
Sets        ::= Set { Op Set }*
Op          ::= "+" | "-"
Ident       ::= [a-zA-Z_][a-zA-Z0-9_]*
String      ::= [^,)]+
Number      ::= [-+]?[0-9]+
Comp        ::= "<" | "<=" | ">" | ">=" | "=" | "!="

Ident: An identifier starting with a letter or underscore, followed by any combination of letters, digits, and underscores.
String: A sequence of characters representing tag values, excluding commas and closing parentheses.
Number: An integer, which may include an optional sign.
Comp: Comparison operators used in level filters.
Examples
Query a single area by its identifier:
```

bash
Copy code
#Downtown
Query all areas:

markdown
Copy code

- Apply a tag filter to an area:

scss
Copy code
ParkArea(tag:type=park)
Apply multiple filters to an area:

scss
Copy code
ResidentialArea(tag:status=active, level>=2)
Combine multiple area queries into a set with operations:

less
Copy code
{ #Commercial + #Industrial - #Restricted }
