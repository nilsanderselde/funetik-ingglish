var extract_text = require("extract-text-webpack-plugin")

module.exports = {
    entry: ["./scss/main.scss"], 
    output: {filename: "output.js"},
    plugins: [
        new extract_text({
            filename: "static/stylesheets/main.css",
            allChunks: true,
        })],
    module: {
        rules: [{
            test: /\.css$/,
            use: extract_text.extract({use: [{
                loader: "css-loader",
                options: {url: false},
            },{
                loader: "postcss-loader",
                options: {plugins: (loader) => [
                    require("autoprefixer")(),
                    require("postcss-clean")(),
                ]}
            }]})
        },{
            test: /\.scss$/,
            use: extract_text.extract({use: [{
                loader: "css-loader",
                options: {url: false},
            },{
                loader: "postcss-loader",
                options: {plugins: (loader) => [
                    require("autoprefixer")(),
                    require("postcss-clean")(),
                ]}
            },{
                loader: "sass-loader",
                options: {url: false},
            }]})

        }]
    }
}