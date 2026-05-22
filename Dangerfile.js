const { danger, fail, warn, message } = require("danger");

const pr = danger.github.pr;

if (pr.body.length < 20) {
  warn("PR description is too short.");
}

const modifiedFiles = danger.git.modified_files;
const createdFiles = danger.git.created_files;

const allFiles = [...modifiedFiles, ...createdFiles];

if (allFiles.includes("go.mod") && !allFiles.includes("go.sum")) {
  fail("go.sum must be updated with go.mod changes.");
}

if (allFiles.some((f) => f.endsWith(".go")) &&
    !allFiles.some((f) => f.includes("_test.go"))) {
  warn("Go files changed without adding tests.");
}

if (pr.additions > 1000) {
  warn("Large PR detected. Consider splitting it.");
}

message("Thanks for contributing to ebook-lcp-server.");
