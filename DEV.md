# golang-language-server

> nvim: lsp/golsp.lua
```lua
local lspconfig = require("lspconfig.configs")
lspconfig.golsp = {
  default_config = {
    cmd = {"golang-language-server"},
    filetypes = {"go"},
    single_file_support = true,
    root_markers = { '.git', 'build', 'cmake' },
  }
}

return {
  cmd = vim.lsp.rpc.connect('127.0.0.1',9999),
  filetypes = { "go" },
  root_markers = { '.git', 'build', 'go.mod', 'go.sum'},
}
```

### start
make run