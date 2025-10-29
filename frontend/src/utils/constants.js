export const LANGUAGES = [
  { id: 'python', name: 'Python', icon: '🐍' },
  { id: 'javascript', name: 'JavaScript', icon: '📜' },
  { id: 'cpp', name: 'C++', icon: '⚡' },
  { id: 'java', name: 'Java', icon: '☕' }
]

export const DIFFICULTY_LEVELS = {
  beginner: { name: 'Начальный', color: '#28a745', icon: '🟢' },
  intermediate: { name: 'Средний', color: '#fd7e14', icon: '🟠' },
  advanced: { name: 'Продвинутый', color: '#dc3545', icon: '🔴' }
}

export const TOPICS = {
  output: { name: 'Вывод данных', order: 1 },
  input: { name: 'Ввод данных', order: 2 },
  conditions: { name: 'Условия', order: 3 },
  loops: { name: 'Циклы', order: 4 },
  arrays: { name: 'Массивы', order: 5 },
  functions: { name: 'Функции', order: 6 },
  oop: { name: 'ООП', order: 7 },
  algorithms: { name: 'Алгоритмы', order: 8 }
}

export const TASK_DATA = {
  python: {
    output: [
            {
        id: 1,
        title: "Hello World",
        description: "Напишите программу, которая выводит 'Hello World' в консоль.",
        difficulty: "beginner",
        defaultCode: '# Ваш код здесь\nprint("Hello World")',
        testPattern: /print\s*\(\s*["']Hello World["']\s*\)\s*$/m, // Только Hello World
        solution: 'print("Hello World")'
      },
      {
        id: 2,
        title: "Вывод чисел",
        description: "Выведите числа от 1 до 5, каждое на новой строке.",
        difficulty: "beginner", 
        defaultCode: '# Ваш код здесь\nprint(1)\nprint(2)\nprint(3)\nprint(4)\nprint(5)',
        testPattern: /print\s*\(\s*1\s*\)[\s\S]*print\s*\(\s*2\s*\)[\s\S]*print\s*\(\s*3\s*\)[\s\S]*print\s*\(\s*4\s*\)[\s\S]*print\s*\(\s*5\s*\)/,
        solution: 'for i in range(1, 6):\n    print(i)'
      }
    ],
    input: [
      {
        id: 3,
        title: "Сумма двух чисел",
        description: "Программа должна принимать два числа и выводить их сумму.",
        difficulty: "beginner",
        defaultCode: '# Ваш код здесь\na = int(input())\nb = int(input())',
        testPattern: /input\s*\(\s*\)[\s\S]*\+\s*/,
        solution: 'a = int(input())\nb = int(input())\nprint(a + b)'
      }
    ],
    conditions: [
      {
        id: 4, 
        title: "Проверка четности",
        description: "Напишите программу, которая проверяет, является ли число четным.",
        difficulty: "beginner",
        defaultCode: '# Ваш код здесь\nnum = int(input())',
        testPattern: /if\s+.*%\s*2\s*==\s*0/,
        solution: 'num = int(input())\nif num % 2 == 0:\n    print("Четное")\nelse:\n    print("Нечетное")'
      }
    ],
    loops: [
      {
        id: 5,
        title: "Таблица умножения",
        description: "Выведите таблицу умножения для числа 5 (от 1 до 10).",
        difficulty: "intermediate", 
        defaultCode: '# Ваш код здесь',
        testPattern: /for\s+.*in\s+range.*print.*\*/,
        solution: 'for i in range(1, 11):\n    print(f"5 * {i} = {5 * i}")'
      }
    ]
  },
  javascript: {
    output: [
      {
        id: 1,
        title: "Hello World", 
        description: "Напишите программу, которая выводит 'Hello World' в консоль.",
        difficulty: "beginner",
        defaultCode: '// Ваш код здесь\nconsole.log("Hello World")',
        testPattern: /console\.log\s*\(\s*["']Hello World["']\s*\)/,
        solution: 'console.log("Hello World")'
      }
    ],
    input: [
      {
        id: 2,
        title: "Сумма двух чисел",
        description: "Программа должна принимать два числа и выводить их сумму.",
        difficulty: "beginner",
        defaultCode: '// Ваш код здесь\nconst a = parseInt(prompt());\nconst b = parseInt(prompt());',
        testPattern: /prompt\s*\(\s*\)[\s\S]*\+\s*/,
        solution: 'const a = parseInt(prompt());\nconst b = parseInt(prompt());\nconsole.log(a + b);'
      }
    ]
  },
  cpp: {
    output: [
      {
        id: 1,
        title: "Hello World",
        description: "Напишите программу, которая выводит 'Hello World' в консоль.",
        difficulty: "beginner",
        defaultCode: '// Ваш код здесь\n#include <iostream>\n\nint main() {\n    std::cout << "Hello World" << std::endl;\n    return 0;\n}',
        testPattern: /cout\s*<<\s*["']Hello World["']/,
        solution: '#include <iostream>\n\nint main() {\n    std::cout << "Hello World" << std::endl;\n    return 0;\n}'
      }
    ]
  },
  java: {
    output: [
      {
        id: 1,
        title: "Hello World",
        description: "Напишите программу, которая выводит 'Hello World' в консоль.", 
        difficulty: "beginner",
        defaultCode: '// Ваш код здесь\npublic class Main {\n    public static void main(String[] args) {\n        System.out.println("Hello World");\n    }\n}',
        testPattern: /System\.out\.println\s*\(\s*["']Hello World["']\s*\)/,
        solution: 'public class Main {\n    public static void main(String[] args) {\n        System.out.println("Hello World");\n    }\n}'
      }
    ]
  }
}

// Вспомогательные функции для работы с задачами
export const getTasksByLanguage = (language) => {
  return TASK_DATA[language] || {}
}

export const getTask = (language, topic, taskId) => {
  const tasks = TASK_DATA[language]?.[topic] || []
  return tasks.find(task => task.id === parseInt(taskId))
}

export const getAllTopics = () => {
  return Object.entries(TOPICS).sort((a, b) => a[1].order - b[1].order)
}