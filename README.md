# sortjson.nvim

NeoVIM plugin port of [vscode-sort-json](https://github.com/richie5um/vscode-sort-json).

https://github.com/2nthony/sortjson.nvim/assets/19513289/5d425e1b-28c5-4c3b-8d42-3e4b4d1dd266

## Requirements

- [jq](https://jqlang.github.io/jq/) use `jq` as the processor by default.

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
  -- options with default values
  opts = {
    jq = "jq", -- jq command, you can try `jaq` `gojq` etc.
  },
}
```

```sh
# Supported commands see `cmd = {}` above
:SortJSONByAlphaNum
```

## License

MIT &copy; [2nthony](https://github.com/sponsors/2nthony)
