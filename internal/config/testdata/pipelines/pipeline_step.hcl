project = "foo"

pipeline "foo" {
  step {
    image_url = "example.com/test"

    use "test" {
      foo = "bar"
    }
  }
}

app "foo" {
}
