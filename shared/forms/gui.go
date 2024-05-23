package forms

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func refreshForm(form *Form, box *fyne.Container, dialog *dialog.CustomDialog, window fyne.Window) {
	FormToFyneForm(form, box, dialog, window)
}

func fieldsToFyneForm(fields []Field, form *Form, box *fyne.Container, fyneForm *widget.Form, dialog *dialog.CustomDialog, window fyne.Window) {
	for _, field := range fields {
		switch field := field.(type) {
		case *FieldBaseType:
			// Do nothing
		case *TextField:
			entry := widget.NewEntry()
			entry.SetText(field.GetValue())
			entry.SetPlaceHolder(field.GetPlaceholder())
			entry.OnChanged = func(text string) {
				field.SetValue(text)
				refreshForm(form, box, dialog, window)
			}
			entry.Validator = func(text string) error {
				if !field.IsValid() {
					return field.GetError()
				}
				return nil
			}
			fyneForm.Append(field.GetPrompt(), entry)
		case *MultipleChoiceField:
			labelsToKeys := make(map[string]string)
			options := make([]string, 0, len(field.GetOptions()))
			for key, option := range field.GetOptions() {
				options = append(options, option.Label)
				labelsToKeys[option.Label] = key
			}
			selectWidget := widget.NewSelect(options, func(value string) {
				key := labelsToKeys[value]
				field.SetValue(key)
			})
			selectWidget.SetSelected(field.Options[field.GetValue()].Label)
			selectWidget.OnChanged = func(value string) {
				key := labelsToKeys[value]
				field.SetValue(key)
				refreshForm(form, box, dialog, window)
			}
			fyneForm.Append(field.GetPrompt(), selectWidget)
		case *Message:
			fyneForm.Append(field.GetValue(), widget.NewLabel(""))
		case *NumberField:
			entry := widget.NewEntry()
			entry.SetText(field.GetValue())
			entry.SetPlaceHolder(field.GetPlaceholder())
			entry.OnChanged = func(text string) {
				field.SetValue(text)
				refreshForm(form, box, dialog, window)
			}
			entry.Validator = func(text string) error {
				if !field.IsValid() {
					return field.GetError()
				}
				return nil
			}
			fyneForm.Append(field.GetId(), entry)
		case *FieldGroup:
			if field.GetValue() != "" {
				fyneForm.Append(field.GetValue(), widget.NewLabel(""))
			}
			fieldsToFyneForm(field.GetFieldsToDisplay(), form, box, fyneForm, dialog, window)
		default:
			panic("Unknown field type")
		}
	}
}

func FormToFyneForm(form *Form, box *fyne.Container, parentDialog *dialog.CustomDialog, window fyne.Window) {
	fyneForm := widget.NewForm()
	fields := form.GetFieldsToDisplay()
	fieldsToFyneForm(fields, form, box, fyneForm, parentDialog, window)
	fyneForm.OnSubmit = func() {
		if form.IsValid() {
			parentDialog.Hide()
		} else {
			dialog.ShowError(form.GetError(), window)
		}
	}
	fyneForm.OnCancel = func() {
		parentDialog.Hide()
	}
	fyneForm.Resize(fyne.NewSize(700, 400))
	box.RemoveAll()
	box.Add(fyneForm)
	box.Refresh()
}
