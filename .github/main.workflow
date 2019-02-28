workflow "Publish container" {
  on       = "push"
  resolves = ["Test"]
}

action "Lint" {
  uses = "./.github/actions/golang"
  args = "fmt"
}

action "Test" {
  needs = ["Lint"]
  uses  = "./.github/actions/golang"
  args  = "test"
  env = {
    GO111MODULE = "on"
  }
}
