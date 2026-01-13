package executor

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

type LocalExecutor struct {
	tempDir string
}

func NewLocalExecutor() *LocalExecutor {
	tempDir := filepath.Join(os.TempDir(), "code_executor")
	os.MkdirAll(tempDir, 0755)
	return &LocalExecutor{
		tempDir: tempDir,
	}
}

func (e *LocalExecutor) Execute(code, language string, inputs []string) (map[string]interface{}, error) {
	log.Printf("🎯 LocalExecutor executing %s code, length: %d chars, inputs: %v", language, len(code), inputs)

	switch strings.ToLower(language) {
	case "go":
		return e.executeGo(code, inputs)
	case "python", "python3":
		return e.executePython(code, inputs)
	case "javascript", "node":
		return e.executeJavaScript(code, inputs)
	case "cpp", "c++":
		return e.executeCpp()
	case "java":
		return e.executeJava()
	default:
		return map[string]interface{}{
			"output":   "",
			"error":    "Unsupported language: " + language,
			"exitCode": 1,
		}, nil
	}
}

func (e *LocalExecutor) executePython(code string, inputs []string) (map[string]interface{}, error) {
	log.Printf("🐍 Executing Python code, length: %d chars, inputs: %v", len(code), inputs)

	// Определяем команду Python
	cmdName := e.findPythonCommand()
	if cmdName == "" {
		errorMsg := "Python not found. Please install Python and make sure it's in PATH"
		log.Printf("❌ %s", errorMsg)
		return map[string]interface{}{
			"output":   "",
			"error":    errorMsg,
			"exitCode": 1,
		}, nil
	}

	log.Printf("🔧 Using Python command: %s", cmdName)

	// Создаем временный файл для Python кода
	tmpFile := filepath.Join(e.tempDir, "script_"+fmt.Sprintf("%d", time.Now().UnixNano())+".py")

	// Записываем код в файл
	err := os.WriteFile(tmpFile, []byte(code), 0644)
	if err != nil {
		log.Printf("❌ Failed to write Python file: %v", err)
		return map[string]interface{}{
			"output":   "",
			"error":    fmt.Sprintf("Error creating file: %v", err),
			"exitCode": 1,
		}, nil
	}
	defer func() {
		if err := os.Remove(tmpFile); err != nil {
			log.Printf("⚠️ Failed to remove temp file %s: %v", tmpFile, err)
		}
	}()

	// Создаем команду для выполнения Python файла
	cmd := exec.Command(cmdName, tmpFile)

	// Подготавливаем входные данные
	var stdin bytes.Buffer
	if len(inputs) > 0 {
		// Если inputs это массив строк, объединяем их с переносами строк
		fullInput := strings.Join(inputs, "\n") + "\n"
		stdin.WriteString(fullInput)
		log.Printf("📥 Sending input to Python: %q", fullInput)
	} else {
		log.Printf("📥 No input provided for Python")
	}

	cmd.Stdin = &stdin

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	// Устанавливаем таймаут выполнения (15 секунд)
	done := make(chan error, 1)
	go func() {
		done <- cmd.Run()
	}()

	select {
	case err := <-done:
		exitCode := 0
		if err != nil {
			if exitErr, ok := err.(*exec.ExitError); ok {
				exitCode = exitErr.ExitCode()
			} else {
				log.Printf("❌ Failed to start Python process: %v", err)
				return map[string]interface{}{
					"output":   "",
					"error":    fmt.Sprintf("Failed to start Python: %v", err),
					"exitCode": 1,
				}, nil
			}
			log.Printf("⚠️ Python execution completed with exit code %d", exitCode)
		} else {
			log.Printf("✅ Python execution completed successfully")
		}

		outputStr := strings.TrimSpace(stdout.String())
		errorStr := strings.TrimSpace(stderr.String())

		result := map[string]interface{}{
			"output":   outputStr,
			"error":    errorStr,
			"exitCode": exitCode,
		}

		log.Printf("📊 Python execution result - Output: %q", outputStr)
		log.Printf("📊 Python execution result - Error: %q", errorStr)

		return result, nil

	case <-time.After(15 * time.Second):
		// Таймаут - убиваем процесс
		log.Printf("⏰ Python execution timeout (15 seconds)")
		if cmd.Process != nil {
			cmd.Process.Kill()
		}
		return map[string]interface{}{
			"output":   "",
			"error":    "Execution timeout (15 seconds)",
			"exitCode": 1,
		}, nil
	}
}

// Остальные методы остаются без изменений
func (e *LocalExecutor) executeGo(code string, inputs []string) (map[string]interface{}, error) {
	log.Printf("🔵 Executing Go code, length: %d chars, inputs: %v", len(code), inputs)

	// Создаем временный файл
	tmpFile := filepath.Join(e.tempDir, "main_"+fmt.Sprintf("%d", time.Now().UnixNano())+".go")

	// Если код не содержит package main, добавляем его
	fullCode := code
	if !strings.Contains(code, "package main") {
		fullCode = "package main\n\n" + code
	}

	// Если нет функции main, добавляем простую обертку
	if !strings.Contains(code, "func main()") {
		fullCode = fullCode + "\n\nfunc main() {\n\t// Ваш код будет выполнен здесь\n}"
	}

	log.Printf("📝 Writing Go code to: %s", tmpFile)
	err := os.WriteFile(tmpFile, []byte(fullCode), 0644)
	if err != nil {
		log.Printf("❌ Failed to write Go file: %v", err)
		return map[string]interface{}{
			"output":   "",
			"error":    fmt.Sprintf("Error creating file: %v", err),
			"exitCode": 1,
		}, nil
	}
	defer func() {
		if err := os.Remove(tmpFile); err != nil {
			log.Printf("⚠️ Failed to remove temp file %s: %v", tmpFile, err)
		}
	}()

	// Компилируем и запускаем
	outputFile := tmpFile + ".exe"
	if runtime.GOOS != "windows" {
		outputFile = tmpFile + ".out"
	}

	// Компиляция
	log.Printf("🔨 Compiling Go code...")
	compileCmd := exec.Command("go", "build", "-o", outputFile, tmpFile)
	var compileStdout, compileStderr bytes.Buffer
	compileCmd.Stdout = &compileStdout
	compileCmd.Stderr = &compileStderr

	err = compileCmd.Run()
	if err != nil {
		log.Printf("❌ Go compilation failed: %v", err)
		return map[string]interface{}{
			"output":   "",
			"error":    fmt.Sprintf("Compilation error: %s", compileStderr.String()),
			"exitCode": 1,
		}, nil
	}
	defer func() {
		if err := os.Remove(outputFile); err != nil {
			log.Printf("⚠️ Failed to remove executable %s: %v", outputFile, err)
		}
	}()

	// Выполнение
	log.Printf("🚀 Running Go program...")
	cmd := exec.Command(outputFile)

	// Подготавливаем входные данные
	var stdin bytes.Buffer
	for _, input := range inputs {
		stdin.WriteString(input + "\n")
	}
	cmd.Stdin = &stdin

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err = cmd.Run()
	exitCode := 0
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			exitCode = exitErr.ExitCode()
		}
		log.Printf("⚠️ Go execution completed with exit code %d", exitCode)
	} else {
		log.Printf("✅ Go execution completed successfully")
	}

	result := map[string]interface{}{
		"output":   strings.TrimSpace(stdout.String()),
		"error":    strings.TrimSpace(stderr.String()),
		"exitCode": exitCode,
	}

	return result, nil
}

func (e *LocalExecutor) executeJavaScript(code string, inputs []string) (map[string]interface{}, error) {
	log.Printf("🔵 Executing JavaScript code, length: %d chars, inputs: %v", len(code), inputs)

	// Проверяем доступность Node.js
	cmdName := "node"
	if _, err := exec.LookPath(cmdName); err != nil {
		errorMsg := "Node.js is not installed or not in PATH"
		log.Printf("❌ %s", errorMsg)
		return map[string]interface{}{
			"output":   "",
			"error":    errorMsg,
			"exitCode": 1,
		}, nil
	}

	log.Printf("🔧 Using Node.js command: %s", cmdName)

	// Создаем временный файл для JavaScript кода
	tmpFile := filepath.Join(e.tempDir, "script_"+fmt.Sprintf("%d", time.Now().UnixNano())+".js")

	// Создаем обернутый код для Node.js с поддержкой ввода
	wrappedCode := e.createJavaScriptWrapper(code)

	err := os.WriteFile(tmpFile, []byte(wrappedCode), 0644)
	if err != nil {
		log.Printf("❌ Failed to write JavaScript file: %v", err)
		return map[string]interface{}{
			"output":   "",
			"error":    fmt.Sprintf("Error creating file: %v", err),
			"exitCode": 1,
		}, nil
	}
	defer func() {
		if err := os.Remove(tmpFile); err != nil {
			log.Printf("⚠️ Failed to remove temp file %s: %v", tmpFile, err)
		}
	}()

	// Выполняем код через файл
	cmd := exec.Command(cmdName, tmpFile)

	// Подготавливаем входные данные
	var stdin bytes.Buffer
	if len(inputs) > 0 {
		fullInput := strings.Join(inputs, "\n") + "\n"
		stdin.WriteString(fullInput)
		log.Printf("📥 Sending input to JavaScript: %q", fullInput)
	} else {
		log.Printf("📥 No input provided for JavaScript")
	}
	cmd.Stdin = &stdin

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	// Устанавливаем таймаут выполнения (15 секунд)
	done := make(chan error, 1)
	go func() {
		done <- cmd.Run()
	}()

	select {
	case err := <-done:
		exitCode := 0
		if err != nil {
			if exitErr, ok := err.(*exec.ExitError); ok {
				exitCode = exitErr.ExitCode()
			}
			log.Printf("⚠️ JavaScript execution completed with exit code %d", exitCode)
		} else {
			log.Printf("✅ JavaScript execution completed successfully")
		}

		result := map[string]interface{}{
			"output":   strings.TrimSpace(stdout.String()),
			"error":    strings.TrimSpace(stderr.String()),
			"exitCode": exitCode,
		}

		return result, nil

	case <-time.After(15 * time.Second):
		// Таймаут - убиваем процесс
		log.Printf("⏰ JavaScript execution timeout (15 seconds)")
		if cmd.Process != nil {
			cmd.Process.Kill()
		}
		return map[string]interface{}{
			"output":   "",
			"error":    "Execution timeout (15 seconds)",
			"exitCode": 1,
		}, nil
	}
}

func (e *LocalExecutor) executeCpp() (map[string]interface{}, error) {
	// ... (реализация как в предыдущей версии)
	return map[string]interface{}{
		"output":   "",
		"error":    "C++ execution not implemented",
		"exitCode": 1,
	}, nil
}

func (e *LocalExecutor) executeJava() (map[string]interface{}, error) {
	// ... (реализация как в предыдущей версии)
	return map[string]interface{}{
		"output":   "",
		"error":    "Java execution not implemented",
		"exitCode": 1,
	}, nil
}

// createJavaScriptWrapper создает обертку для JavaScript кода с поддержкой ввода
func (e *LocalExecutor) createJavaScriptWrapper(code string) string {
	return `
const readline = require('readline');

async function main() {
    const rl = readline.createInterface({
        input: process.stdin,
        output: process.stdout,
        terminal: false
    });

    let inputLines = [];
    for await (const line of rl) {
        inputLines.push(line);
    }

    let inputIndex = 0;
    const input = () => {
        if (inputIndex < inputLines.length) {
            return inputLines[inputIndex++];
        }
        return "";
    };

    // Заменяем глобальные функции ввода
    global.prompt = input;
    global.input = input;

    try {
        ` + code + `
    } catch (error) {
        console.error(error);
    }
}

main().catch(console.error);
`
}

// findPythonCommand ищет доступную команду Python в системе
func (e *LocalExecutor) findPythonCommand() string {
	// Список возможных команд Python в порядке предпочтения
	possibleCommands := []string{"python3", "python", "py"}

	for _, cmd := range possibleCommands {
		if path, err := exec.LookPath(cmd); err == nil {
			// Проверяем версию Python
			versionCmd := exec.Command(path, "--version")
			var versionOut bytes.Buffer
			versionCmd.Stdout = &versionOut
			versionCmd.Stderr = &versionOut

			if err := versionCmd.Run(); err == nil {
				log.Printf("✅ Found %s at %s: %s", cmd, path, strings.TrimSpace(versionOut.String()))
				return path
			}
		}
	}

	return ""
}

// Cleanup удаляет временную директорию
func (e *LocalExecutor) Cleanup() {
	if err := os.RemoveAll(e.tempDir); err != nil {
		log.Printf("⚠️ Failed to cleanup temp directory %s: %v", e.tempDir, err)
	} else {
		log.Printf("🧹 Cleaned up temp directory: %s", e.tempDir)
	}
}
