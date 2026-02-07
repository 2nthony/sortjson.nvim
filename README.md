# sortjson.nvim

A NeoVIM plugin that can sort current JSON file by key name.

https://github.com/2nthony/sortjson.nvim/assets/19513289/5d425e1b-28c5-4c3b-8d42-3e4b4d1dd266

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
    log_level = "WARN", -- log level, see `:h vim.log.levels`, print error info when parsing json failed
  },
}
```

```sh
# Supported commands see `cmd = {}` above
:SortJSONByAlphaNum
```

## License

MIT &copy; [2nthony](https://github.com/sponsors/2nthony)
