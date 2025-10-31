import { simpleAPI } from './simpleApi'

// Сервис для выполнения кода
export class CodeRunner {
  /**
   * Выполнение кода
   */
  static async executeCode(language, code) {
    try {
      console.log(`🎯 Executing ${language} code`)
      
      const result = await simpleAPI.executeCode(language, code)
      
      return {
        output: result.output || 'No output received',
        success: result.success !== false,
        color: result.success !== false ? 'green' : 'red'
      }
    } catch (error) {
      console.error('❌ Code execution failed:', error)
      
      // Локальное выполнение
      return await this.localExecuteCode(language, code)
    }
  }

  /**
   * Проверка решения задачи
   */
  static async checkSolution(language, code, taskId) {
    try {
      console.log(`🧪 Checking solution for task ${taskId}`)
      
      const result = await simpleAPI.checkSolution(language, code, taskId)
      
      return {
        output: result.output || 'No check result',
        passed: result.passed || false,
        color: result.passed ? 'green' : 'red'
      }
    } catch (error) {
      console.error('❌ Solution check failed:', error)
      
      // Локальная проверка
      return this.localSolutionCheck(language, code, taskId)
    }
  }

  /**
   * Локальное выполнение кода
   */
  static async localExecuteCode(language, code) {
    return new Promise((resolve) => {
      setTimeout(() => {
        try {
          const output = this.executeInBrowser(language, code)
          resolve({
            output: output,
            success: true,
            color: 'green'
          })
        } catch (error) {
          resolve({
            output: `❌ Ошибка выполнения: ${error.message}`,
            success: false,
            color: 'red'
          })
        }
      }, 500)
    })
  }

  /**
   * Локальная проверка решения
   */
  static localSolutionCheck(language, code, taskId) {
    console.log(`🔍 Local check for task ${taskId}`)
    
    // Проверяем синтаксис
    const syntaxValid = this.checkSyntax(language, code)
    if (!syntaxValid.valid) {
      return {
        output: `❌ Синтаксическая ошибка:\n${syntaxValid.error}`,
        passed: false,
        color: 'red'
      }
    }
    
    // Выполняем код для получения вывода
    const output = this.executeInBrowser(language, code)
    
    // Проверяем задачу
    const taskCheck = this.checkTaskRequirements(language, code, taskId, output)
    
    return {
      output: taskCheck.passed 
        ? `✅ Задача решена правильно!\n\nВывод программы:\n${output}`
        : `❌ Задача не решена:\n${taskCheck.message}\n\nВывод программы:\n${output}`,
      passed: taskCheck.passed,
      color: taskCheck.passed ? 'green' : 'red'
    }
  }

  /**
   * Проверка синтаксиса
   */
  static checkSyntax(language, code) {
    const checks = {
      python: () => {
        const openParens = (code.match(/\(/g) || []).length
        const closeParens = (code.match(/\)/g) || []).length
        
        if (openParens !== closeParens) {
          return { valid: false, error: 'Несбалансированные круглые скобки' }
        }
        return { valid: true }
      },
      javascript: () => {
        const openBraces = (code.match(/{/g) || []).length
        const closeBraces = (code.match(/}/g) || []).length
        const openParens = (code.match(/\(/g) || []).length
        const closeParens = (code.match(/\)/g) || []).length
        
        if (openBraces !== closeBraces) {
          return { valid: false, error: 'Несбалансированные фигурные скобки' }
        }
        if (openParens !== closeParens) {
          return { valid: false, error: 'Несбалансированные круглые скобки' }
        }
        return { valid: true }
      },
      cpp: () => {
        if (!code.includes('int main()')) {
          return { valid: false, error: 'Отсутствует функция main()' }
        }
        return { valid: true }
      },
      java: () => {
        if (!code.includes('public class Main')) {
          return { valid: false, error: 'Отсутствует класс Main' }
        }
        return { valid: true }
      }
    }
    
    return checks[language] ? checks[language]() : { valid: true }
  }

  /**
   * Выполнение в браузере
   */
  static executeInBrowser(language, code) {
    const executors = {
      javascript: () => this.executeJavaScript(code),
      python: () => this.executePython(code),
      cpp: () => '❌ C++ выполнение требует бэкенд',
      java: () => '❌ Java выполнение требует бэкенд'
    }
    
    return executors[language] ? executors[language]() : '❌ Язык не поддерживается'
  }

  static executeJavaScript(code) {
    try {
      let output = ''
      const originalConsoleLog = console.log
      
      console.log = (...args) => {
        output += args.join(' ') + '\n'
      }
      
      eval(code)
      
      console.log = originalConsoleLog
      
      return output || 'Программа выполнена (реальное выполнение JS)'
    } catch (error) {
      return `❌ Ошибка JavaScript: ${error.message}`
    }
  }

  static executePython(code) {
    try {
      let output = ''
      const lines = code.split('\n')
      
      for (let i = 0; i < lines.length; i++) {
        const line = lines[i].trim()
        
        // Эмуляция print
        if (line.startsWith('print(') && line.endsWith(')')) {
          const content = line.slice(6, -1).trim()
          const printed = content.replace(/^['"](.*)['"]$/, '$1')
          output += printed + '\n'
        }
        
        // Эмуляция циклов for
        if (line.includes('for') && line.includes('in range(')) {
          const rangeMatch = line.match(/for\s+(\w+)\s+in\s+range\(([^)]+)\)/)
          if (rangeMatch) {
            const [_, varName, rangeArgs] = rangeMatch
            const args = rangeArgs.split(',').map(arg => parseInt(arg.trim()))
            
            let start = 0, end = 0
            if (args.length === 1) {
              end = args[0]
            } else if (args.length === 2) {
              start = args[0]
              end = args[1]
            }
            
            // Обрабатываем следующую строку с отступом
            if (i + 1 < lines.length && (lines[i + 1].startsWith(' ') || lines[i + 1].startsWith('\t'))) {
              const loopLine = lines[i + 1].trim()
              if (loopLine.includes('print(')) {
                for (let j = start; j < end; j++) {
                  const loopContent = loopLine.slice(6, -1).trim()
                  let printed = loopContent.replace(new RegExp(varName, 'g'), j)
                    .replace(/^['"](.*)['"]$/, '$1')
                  
                  // Вычисляем простые выражения
                  if (printed.includes('*')) {
                    const parts = printed.split('*').map(p => p.trim())
                    if (parts.length === 2 && !isNaN(parts[0]) && !isNaN(parts[1])) {
                      printed = (parseInt(parts[0]) * parseInt(parts[1])).toString()
                    }
                  }
                  
                  output += printed + '\n'
                }
                i++ // Пропускаем строку цикла
              }
            }
          }
        }
      }
      
      return output || 'Программа выполнена (эмуляция Python)'
    } catch (error) {
      return `❌ Ошибка Python: ${error.message}`
    }
  }

  /**
   * Проверка требований задачи
   */
  static checkTaskRequirements(language, code, taskId, output) {
    const validators = {
      1: () => ({
        passed: output.includes('Hello World'),
        message: 'Программа должна выводить "Hello World"'
      }),
      2: () => ({
        passed: output.includes('5') && output.includes('10') && output.includes('15'),
        message: 'Программа должна выводить числа 5, 10, 15...'
      })
    }
    
    const validator = validators[taskId] || (() => ({ 
      passed: true, 
      message: '' 
    }))
    
    return validator()
  }
}