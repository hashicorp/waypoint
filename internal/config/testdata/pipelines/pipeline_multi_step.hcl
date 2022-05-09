project = "foo"

pipeline "foo" {
  step {
    image_url = "example.com/test"
    name      = "zero"

    use "test" {
      foo = "qubit"
    }
  }

  step {
    image_url = "example.com/second"

    use "test/exec" {
      foo = "few"
      bar = "bar"
    }
  }

  step {
    image_url = "example.com/different"
    depends_on = ["zero"]

    use "test/hunger" {
      foo = "food"
      bar = "drink"
    }
  }
}

app "web" {
    config {
        env = {
            static = "hello"
        }
    }

    build {}

    deploy {}
}
