local defaults = {
  log_level = "WARN",
}

local function plugin_dir()
  local info = debug.getinfo(1, "S")
  local source = info and info.source or ""

  if vim.startswith(source, "@") then
    source = source:sub(2)
  end

  return vim.fn.fnamemodify(source, ":h:h:h")
end

local M = {}

function M.sort_json(opts, length, reverse)
  local bufnr = vim.api.nvim_get_current_buf()
  local lines = vim.api.nvim_buf_get_lines(bufnr, 0, -1, false)
  local json_text = table.concat(lines, "\n")

  local command = plugin_dir() .. "/cli"
  if vim.fn.executable(command) ~= 1 then
    local error_msg = string.format("[sortjson.nvim] Failed to sort JSON. Error: %s is not executable", command)
    vim.notify(error_msg, vim.log.levels[opts.log_level])
    return
  end

  if length then
    command = command .. " --length"
  end
  if reverse then
    command = command .. " --reverse"
  end

  local wrapped_json_text = "'" .. json_text .. "'"
  local output = vim.fn.system(command .. " " .. wrapped_json_text)

  if vim.v.shell_error ~= 0 then
    local error_msg = string.format("[sortjson.nvim] Failed to sort JSON. Error: %s", vim.trim(output))
    vim.notify(error_msg, vim.log.levels[opts.log_level])
    return
  end

  -- Check if the output is valid JSON
  local success, _ = pcall(vim.json and vim.json.decode or vim.fn.json_decode, output)
  if not success then
    vim.notify("[sortjson.nvim] Failed to sort JSON: Invalid JSON output", vim.log.levels[opts.log_level])
    return
  end

  vim.api.nvim_buf_set_lines(bufnr, 0, -1, false, vim.split(output, "\n"))
end

function M.setup(user_opts)
  user_opts = user_opts or {}
  local opts = vim.tbl_extend("force", defaults, user_opts)

  vim.api.nvim_create_user_command("SortJSONByAlphaNum", function()
    M.sort_json(opts, false, false)
  end, {})

  vim.api.nvim_create_user_command("SortJSONByAlphaNumReverse", function()
    M.sort_json(opts, false, true)
  end, {})

  vim.api.nvim_create_user_command("SortJSONByKeyLength", function()
    M.sort_json(opts, true, false)
  end, {})

  vim.api.nvim_create_user_command("SortJSONByKeyLengthReverse", function()
    M.sort_json(opts, true, true)
  end, {})
end

return M
