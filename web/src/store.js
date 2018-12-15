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
            token: ""
        }
    },
    getters: {
        LOGGED_IN: function(state) {
            return state.auth.token.length > 0;
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
            state.auth.token = `Bearer ${token}`;
        }
    },
    actions: {
        OAUTH2_GET_TOKEN: function({ commit }, code) {
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
        }
    }
});
