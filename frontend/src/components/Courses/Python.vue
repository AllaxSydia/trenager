<template>
  <CourseLayout
    course-title="Курс Python"
    :lessons="pythonLessons"
    language="python"
  />
</template>

<script>
import CourseLayout from './CourseLayout.vue';

export default {
  name: 'Python',
  components: {
    CourseLayout
  },
  data() {
    return {
      pythonLessons: [
        {
          id: 1,
          title: 'Привет, Python!',
          description: 'Напишите программу, которая выводит "Hello, World!" на экран.',
          difficulty: 'easy',
          completed: true,
          tests: [
            { input: '', expected: 'Hello, World!', actual: null, status: null, error: null }
          ],
          starterCode: 'print("Hello, World!")'
        },
        {
          id: 2,
          title: 'Сумма двух чисел',
          description: 'Напишите программу, которая складывает два числа (5 и 3).',
          difficulty: 'easy',
          completed: false,
          tests: [
            { input: '', expected: '8', actual: null, status: null, error: null }
          ],
          starterCode: '# Сложите два числа\nnum1 = 5\nnum2 = 3\nresult = num1 + num2\nprint(result)'
        },
        {
          id: 3,
          title: 'Факториал числа',
          description: 'Напишите программу, которая вычисляет факториал числа 5.',
          difficulty: 'medium',
          completed: false,
          tests: [
            { input: '', expected: '120', actual: null, status: null, error: null }
          ],
          starterCode: '# Вычислите факториал числа 5\nn = 5\nfactorial = 1\nfor i in range(1, n + 1):\n    factorial *= i\nprint(factorial)'
        },
        {
          id: 4,
          title: 'Проверка на четность',
          description: 'Напишите программу, которая проверяет, является ли число 7 четным.',
          difficulty: 'easy',
          completed: false,
          tests: [
            { input: '', expected: 'Нечетное', actual: null, status: null, error: null }
          ],
          starterCode: '# Проверка на четность\nnumber = 7\nif number % 2 == 0:\n    print("Четное")\nelse:\n    print("Нечетное")'
        },
        {
          id: 5,
          title: 'Поиск максимума',
          description: 'Напишите программу, которая находит максимальное из трех чисел (5, 3, 8).',
          difficulty: 'medium',
          completed: false,
          tests: [
            { input: '', expected: '8', actual: null, status: null, error: null }
          ],
          starterCode: '# Найдите максимальное из трех чисел\nnum1, num2, num3 = 5, 3, 8\nmax_num = max(num1, num2, num3)\nprint(max_num)'
        }
      ]
    }
  }
}
</script>