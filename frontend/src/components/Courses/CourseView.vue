<template>
  <CourseLayout
    :course-title="courseTitle"
    :language="courseLanguage"
  />
</template>

<script>
import CourseLayout from './CourseLayout.vue';

export default {
  name: 'CourseView',
  components: {
    CourseLayout
  },
  props: {
    language: {
      type: String,
      required: true
    }
  },
  computed: {
    courseLanguage() {
      const langMap = {
        'python': 'python',
        'javascript': 'javascript',
        'js': 'javascript',
        'typescript': 'typescript',
        'ts': 'typescript',
        'java': 'java',
        'cpp': 'cpp',
        'c++': 'cpp',
        'csharp': 'csharp',
        'c#': 'csharp',
        'cs': 'csharp',
        'go': 'go',
        'golang': 'go',
        'rust': 'rust',
        'php': 'php',
        'ruby': 'ruby'
      }
      return langMap[this.language.toLowerCase()] || this.language.toLowerCase()
    },
    
    courseTitle() {
      const titles = {
        'python': 'Курс Python для начинающих',
        'javascript': 'Курс JavaScript',
        'typescript': 'Курс TypeScript',
        'java': 'Курс Java',
        'cpp': 'Курс C++',
        'csharp': 'Курс C#',
        'go': 'Курс Go',
        'rust': 'Курс Rust',
        'php': 'Курс PHP',
        'ruby': 'Курс Ruby'
      }
      return titles[this.courseLanguage] || `Курс ${this.courseLanguage}`
    }
  },
  created() {
    console.log(`Курс: ${this.language} -> ${this.courseLanguage}`)
  }
}
</script>