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
		return e.executeCpp(code, inputs)
	case "java":
		return e.executeJava(code, inputs)
	default:
		return map[string]interface{}{
			"output":   "",
			"error":    "Unsupported language: " + language,
			"exitCode": 1,
		}, nil
	}
}

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
		"output":   stdout.String(),
		"error":    stderr.String(),
		"exitCode": exitCode,
	}

	log.Printf("📊 Go execution result - Output: %d chars, Error: %d chars",
		len(stdout.String()), len(stderr.String()))

	return result, nil
}

func (e *LocalExecutor) executePython(code string, inputs []string) (map[string]interface{}, error) {
	log.Printf("🐍 Executing Python code, length: %d chars, inputs: %v", len(code), inputs)

	// Создаем временный файл
	tmpFile := filepath.Join(e.tempDir, "script_"+fmt.Sprintf("%d", time.Now().UnixNano())+".py")

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

	// Выполняем код с таймаутом
	cmd := exec.Command(cmdName, tmpFile)

	// Подготавливаем входные данные
	var stdin bytes.Buffer
	for _, input := range inputs {
		stdin.WriteString(input + "\n")
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

		result := map[string]interface{}{
			"output":   stdout.String(),
			"error":    stderr.String(),
			"exitCode": exitCode,
		}

		log.Printf("📊 Python execution result - Output: %d chars, Error: %d chars",
			len(stdout.String()), len(stderr.String()))

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

func (e *LocalExecutor) executeJavaScript(code string, inputs []string) (map[string]interface{}, error) {
	log.Printf("🔵 Executing JavaScript code, length: %d chars, inputs: %v", len(code), inputs)

	// Для Node.js нужно создать обертку для эмуляции prompt()
	wrappedCode := e.wrapJavaScriptCode(code)

	tmpFile := filepath.Join(e.tempDir, "script_"+fmt.Sprintf("%d", time.Now().UnixNano())+".js")

	log.Printf("📝 Writing JavaScript code to: %s", tmpFile)
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

	// Выполняем код с таймаутом
	cmd := exec.Command(cmdName, tmpFile)

	// Подготавливаем входные данные
	var stdin bytes.Buffer
	for _, input := range inputs {
		stdin.WriteString(input + "\n")
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
			"output":   stdout.String(),
			"error":    stderr.String(),
			"exitCode": exitCode,
		}

		log.Printf("📊 JavaScript execution result - Output: %d chars, Error: %d chars",
			len(stdout.String()), len(stderr.String()))

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

func (e *LocalExecutor) executeCpp(code string, inputs []string) (map[string]interface{}, error) {
	log.Printf("🔷 Executing C++ code, length: %d chars, inputs: %v", len(code), inputs)

	// Создаем временный файл
	tmpFile := filepath.Join(e.tempDir, "program_"+fmt.Sprintf("%d", time.Now().UnixNano())+".cpp")

	// Если нет #include <iostream>, добавляем его
	fullCode := code
	if !strings.Contains(code, "#include") {
		fullCode = "#include <iostream>\nusing namespace std;\n\n" + code
	}

	// Если нет функции main, добавляем
	if !strings.Contains(code, "int main()") && !strings.Contains(code, "void main()") {
		fullCode = fullCode + "\n\nint main() {\n\t// Ваш код будет выполнен здесь\n\treturn 0;\n}"
	}

	log.Printf("📝 Writing C++ code to: %s", tmpFile)
	err := os.WriteFile(tmpFile, []byte(fullCode), 0644)
	if err != nil {
		log.Printf("❌ Failed to write C++ file: %v", err)
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

	// Компилируем
	outputFile := tmpFile + ".exe"
	if runtime.GOOS != "windows" {
		outputFile = tmpFile + ".out"
	}

	// Пробуем разные компиляторы
	compilers := []string{"g++", "clang++", "c++"}
	var compileErr error
	var compileStderr bytes.Buffer

	for _, compiler := range compilers {
		if _, err := exec.LookPath(compiler); err != nil {
			continue
		}

		log.Printf("🔨 Compiling C++ code with %s...", compiler)
		compileCmd := exec.Command(compiler, "-o", outputFile, tmpFile)
		compileStderr.Reset()
		compileCmd.Stderr = &compileStderr

		compileErr = compileCmd.Run()
		if compileErr == nil {
			log.Printf("✅ C++ compilation successful with %s", compiler)
			break
		}
		log.Printf("⚠️ C++ compilation failed with %s: %v", compiler, compileErr)
	}

	if compileErr != nil {
		log.Printf("❌ All C++ compilers failed")
		return map[string]interface{}{
			"output":   "",
			"error":    fmt.Sprintf("Compilation error: %s\nNo C++ compiler found (tried: %v)", compileStderr.String(), compilers),
			"exitCode": 1,
		}, nil
	}
	defer func() {
		if err := os.Remove(outputFile); err != nil {
			log.Printf("⚠️ Failed to remove executable %s: %v", outputFile, err)
		}
	}()

	// Выполняем
	log.Printf("🚀 Running C++ program...")
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
		log.Printf("⚠️ C++ execution completed with exit code %d", exitCode)
	} else {
		log.Printf("✅ C++ execution completed successfully")
	}

	result := map[string]interface{}{
		"output":   stdout.String(),
		"error":    stderr.String(),
		"exitCode": exitCode,
	}

	log.Printf("📊 C++ execution result - Output: %d chars, Error: %d chars",
		len(stdout.String()), len(stderr.String()))

	return result, nil
}

func (e *LocalExecutor) executeJava(code string, inputs []string) (map[string]interface{}, error) {
	log.Printf("☕ Executing Java code, length: %d chars, inputs: %v", len(code), inputs)

	// Создаем временный файл с уникальным именем класса
	className := "Main_" + fmt.Sprintf("%d", time.Now().UnixNano())
	tmpFile := filepath.Join(e.tempDir, className+".java")

	// Если код не содержит public class, добавляем обертку
	fullCode := code
	if !strings.Contains(code, "public class") {
		fullCode = "public class " + className + " {\n    public static void main(String[] args) {\n        " +
			strings.ReplaceAll(code, "\n", "\n        ") +
			"\n    }\n}"
	} else {
		// Если уже есть public class, обновляем имя класса
		fullCode = strings.Replace(fullCode, "public class Main", "public class "+className, 1)
		fullCode = strings.Replace(fullCode, "public class main", "public class "+className, 1)
	}

	log.Printf("📝 Writing Java code to: %s", tmpFile)
	err := os.WriteFile(tmpFile, []byte(fullCode), 0644)
	if err != nil {
		log.Printf("❌ Failed to write Java file: %v", err)
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

	// Компилируем
	log.Printf("🔨 Compiling Java code...")
	compileCmd := exec.Command("javac", tmpFile)
	var compileStderr bytes.Buffer
	compileCmd.Stderr = &compileStderr

	err = compileCmd.Run()
	if err != nil {
		log.Printf("❌ Java compilation failed: %v", err)
		return map[string]interface{}{
			"output":   "",
			"error":    fmt.Sprintf("Compilation error: %s", compileStderr.String()),
			"exitCode": 1,
		}, nil
	}
	defer func() {
		classFile := filepath.Join(e.tempDir, className+".class")
		if err := os.Remove(classFile); err != nil {
			log.Printf("⚠️ Failed to remove class file %s: %v", classFile, err)
		}
	}()

	// Выполняем
	log.Printf("🚀 Running Java program...")
	cmd := exec.Command("java", "-cp", e.tempDir, className)

	// Подготавливаем входные данные
	var stdin bytes.Buffer
	for _, input := range inputs {
		stdin.WriteString(input + "\n")
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
			log.Printf("⚠️ Java execution completed with exit code %d", exitCode)
		} else {
			log.Printf("✅ Java execution completed successfully")
		}

		result := map[string]interface{}{
			"output":   stdout.String(),
			"error":    stderr.String(),
			"exitCode": exitCode,
		}

		log.Printf("📊 Java execution result - Output: %d chars, Error: %d chars",
			len(stdout.String()), len(stderr.String()))

		return result, nil

	case <-time.After(15 * time.Second):
		// Таймаут - убиваем процесс
		log.Printf("⏰ Java execution timeout (15 seconds)")
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

// wrapJavaScriptCode создает обертку для эмуляции prompt() в Node.js
func (e *LocalExecutor) wrapJavaScriptCode(code string) string {
	return `
const readline = require('readline');

// Эмуляция prompt() для Node.js
function prompt() {
    const rl = readline.createInterface({
        input: process.stdin,
        output: process.stdout
    });
    
    return new Promise((resolve) => {
        rl.question('', (answer) => {
            rl.close();
            resolve(answer);
        });
    });
}

// Заменяем глобальный prompt
global.prompt = prompt;

// Основной код
(async function() {
    ` + code + `
})().catch(console.error);
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
