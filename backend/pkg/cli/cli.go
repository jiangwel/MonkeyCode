package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/chaitin/MonkeyCode/backend/domain"
)

// RunParseCLI
//  @Description          ctcode-cli 支持通过 CLI 的方式获取文件的 AST 信息
//	@param command        ["parse"] - 解析文件为 AST 树
//	@param flag           ["-s", "--successOnly"] - 是否只返回成功的结果
//	@param paths          要解析的文件路径
//	@return []ParseResult 解析结果数组
//	@return error         错误信息
func RunParseCLI(command string, flag string, paths ...string) ([]domain.ParseResult, error) {
	cmd := exec.Command("ctcode-cli", command, flag, strings.Join(paths, " "))
	cmd.Env = os.Environ() 
	output, err := cmd.CombinedOutput()
	fmt.Printf(`err: %s, output: %s\n`, fmt.Sprint(err), string(output))
	if err != nil {
		return []domain.ParseResult{}, err
	}
	var results []domain.ParseResult
	if err := json.Unmarshal(output, &results); err != nil {
		panic(err)
	}
	return results, nil
}
