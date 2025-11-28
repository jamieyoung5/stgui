package widgets

import "fmt"

type Button struct {
	Label    string
	Callback func()
}

func NewButton(label string, callback func()) *Button {
	return &Button{
		Label:    label,
		Callback: callback,
	}
}

func (b *Button) RenderLines() []string {
	return []string{fmt.Sprintf("[ %s ]", b.Label)}
}

func (b *Button) OnClick() {
	if b.Callback != nil {
		b.Callback()
	}
}
