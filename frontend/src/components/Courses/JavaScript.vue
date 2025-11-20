<template>
  <CourseLayout
    course-title="Курс JavaScript"
    :lessons="javascriptLessons"
    language="javascript"
  />
</template>

<script>
import CourseLayout from './CourseLayout.vue';

export default {
  name: 'JavaScript',
  components: {
    CourseLayout
  },
  data() {
    return {
      javascriptLessons: [
        {
          id: 1,
          title: 'Hello World',
          description: 'Напишите программу которая выводит "Hello, World!" на экран.',
          difficulty: 'easy',
          completed: false,
          tests: [
            { input: '', expected: 'Hello, World!', actual: null, status: null, error: null }
          ],
          starterCode: 'console.log("Hello, World!");'
        },
        {
          id: 2,
          title: 'Сумма с вводом',
          description: 'Напишите программу которая читает два числа из stdin и выводит их сумму.',
          difficulty: 'easy', 
          completed: false,
          tests: [
            { input: '5\n3', expected: '8', actual: null, status: null, error: null }
          ],
          starterCode: `// Рабочий способ чтения stdin для Windows
        const readline = require('readline');

        const rl = readline.createInterface({
          input: process.stdin,
          output: process.stdout
        });

        let numbers = [];

        rl.on('line', (line) => {
          numbers.push(parseInt(line.trim()));
          if (numbers.length === 2) {
            const sum = numbers[0] + numbers[1];
            console.log(sum);
            rl.close();
          }
        });`
        }
      ]
    }
  }
}
</script>