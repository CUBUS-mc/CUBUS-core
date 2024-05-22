package forms

import (
	"net"
	"strconv"
)

// TODO: Fix all validators to use field any instead of value string

type Field interface {
	GetId() string
	ShouldDisplay() bool
	IsValid() bool
	GetValue() string
	SetValue(value string)
}

type Validator interface {
	Validate(field any) bool
}

type DisplayCondition interface {
	DisplayCondition(field any) bool
}

// Defining the Base Field Type

type FieldBaseType struct {
	Id                string
	DisplayConditions []DisplayCondition
	Validators        []Validator
	Value             string
	form              *Form
}

func (f *FieldBaseType) GetId() string {
	return f.Id
}

func (f *FieldBaseType) ShouldDisplay() bool {
	for _, displayCondition := range f.DisplayConditions {
		if !displayCondition.DisplayCondition(f) {
			return false
		}
	}
	return true
}

func (f *FieldBaseType) IsValid() bool {
	for _, validator := range f.Validators {
		if !validator.Validate(f) {
			return false
		}
	}
	return true
}

func (f *FieldBaseType) GetValue() string {
	return f.Value
}

func (f *FieldBaseType) SetValue(value string) {
	f.Value = value
	if f.form != nil {
		f.form.onChange()
	}
}

type CustomValidator struct {
	Validator func(field any) bool
}

func (v *CustomValidator) Validate(field any) bool {
	return v.Validator(field)
}

type AllFieldsValid struct{}

func (v *AllFieldsValid) Validate(field any) bool {
	fields := field.(*FieldBaseType).form.Fields
	for _, f := range fields {
		if !f.IsValid() {
			return false
		}
	}
	return true
}

type IsValidValidator struct {
	fieldIds []string
}

func (v *IsValidValidator) Validate(field any) bool {
	fields := field.(*FieldBaseType).form.Fields
	for _, f := range fields {
		for _, id := range v.fieldIds {
			if f.GetId() == id && !f.IsValid() {
				return false
			}
		}
	}
	return true
}

type AlwaysDisplay struct{}

func (d *AlwaysDisplay) DisplayCondition(_ any) bool {
	return true
}

type CustomDisplayCondition struct {
	Condition func(field any) bool
}

func (d *CustomDisplayCondition) DisplayCondition(field any) bool {
	return d.Condition(field)
}

type IsValidDisplayCondition struct {
	fieldIds []string
}

func (d *IsValidDisplayCondition) DisplayCondition(field any) bool {
	fields := field.(*FieldBaseType).form.Fields
	for _, f := range fields {
		for _, id := range d.fieldIds {
			if f.GetId() == id && !f.IsValid() {
				return false
			}
		}
	}
	return true
}

type IsInvalidDisplayCondition struct {
	fieldIds []string
}

func (d *IsInvalidDisplayCondition) DisplayCondition(field any) bool {
	fields := field.(*FieldBaseType).form.Fields
	for _, f := range fields {
		for _, id := range d.fieldIds {
			if f.GetId() == id && f.IsValid() {
				return false
			}
		}
	}
	return true
}

type AllFieldsValidDisplayCondition struct{}

func (d *AllFieldsValidDisplayCondition) DisplayCondition(field any) bool {
	fields := field.(*FieldBaseType).form.Fields
	for _, f := range fields {
		if !f.IsValid() {
			return false
		}
	}
	return true
}

type HasValueDisplayCondition struct {
	fieldId string
	value   string
}

func (d *HasValueDisplayCondition) DisplayCondition(field any) bool {
	fields := field.(*FieldBaseType).form.Fields
	for _, f := range fields {
		if f.GetId() == d.fieldId && f.GetValue() == d.value {
			return true
		}
	}
	return false

}

// Defining the Field Types based on the Base Field Type

type Message struct {
	*FieldBaseType
}

// Defining the Text Field Type based on the Base Field Type

type TextField struct {
	*FieldBaseType
	Placeholder string
	Prompt      string
}

type NotEmptyValidator struct{}

func (v *NotEmptyValidator) Validate(field any) bool {
	value := field.(*TextField).Value
	return value != ""
}

type MaxLengthValidator struct {
	MaxLength int
}

func (v *MaxLengthValidator) Validate(value string) bool {
	return len(value) <= v.MaxLength
}

type MinLengthValidator struct {
	MinLength int
}

func (v *MinLengthValidator) Validate(value string) bool {
	return v.MinLength <= len(value)
}

type IpValidator struct{}

func (v *IpValidator) Validate(field *TextField) bool {
	return net.ParseIP(field.Value) != nil
}

func (t *TextField) GetPlaceholder() string {
	return t.Placeholder
}

func (t *TextField) GetPrompt() string {
	return t.Prompt
}

// Defining the Number Field Type based on the Base Field Type

type NumberField struct {
	*TextField
}

type MinValidator struct {
	Min int
}

func (v *MinValidator) Validate(field *NumberField) bool {
	valueAsInt, err := strconv.Atoi(field.Value)
	if err != nil {
		return false
	}
	return v.Min <= valueAsInt
}

type MaxValidator struct {
	Max int
}

func (v *MaxValidator) Validate(field *NumberField) bool {
	valueAsInt, err := strconv.Atoi(field.Value)
	if err != nil {
		return false
	}
	return valueAsInt <= v.Max
}

type IsIntegerValidator struct{}

func (v *IsIntegerValidator) Validate(field *NumberField) bool {
	_, err := strconv.Atoi(field.Value)
	return err == nil
}

// Defining the Multiple Choice Field Type based on the Text Field Type

type MultipleChoiceField struct {
	*TextField
	Options map[string]Option
}

type Option struct {
	Label       string
	Description string
}

type ChoiceValidator struct{}

func (v *ChoiceValidator) Validate(field any) bool {
	multipleChoiceField, ok := field.(*MultipleChoiceField)
	if !ok {
		return false
	}
	_, ok = multipleChoiceField.Options[multipleChoiceField.Value]
	return ok
}

// Defining the Field Group Type based on the Base Field Type

type FieldGroup struct {
	*FieldBaseType
	Fields []Field
}

// Defining the Form Type

type Form struct {
	Fields   []Field
	onChange func()
}

func (f *Form) IsValid() bool {
	for _, field := range f.Fields {
		if !field.IsValid() {
			return false
		}
	}
	return true
}

func (f *Form) GetFieldById(id string) Field {
	for _, field := range f.Fields {
		if field.GetId() == id {
			return field
		}
	}
	return nil
}

func (f *Form) GetFieldsToDisplay() []Field {
	var fieldsToDisplay []Field
	for _, field := range f.Fields {
		if field.ShouldDisplay() {
			fieldsToDisplay = append(fieldsToDisplay, field)
		}
	}
	return fieldsToDisplay
}

func (f *Form) GetFieldValues() map[string]string {
	fieldValues := make(map[string]string)
	for _, field := range f.Fields {
		fieldValues[field.GetId()] = field.GetValue()
	}
	return fieldValues
}

func (f *Form) SetOnChangeCallback(onChange func()) {
	f.onChange = onChange
}

// Defining the Form Builder Functions

func NewForm(fields ...Field) *Form {
	form := &Form{Fields: fields, onChange: func() {}}
	for _, field := range fields {
		field.(*FieldBaseType).form = form
	}
	return form
}

func NewFieldGroup(id string, displayConditions []DisplayCondition, validators []Validator, heading string, fields ...Field) *FieldGroup {
	return &FieldGroup{FieldBaseType: &FieldBaseType{Id: id, DisplayConditions: displayConditions, Validators: validators, Value: heading}, Fields: fields}
}

func NewTextField(id string, displayConditions []DisplayCondition, validators []Validator, placeholder string, prompt string, defaultValue string) *TextField {
	return &TextField{FieldBaseType: &FieldBaseType{Id: id, DisplayConditions: displayConditions, Validators: validators, Value: defaultValue}, Placeholder: placeholder, Prompt: prompt}
}

func NewNumberField(id string, displayConditions []DisplayCondition, validators []Validator, placeholder string, prompt string, defaultValue int) *NumberField {
	return &NumberField{TextField: &TextField{FieldBaseType: &FieldBaseType{Id: id, DisplayConditions: displayConditions, Validators: validators, Value: strconv.Itoa(defaultValue)}, Placeholder: placeholder, Prompt: prompt}}
}

func NewMultipleChoiceField(id string, displayConditions []DisplayCondition, validators []Validator, placeholder string, prompt string, options map[string]Option, defaultValue string) *MultipleChoiceField {
	return &MultipleChoiceField{TextField: &TextField{FieldBaseType: &FieldBaseType{Id: id, DisplayConditions: displayConditions, Validators: validators, Value: defaultValue}, Placeholder: placeholder, Prompt: prompt}, Options: options}
}

func NewMessage(id string, displayConditions []DisplayCondition, validators []Validator, message string) *Message {
	return &Message{FieldBaseType: &FieldBaseType{Id: id, DisplayConditions: displayConditions, Validators: validators, Value: message}}
}
