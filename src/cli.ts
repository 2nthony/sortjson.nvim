import { CAC } from 'cac'
import { loadJsonFile } from 'load-json-file'
import type { Options } from 'write-json-file'
import { writeJsonFile } from 'write-json-file'

type Mutable<T> = { -readonly [P in keyof T]: T[P] }

main()

async function main() {
  const cli = new CAC()

  cli
    .option('--alpha-num', '')
    .option('--alpha-num-reverse', '')
    .option('--key-length', '')
    .option('--key-length-reverse', '')
    .help()

  const { args, options } = cli.parse()
  const file = args[0]

  // no file provided, return early
  if (!file)
    return

  // if filepath is a json file
  if (file.endsWith('.json')) {
    const sortOptions: Mutable<Options> = {
      detectIndent: true,
    }

    try {
      const json = await loadJsonFile(file)

      if (options.alphaNum)
        sortOptions.sortKeys = (a, b) => a.localeCompare(b)

      if (options.alphaNumReverse)
        sortOptions.sortKeys = (a, b) => b.localeCompare(a)

      if (options.keyLength)
        sortOptions.sortKeys = (a, b) => a.length - b.length

      if (options.keyLengthReverse)
        sortOptions.sortKeys = (a, b) => b.length - a.length

      await writeJsonFile(file, json, sortOptions)
    }
    catch {
      console.error('[sortjson.nvim]: Error read/write json file')
    }
  }
  else {
    console.error('[sortjson.nvim]: Please provide a json file')
  }
}
