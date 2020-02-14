package status

type Updater interface {
	Update(str string)
	Close()
}
