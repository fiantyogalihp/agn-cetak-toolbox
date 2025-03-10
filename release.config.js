const config = {
  branches: ["main"],
  plugins: [
    "@semantic-release/commit-analyzer",
    "@semantic-release/release-notes-generator",
    // ["@semantic-release/git", {
    //   "assets": ["build/**"],
    //   // "message": "chore(release): ${nextRelease.version} [skip ci]\n\n${nextRelease.notes}"
    // }],
    ["@semantic-release/github", {
      "assets": [
        { "path": "build/agn-cetak-toolbox.zip" },
      ]
    }],
  ]
};

module.exports = config;
