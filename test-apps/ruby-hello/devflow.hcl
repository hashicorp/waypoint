app "ruby-hello" {
  build "lambda" {
    runtime = "ruby2.5"
  }

  deploy "lambda" {
    bucket = "emp-devflow-lambda-test"
  }
}
