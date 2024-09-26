// apollo.config.js
module.exports = {
  client: {
    service: {
      name: 'client',
        // URL to the GraphQL API
      url: 'http://localhost:5148/graphql',
    },
      // Files processed by the extension
    includes: [
      'src/**/*.vue',
      'src/**/*.js',
    ],
  },
}