module.exports = {
    devServer: {
        proxy: {
            "/auth": {
                target: "http://api:5050"
            }
        }
    }
};
