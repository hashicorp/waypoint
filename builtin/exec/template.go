package exec

// tplData is the structure given to Go's text/template when rendering
// templates.
type tplData struct {
	// Input comes from the function input.
	Input map[string]interface{}
}

func (d *tplData) Populate(input *Input) {
	if input == nil {
		return
	}

	if d.Input == nil {
		d.Input = map[string]interface{}{}
	}

	for k, value := range input.Data {
		d.Input[k] = inputValueToInterface(value)
	}
}

func inputValueToInterface(v *Input_Value) interface{} {
	switch v := v.Value.(type) {
	case *Input_Value_Text:
		return v.Text

	default:
		return nil
	}
}
