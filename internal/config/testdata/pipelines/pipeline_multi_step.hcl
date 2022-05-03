project = "foo"

pipeline "foo" {
  step {
    image_url = "example.com/test"

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

    use "test/hunger" {
      foo = "food"
      bar = "drink"
    }
  }
}

app "foo" {
}
