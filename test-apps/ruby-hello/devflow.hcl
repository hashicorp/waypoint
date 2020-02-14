app "ruby-hello" {
  build "lambda" {
    runtime = "ruby2.5"
    setup = "yum install -y postgresql-devel"
  }

  deploy "lambda" {
    bucket = "emp-devflow-lambda-test"
  }
}
