import { WebContainer } from '@webcontainer/api'

class WebCompiler {
  static webcontainerInstance = null
  static isInitialized = false

  // Инициализация WebContainer
  static async initialize() {
    if (this.isInitialized) return
    
    try {
      console.log('🚀 Initializing WebContainer...')
      this.webcontainerInstance = await WebContainer.boot()
      this.isInitialized = true
      console.log('✅ WebContainer initialized')
    } catch (error) {
      console.error('❌ WebContainer initialization failed:', error)
      throw error
    }
  }

  // Выполнение кода
  static async executeCode(language, code) {
    if (!this.isInitialized) {
      await this.initialize()
    }

    try {
      const fileName = this.getFileName(language)
      const command = this.getRunCommand(language)
      
      // Записываем код в файл
      await this.webcontainerInstance.fs.writeFile(fileName, code)
      
      // Выполняем команду
      const process = await this.webcontainerInstance.spawn(command.cmd, command.args)
      
      let output = ''
      let errorOutput = ''
      
      process.output.pipeTo(new WritableStream({
        write(data) {
          output += data
        }
      }))
      
      process.stderr.pipeTo(new WritableStream({
        write(data) {
          errorOutput += data
        }
      }))
      
      const exitCode = await process.exit
      
      return {
        success: exitCode === 0,
        output: exitCode === 0 ? output : errorOutput,
        exitCode: exitCode
      }
      
    } catch (error) {
      return {
        success: false,
        output: `❌ Execution error: ${error.message}`,
        exitCode: -1
      }
    }
  }

  static getFileName(language) {
    const files = {
      python: 'main.py',
      javascript: 'main.js',
      cpp: 'main.cpp',
      java: 'Main.java'
    }
    return files[language] || 'main.txt'
  }

  static getRunCommand(language) {
    const commands = {
      python: { cmd: 'python', args: ['main.py'] },
      javascript: { cmd: 'node', args: ['main.js'] },
      cpp: { 
        cmd: 'sh', 
        args: ['-c', 'g++ main.cpp -o main && ./main']
      },
      java: { cmd: 'java', args: ['Main.java'] }
    }
    return commands[language] || { cmd: 'echo', args: ['Unsupported language'] }
  }
}

export { WebCompiler }
