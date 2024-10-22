local M = {}

function M.sort_json(sort_option, reverse)
  local bufnr = vim.api.nvim_get_current_buf()
  local lines = vim.api.nvim_buf_get_lines(bufnr, 0, -1, false)
  local json_text = table.concat(lines, "\n")

  local jq_command = string.format(
    "jq 'def sort_recursive: if type == \"object\" then to_entries | sort_by(%s) | %s | from_entries | map_values(sort_recursive) elif type == \"array\" then map(sort_recursive) else . end; sort_recursive'",
    sort_option, reverse and "reverse" or "."
  )

  local output = vim.fn.system(jq_command, json_text)


  if vim.v.shell_error ~= 0 then
    local error_msg = string.format("[sortjson.nvim] Failed to sort JSON. Error: %s", vim.trim(output))
    vim.notify(error_msg, vim.log.levels.WARN)
    return
  end

  -- Check if the output is valid JSON
  local success, _ = pcall(vim.fn.json_decode, output)
  if not success then
    vim.notify("[sortjson.nvim] Failed to sort JSON: Invalid JSON output", vim.log.levels.WARN)
    return
  end

  vim.api.nvim_buf_set_lines(bufnr, 0, -1, false, vim.split(output, "\n"))
end

function M.setup()
  vim.api.nvim_create_user_command("SortJSONByAlphaNum", function()
    M.sort_json(".key", false)
  end, {})

  vim.api.nvim_create_user_command("SortJSONByAlphaNumReverse", function()
    M.sort_json(".key", true)
  end, {})

  vim.api.nvim_create_user_command("SortJSONByKeyLength", function()
    M.sort_json(".key | length", false)
  end, {})

  vim.api.nvim_create_user_command("SortJSONByKeyLengthReverse", function()
    M.sort_json(".key | length", true)
  end, {})
end

return M
