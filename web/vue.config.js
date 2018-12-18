module.exports = {
    baseUrl: "/",

    // where to output built files
    outputDir: "dist",
    devServer: {
        proxy: {
            "/*": {
                target: "http://api:5050",
                ws: true
            }
        }
    }
};
