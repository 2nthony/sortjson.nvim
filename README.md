# sortjson.nvim

NeoVIM plugin port of [vscode-sort-json](https://github.com/richie5um/vscode-sort-json).

## Requirements

- [Node.JS](https://nodejs.org/en/) LTS, native welcome.

## Usage

```lua
-- lazy.nvim
return {
  "2nthony/sortjson.nvim",
  cmd = {
    "SortJSONByAlphaNum",
    "SortJSONByAlphaNumReverse",
    "SortJSONByKeyLength",
    "SortJSONByKeyLengthReverse",
  },
  config = true,
}
```

```sh
# Supported commands see `cmd = {}` above
:SortJSONByAlphaNum
```

## License

MIT &copy; [2nthony](https://github.com/sponsors/2nthony)
