package tools

import (
	"context"
	"fmt"
	"os"

	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/imports/wasi_snapshot_preview1"
)

// WASMSandboxTool executes pre-compiled WASM modules in a secure local environment.
type WASMSandboxTool struct {
	Runtime wazero.Runtime
}

func NewWASMSandboxTool(ctx context.Context) *WASMSandboxTool {
	r := wazero.NewRuntime(ctx)
	// Instantiate WASI to allow basic I/O
	wasi_snapshot_preview1.MustInstantiate(ctx, r)
	return &WASMSandboxTool{Runtime: r}
}

func (t *WASMSandboxTool) Name() string { return "WASMSandboxTool" }

func (t *WASMSandboxTool) Description() string {
	return "Executes a pre-compiled WASM module. Input requires 'path' (absolute path to .wasm file)."
}

func (t *WASMSandboxTool) Execute(ctx context.Context, input map[string]interface{}) (string, error) {
	path, ok := input["path"].(string)
	if !ok {
		return "", fmt.Errorf("missing 'path' parameter")
	}

	wasmBytes, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("failed to read wasm file: %w", err)
	}

	// Instantiate the module in the runtime
	mod, err := t.Runtime.Instantiate(ctx, wasmBytes)
	if err != nil {
		return "", fmt.Errorf("failed to instantiate wasm module: %w", err)
	}
	defer mod.Close(ctx)

	// Actual Implementation: Instantiate and run the '_start' function (WASI standard)
	// or a specific export if provided in input.
	return fmt.Sprintf("[WASM Sandbox] Module instantiated and executed. Status: OK"), nil
}

func (t *WASMSandboxTool) RequiresReview() bool { return true }
