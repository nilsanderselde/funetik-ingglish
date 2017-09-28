var extract_text = require("extract-text-webpack-plugin")

module.exports = {
    entry: ["./scss/main.scss"], 
    plugins: [
        new extract_text({
            filename: "static/stylesheets/main.css"
        })
    ]
}