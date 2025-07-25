package cli

import (
	"encoding/json"
	"os"
	"os/exec"

	"github.com/chaitin/MonkeyCode/backend/domain"
)

// RunCli
//  @Tags			                      WorkspaceFile
//  @Description:                         运行ctcode-cli命令
//	@param command                        命令
//	@param flag	                          [ -m | --maxlines <int> ] 解析时允许的最大 CodeSnippet 行数
//	@param codeFiles                      代码文件信息
//	@return string                        输出结果
//	@return error
func RunCli(command string, flag string, codeFiles domain.CodeFiles) ([]domain.IndexResult, error) {
	inputJson, err := json.Marshal(codeFiles) 
	if err != nil {
		return []domain.IndexResult{}, err 
	}
	cmd := exec.Command("ctcode-cli", command, flag, string(inputJson)) 
	cmd.Env = os.Environ() 
	output, err := cmd.CombinedOutput() 
	
	if err != nil {
		return []domain.IndexResult{}, err 
	}
	var res []domain.IndexResult 
	err = json.Unmarshal(output, &res) 
	if err!= nil {
		return []domain.IndexResult{}, err
	} 
	return res, nil 
}
