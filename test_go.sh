obj='{
  "a": "1",
  "c": "3",
  "nest": {
    "f": "1",
    "g": "3",
    "e": {
      "b": 2,
      "c": null,
      "a": "1"
    },
    "h": "4"
  },
  "d": "4",
  "b": {},
  "zzz": "zzzzzzzz",
  "name": "sortjson.nvim"
}
'
go run . $obj
go run . --reverse $obj
go run . --length $obj
go run . --length --reverse $obj

arr1='[
  {
    "a": "1",
    "c": "3",
    "nest": {
      "f": "1",
      "g": "3",
      "e": {
        "b": 2,
        "c": null,
        "a": "1"
      },
      "h": "4"
    },
    "d": "4",
    "b": {},
    "zzz": "zzzzzzzz",
    "name": "sortjson.nvim"
  }
]
'
go run . $arr1
go run . --reverse $arr1
go run . --length $arr1
go run . --length --reverse $arr1
