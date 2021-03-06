const { defineConfig } = require("cypress")

module.exports = defineConfig({
    projectId: "j29dep",
    video: false,
    e2e: {
        // We've imported your old cypress plugins here.
        // You may want to clean this up later by importing these.
        setupNodeEvents(on, config) {
            return require("./cypress/plugins/index.cjs")(on, config)
        },
        baseUrl: "http://localhost:3000",
        excludeSpecPattern: "**/examples/**",
    },
    pageLoadTimeout: 10000,
})
