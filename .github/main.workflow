workflow "run tests on push" {
  on = "push"
  resolves = "run tests"
}

workflow "run tests weekly" {
  on = "schedule(0 0 * * 0)"
  resolves = "run tests"
}

action "run tests" {
  uses = "cedrickring/golang-action@1.3.0"
  secrets = [
    "GOOGLE_API_KEY"
  ]
}