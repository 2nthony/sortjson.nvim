local M = {}

function M.sort_json(sort_option, reverse)
  local current_file_path = vim.api.nvim_buf_get_name(0)
  local jq_command = string.format(
  "jq 'def sort_recursive: if type == \"object\" then to_entries | sort_by(%s) | %s | from_entries | map_values(sort_recursive) else . end; sort_recursive'",
    sort_option, reverse and "reverse" or ".")

  local command = string.format("cat '%s' | %s", current_file_path, jq_command)
  local output = vim.fn.system(command)

  if vim.v.shell_error ~= 0 then
    local error_msg = string.format("Error: Failed to sort JSON. Command: %s\nOutput: %s", command, output)
    vim.api.nvim_err_writeln(error_msg)
    return
  end

  vim.api.nvim_buf_set_lines(0, 0, -1, false, vim.split(output, "\n"))
  vim.cmd("write")
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
