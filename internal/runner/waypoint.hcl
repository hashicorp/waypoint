
project = "p_test"
app "a_test" {
  build {
    use "docker" {
    }
  }

  deploy {
    use "docker" {
    }
  }
}
