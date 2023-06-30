local M = {}

local function get_plugin_dir()
  local str = debug.getinfo(2, "S").source:sub(2)
  return vim.fn.fnamemodify(str, ":h")
end

-- nodejs
local __dirname = get_plugin_dir()

local build_ipc_command = function(nvim_command, ipc_option)
  vim.api.nvim_create_user_command(nvim_command, function()
    local current_file_path = vim.api.nvim_buf_get_name(0)
    local output = vim.fn.system({ "node", __dirname .. "/cli.js", current_file_path, ipc_option })
    if output then
      print(output)
    end

    -- reload current file
    vim.cmd([[:e]])
  end, {})
end

M.setup = function()
  build_ipc_command("SortJSONByAlphaNum", "--alpha-num")
  build_ipc_command("SortJSONByAlphaNumReverse", "--alpha-num-reverse")
  build_ipc_command("SortJSONByKeyLength", "--key-length")
  build_ipc_command("SortJSONByKeyLengthReverse", "--key-length-reverse")
end

return M
