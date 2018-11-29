const mutations = {
  SOCKET_CONNECT: "",
  SOCKET_ERROR: "",
  SOCKET_ONERROR: "",
  SOCKET_ONMESSAGE: "",
  SOCKET_ONOPEN: "",
  TOGGLE_TODO_DONE: ""
};

for (let mutation of Object.keys(mutations)) {
  mutations[mutation] = mutations[mutation].toUpperCase();
}

export { mutations };
