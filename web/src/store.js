import Vue from "vue";
import Vuex from "vuex";
import axios from "axios";

Vue.use(Vuex);

const systemMessages = {
    initialize: true,
    system: true
};

export default new Vuex.Store({
    strict: true,
    state: {
        socket: {
            message: null,
            connect: null,
            error: null,
            onerror: null,
            onopen: null,
            onclose: null
        },
        messages: [],
        author: {
            uuid: "",
            name: "",
            email: "",
            avatar: ""
        },
        auth: {
            token: window.localStorage.getItem("token"),
            valid: false
        }
    },
    getters: {
        LOGGED_IN: function(state) {
            return state.auth.valid;
        },
        AUTH_TOKEN: function(state) {
            return state.auth.token;
        }
    },
    mutations: {
        SOCKET_ONOPEN: function(state) {
            state.socket.open = true;
        },
        SOCKET_ONMESSAGE: function(state, message) {
            // initialize messages give us the client information (id and name); however, we
            // don't want the id to persist in this case because if the second message is from
            // the author it will then be grouped with the system message with the last if
            // statement in this block
            if (message.type === "initialize") {
                state.author = Object.assign({}, message.author);
                message.author.id = "";
            }

            if (state.messages.length === 0) {
                state.messages.push(message);
                return;
            }

            let lastMessage = state.messages[state.messages.length - 1];

            // if both the current message and the previous message are system related messages
            // then group them together
            if (
                systemMessages[message.type] &&
                systemMessages[lastMessage.type]
            ) {
                lastMessage.text.push(message.text);
                return;
            }

            if (lastMessage.author.id === message.author.id) {
                lastMessage.text.push(message.text);
                return;
            }

            state.messages.push(message);
        },
        SOCKET_ONERROR: function(state, error) {
            state.socket.onerror = error;
        },
        SOCKET_CONNECT: function(state) {
            state.socket.connect = true;
        },
        SOCKET_ERROR: function(state, error) {
            state.socket.error = error;
        },
        SOCKET_ONCLOSE: function(state) {
            state.socket.onclose = true;
        },
        OAUTH2_SET_TOKEN: function(state, { token }) {
            if (token) {
                state.auth.valid = true;
                state.auth.token = `Bearer ${token}`;
                window.localStorage.setItem("token", state.auth.token);
            } else {
                state.auth.valid = false;
                state.auth.token = "";
                window.localStorage.removeItem("token", state.auth.token);
            }
        },
        OAUTH2_TOKEN_STATE: function(state, valid) {
            state.auth.valid = valid;
        }
    },
    actions: {
        OAUTH2_GET_TOKEN: function({ commit }, code) {
            this.$log.debug("calling api server to exchange token");
            return axios
                .post("/auth/google", {
                    code: code,
                    redirect_uri: "postmessage"
                })
                .then(({ data }) => {
                    commit("OAUTH2_SET_TOKEN", data);
                })
                .catch(error => {
                    this.$log.error(`OAUTH2_GET_TOKEN ${error}`);
                });
        },
        CHECK_TOKEN: function({ commit, getters }) {
            return axios
                .get("/api/health", {
                    headers: {
                        Authorization: getters.AUTH_TOKEN
                    }
                })
                .then(() => {
                    commit("OAUTH2_TOKEN_STATE", true);
                    return true;
                })
                .catch(() => {
                    commit("OAUTH2_SET_TOKEN", "");
                    return false;
                });
        }
    }
});
