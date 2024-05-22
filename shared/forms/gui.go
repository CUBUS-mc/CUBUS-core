package forms

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

func refreshForm(form *Form, box *fyne.Container) {
	FormToFyneForm(form, box)
}

func fieldsToFyneForm(fields []Field, form *Form, box *fyne.Container) *widget.Form {
	fyneForm := widget.NewForm()

	for _, field := range fields { // TODO: Fix the issue that only one field is displayed
		switch field := field.(type) {
		case *FieldBaseType:
			// Do nothing
		case *TextField:
			entry := widget.NewEntry()
			entry.SetText(field.GetValue())
			entry.SetPlaceHolder(field.GetPlaceholder())
			entry.OnChanged = func(text string) {
				field.SetValue(text)
				refreshForm(form, box)
			}
			fyneForm.Append(field.GetId(), entry)
		case *MultipleChoiceField: // TODO: Use label instead of key for the options
			options := make([]string, 0, len(field.Options))
			for key := range field.Options {
				options = append(options, key)
			}
			selectWidget := widget.NewSelect(options, func(value string) { // TODO: translate label back to key
				field.SetValue(value)
			})
			selectWidget.SetSelected(field.GetValue())
			selectWidget.OnChanged = func(value string) {
				field.SetValue(value)
				refreshForm(form, box)
			}
			fyneForm.Append(field.GetId(), selectWidget)
		case *Message:
			label := widget.NewLabel(field.GetValue())
			fyneForm.Append(field.GetId(), label)
		case *NumberField:
			entry := widget.NewEntry()
			entry.SetText(field.GetValue())
			entry.SetPlaceHolder(field.GetPlaceholder())
			entry.OnChanged = func(text string) {
				field.SetValue(text)
				refreshForm(form, box)
			}
			fyneForm.Append(field.GetId(), entry)
		case *FieldGroup:
			if field.GetValue() != "" {
				label := widget.NewLabel(field.GetValue())
				fyneForm.Append(field.GetId(), label)
			}
			subForm := fieldsToFyneForm(field.GetFieldsToDisplay(), form, box)
			fyneForm.Append(field.GetId(), subForm)
		default:
			panic("Unknown field type")
		}
	}

	return fyneForm
}

func FormToFyneForm(form *Form, box *fyne.Container) {
	fields := form.GetFieldsToDisplay()
	box.RemoveAll()
	box.Add(fieldsToFyneForm(fields, form, box))
	box.Refresh()
}
