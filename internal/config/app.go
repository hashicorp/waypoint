package config

// components returns all our components in a map form.
func (app *App) components() map[string]*Component {
	return map[string]*Component{
		"build":    app.Build,
		"registry": app.Registry,
		"deploy":   app.Platform,
		"release":  app.Release,
	}
}
