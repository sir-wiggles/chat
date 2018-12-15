var googleAuth = (function() {
    function install() {
        var apiUrl = "https://apis.google.com/js/api.js";
        return new Promise(resolve => {
            var script = document.createElement("script");
            script.src = apiUrl;
            script.onreadystatechange = script.onload = function() {
                if (
                    !script.readyState ||
                    /loaded|compvare/.test(script.readyState)
                ) {
                    setTimeout(function() {
                        resolve();
                    }, 500);
                }
            };
            document.getElementsByTagName("head")[0].appendChild(script);
        });
    }

    function init(config) {
        return new Promise(resolve => {
            window.gapi.load("auth2", () => {
                resolve(window.gapi);
                //window.gapi.auth2.init(config).then(() => {
                //    resolve(window.gapi);
                //});
            });
        });
    }

    function Auth() {
        if (!(this instanceof Auth)) return new Auth();
        this.GoogleAuth = null;
        this.isAuthorized = false;
        this.store = null;

        this.login = (config, prompt) => {
            config.prompt = prompt;
            config.response_type = "code";
            return new Promise((resolve, reject) => {
                console.log("googleAuth", this.GoogleAuth);
                this.GoogleAuth.authorize(config, function(response) {
                    if (response.error) {
                        reject(response.error);
                    }
                    resolve(response);
                });
            });
        };

        this.waitForGapi = () => {
            return new Promise(resolve => {
                setTimeout(() => {
                    if (window.gapi.auth2) {
                        resolve(window.gapi.auth2);
                    } else {
                        resolve(this.waitForGapi());
                    }
                }, 250);
            });
        };

        this.load = config => {
            this.store = config.store;
            this.waitForGapi()
                .then(gapi => {
                    console.log("gapi", gapi);
                    this.GoogleAuth = gapi;
                })
                .then(() => {
                    return this.login(config, "none");
                })
                .catch(error => {
                    console.error(error);
                    return this.login(config, "");
                })
                .then(response => {
                    console.log(response);
                    this.store.dispatch("OAUTH2_GET_TOKEN", response.code);
                });

            //.then(gapi => {
            //    gapi.auth2.authorize({});
            //    this.GoogleAuth = gapi.auth2.getAuthInstance();
            //    this.GoogleAuth.isSignedIn.listen(this.isSignedInListen);
            //    return this.GoogleAuth.isSignedIn.get();
            //})
            //.then(isSignedIn => {
            //    this.store.commit("OAUTH2_SET_SIGNED_IN", isSignedIn);
            //    if (isSignedIn) {
            //        return this.GoogleAuth.grantOfflineAccess();
            //    } else {
            //        return this.GoogleAuth.grantOfflineAccess();
            //    }
            //})
            //.then(code => {
            //    this.store.dispatch("OAUTH2_GET_TOKEN", code);
            //})
            //.then(() => {});
        };

        this.isSignedInListen = event => {
            console.log(event);
            this.store["commit"]("OAUTH2_SET_SIGNED_IN", event);
        };

        this.getAuthCode = (successCallback, errorCallback, prompt) => {
            return new Promise((resolve, reject) => {
                if (!this.googleAuthReady(errorCallback, reject)) return;
                //this.GoogleAuth.grantOfflineAccess({ prompt: "select_account" })
                this.GoogleAuth.grantOfflineAccess({})
                    .then(function(resp) {
                        if (typeof successCallback === "function")
                            successCallback(resp.code);
                        resolve(resp.code);
                    })
                    .catch(function(error) {
                        if (typeof errorCallback === "function")
                            errorCallback(error);
                        reject(error);
                    });
            });
        };
    }

    return new Auth();
})();

function installGoogleAuthPlugin(Vue, options) {
    //set config
    var GoogleAuthConfig = null;
    var GoogleAuthDefaultConfig = {
        scope: "profile email https://www.googleapis.com/auth/plus.login",
        discoveryDocs: [
            "https://www.googleapis.com/discovery/v1/apis/drive/v3/rest"
        ]
    };
    if (typeof options === "object") {
        GoogleAuthConfig = Object.assign(GoogleAuthDefaultConfig, options);
        if (!options.clientId) {
            /* eslint-disable */
            console.warn("clientId is required");
        }
    } else {
        console.warn("invalid option type. Object type accepted only");
    }

    //Install Vue plugin
    Vue.gAuth = googleAuth;
    Object.defineProperties(Vue.prototype, {
        $gAuth: {
            get: function() {
                return Vue.gAuth;
            }
        }
    });
    Vue.gAuth.load(GoogleAuthConfig);
}

export default installGoogleAuthPlugin;
