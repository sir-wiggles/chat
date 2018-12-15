module.exports = {
    dev: {
        proxyTable: {
            "/*": {
                target: "http://localhost:5050",
                changeOrigin: true
            }
        }
    }
};
