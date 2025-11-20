const state = {
  isLoggedIn: false,
  username: '',
  userAvatar: ''
}

const mutations = {
  setUser(state, userData) {
    state.isLoggedIn = userData.isLoggedIn
    state.username = userData.username
    state.userAvatar = userData.userAvatar
  },
  logout(state) {
    state.isLoggedIn = false
    state.username = ''
    state.userAvatar = ''
  }
}

const actions = {
  login({ commit }, userData) {
    commit('setUser', userData)
    localStorage.setItem('user', JSON.stringify(userData))
  },
  logout({ commit }) {
    commit('logout')
    localStorage.removeItem('user')
  },
  initialize({ commit }) {
    const savedUser = localStorage.getItem('user')
    if (savedUser) {
      try {
        const userData = JSON.parse(savedUser)
        commit('setUser', userData)
      } catch (e) {
        localStorage.removeItem('user')
      }
    }
  }
}

export default {
  state,
  mutations,
  actions,
  namespaced: true
}