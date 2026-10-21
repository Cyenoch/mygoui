package conversation

import (
	"slices"
	"strings"

	"github.com/Cyenoch/mygoui"
	"github.com/Cyenoch/mygoui/form"
	locales "github.com/Cyenoch/mygoui/locale"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/i18n"
	"github.com/egoist/mygo/ui"
)

type QuestionID string
type ChoiceID string

type Choice struct {
	Value              ChoiceID
	Label, Description string
	Disabled           bool
	Hint, Preview      string
}

type Question struct {
	Name                          QuestionID
	Header                        string
	Title, Description            string
	Choices                       []Choice
	Multiple, Optional, AllowText bool
	// TextExclusive makes single-choice and custom text mutually exclusive.
	// Multiple-choice questions always retain both.
	TextExclusive              bool
	TextLabel, TextPlaceholder string
}

type Answer struct {
	Values  []ChoiceID
	Text    string
	Skipped bool
}

type QuestionnaireState struct {
	Index                int
	Answers              map[QuestionID]Answer
	Error                string
	Submitted            bool
	single               ChoiceID
	multiple             map[ChoiceID]*bool
	freeform             string
	active               QuestionID
	initialized, skipped bool
	focus, questionFocus ui.Handle
	answerRequired       bool
	previews             map[ChoiceID]bool
}

type QuestionnaireOptions struct {
	Label       string
	Busy        bool
	SubmitLabel string
	Disabled    bool
	// MaxBodyHeight bounds only the question body; actions remain outside it.
	MaxBodyHeight float32
	OnSubmit      func(map[QuestionID]Answer) error
	OnChange      func(QuestionID, Answer)
	OnCancel      func()
	// KeepCancelEnabled lets a guarded asynchronous consumer explain uncertain
	// delivery without hiding the explicit cancellation action.
	KeepCancelEnabled bool
	LetterShortcuts   bool
}

func (s *QuestionnaireState) current(q Question) {
	if s.Answers == nil {
		s.Answers = make(map[QuestionID]Answer)
	}
	if s.initialized && s.active == q.Name {
		return
	}
	s.active, s.initialized = q.Name, true
	a := s.Answers[q.Name]
	s.single = ""
	if len(a.Values) > 0 {
		s.single = a.Values[0]
	}
	s.multiple = make(map[ChoiceID]*bool, len(q.Choices))
	for _, choice := range q.Choices {
		selected := slices.Contains(a.Values, choice.Value)
		s.multiple[choice.Value] = &selected
	}
	s.freeform, s.skipped, s.Error = a.Text, a.Skipped, ""
	s.previews = nil
	s.answerRequired = false
}

// normalizedAnswer uses the current question as authority. Removed or disabled
// choices and text from a question that no longer allows it cannot be submitted.
func normalizedAnswer(q Question, a Answer) Answer {
	out := Answer{Values: []ChoiceID{}, Skipped: a.Skipped && q.Optional}
	if out.Skipped {
		return out
	}
	if q.AllowText {
		out.Text = a.Text
		if q.TextExclusive && !q.Multiple && strings.TrimSpace(a.Text) != "" {
			return out
		}
	}
	for _, choice := range q.Choices {
		if choice.Disabled || !slices.Contains(a.Values, choice.Value) {
			continue
		}
		out.Values = append(out.Values, choice.Value)
		if !q.Multiple {
			break
		}
	}
	return out
}

func (s *QuestionnaireState) answer(q Question) Answer {
	a := Answer{Values: []ChoiceID{}, Text: s.freeform, Skipped: s.skipped}
	if q.Multiple {
		for _, choice := range q.Choices {
			if selected := s.multiple[choice.Value]; selected != nil && *selected {
				a.Values = append(a.Values, choice.Value)
			}
		}
	} else if s.single != "" {
		a.Values = append(a.Values, s.single)
	}
	return normalizedAnswer(q, a)
}

func answerValid(q Question, a Answer) bool {
	return q.Optional || len(a.Values) > 0 || q.AllowText && strings.TrimSpace(a.Text) != ""
}

func (s *QuestionnaireState) commit(q Question, a Answer, o QuestionnaireOptions) {
	s.Answers[q.Name] = a
	if o.OnChange != nil {
		a.Values = slices.Clone(a.Values)
		o.OnChange(q.Name, a)
	}
}

func (s *QuestionnaireState) Previous(questions []Question, o QuestionnaireOptions) {
	if o.Busy || o.Disabled || s.Index <= 0 || s.Index >= len(questions) {
		return
	}
	q := questions[s.Index]
	s.current(q)
	s.commit(q, s.answer(q), o)
	s.Index--
	s.initialized, s.Submitted = false, false
}

func (s *QuestionnaireState) Next(questions []Question, o QuestionnaireOptions) bool {
	if o.Busy || o.Disabled || s.Index < 0 || s.Index >= len(questions) {
		return false
	}
	q := questions[s.Index]
	s.current(q)
	a := s.answer(q)
	if !answerValid(q, a) {
		s.Error, s.answerRequired = locales.ComponentQuestionAnswerRequired(), true
		return false
	}
	s.commit(q, a, o)
	if s.Index+1 < len(questions) {
		s.Index++
		s.initialized = false
	}
	return true
}

func (s *QuestionnaireState) Skip(questions []Question, o QuestionnaireOptions) bool {
	if o.Busy || o.Disabled || s.Index < 0 || s.Index >= len(questions) || !questions[s.Index].Optional {
		return false
	}
	q := questions[s.Index]
	s.current(q)
	s.skipped = true
	s.commit(q, normalizedAnswer(q, Answer{Skipped: true}), o)
	if s.Index+1 < len(questions) {
		s.Index++
		s.initialized = false
	}
	return true
}

func (s *QuestionnaireState) submit(questions []Question, o QuestionnaireOptions) {
	if o.Busy || o.Disabled || o.OnSubmit == nil || !s.Next(questions, o) {
		return
	}
	result := make(map[QuestionID]Answer, len(questions))
	for i, q := range questions {
		a := normalizedAnswer(q, s.Answers[q.Name])
		if !answerValid(q, a) {
			s.Index, s.initialized = i, false
			s.current(q)
			s.Error, s.answerRequired = locales.ComponentQuestionAnswerRequired(), true
			return
		}
		result[q.Name] = a
	}
	if err := o.OnSubmit(result); err != nil {
		s.answerRequired = false
		s.Error = err.Error()
		return
	}
	s.Error, s.Submitted = "", true
	s.answerRequired = false
}

func (s *QuestionnaireState) edited() {
	s.Error, s.Submitted, s.skipped = "", false, false
	s.answerRequired = false
}

func (s *QuestionnaireState) chose(q Question) {
	if q.TextExclusive && !q.Multiple {
		s.freeform = ""
	}
	s.edited()
}

func (s *QuestionnaireState) typed(q Question) {
	if q.TextExclusive && !q.Multiple && s.freeform != "" {
		s.single = ""
	}
	s.edited()
}

func Questionnaire(c *ui.Context, s *QuestionnaireState, questions []Question, o QuestionnaireOptions) ui.Element {
	colors := theme.Colors(c)
	outer := ui.Column(c).FillWidth().MinWidth(0).Gap(theme.Space16).Role(ui.RoleGroup)
	if o.Label != "" {
		outer.Label(o.Label)
	}
	outer.Children(func() {
		if len(questions) == 0 {
			ui.Text(c, locales.ComponentNoQuestions()).TextColor(colors.MutedForeground)
			return
		}
		s.Index = max(0, min(s.Index, len(questions)-1))
		q := questions[s.Index]
		s.current(q)
		ui.Text(c, locales.ComponentQuestionOf(s.Index+1, len(questions))).
			FontSize(theme.Scale(c, 12)).TextColor(colors.MutedForeground)
		ui.Progress(c, float64(s.Index+1)/float64(len(questions))).FillWidth().Label(locales.ComponentQuestionnaireProgress())
		body := func() {
			if q.Header != "" {
				ui.Text(c, q.Header).FontSize(theme.Scale(c, 12)).TextColor(colors.MutedForeground)
			}
			ui.Text(c, q.Title).FontSize(theme.Scale(c, 18)).FontWeight(500).Role(ui.RoleHeading).Level(2)
			if q.Description != "" {
				ui.Text(c, q.Description).FontSize(theme.Scale(c, 14)).TextColor(colors.MutedForeground)
			}
			choices := func() {
				for i, choice := range q.Choices {
					if q.Multiple {
						value := s.multiple[choice.Value]
						if value == nil {
							value = new(bool)
							s.multiple[choice.Value] = value
						}
						form.Checkbox(c.Key(string(choice.Value)), value, form.CheckboxOptions{
							ControlOptions: form.ControlOptions{Label: choice.Label, Description: choice.Description,
								Disabled: o.Busy || o.Disabled || choice.Disabled},
							OnChange: func(bool) { s.edited() },
						})
					} else {
						selected := s.single == choice.Value
						control := ui.RadioBase(c.Key(string(choice.Value)), &s.single, choice.Value).
							FillWidth().MinHeight(theme.Metrics(c).ControlHeight).Padding(10, 12).Gap(theme.Space8).
							Radius(theme.RadiusPanel).Border(theme.HairlineWidth, colors.Input).
							Justify(ui.Start).AlignItems(ui.Start).
							Disabled(o.Busy || o.Disabled || choice.Disabled).Label(choice.Label).Description(choice.Description).
							OnChange(func() { s.chose(q) })
						if selected {
							control.Background(colors.Accent).BorderColor(colors.Ring)
						}
						control.Children(func() {
							ui.Column(c).FillWidth().MinWidth(0).Gap(theme.Space4).Children(func() {
								ui.Text(c, choice.Label).FontSize(theme.Scale(c, 14)).FontWeight(500)
								if choice.Description != "" {
									form.Description(c, choice.Description)
								}
							})
						})
					}
					if choice.Hint != "" {
						ui.Text(c, choice.Hint).FontSize(theme.Scale(c, 12)).TextColor(colors.MutedForeground)
					}
					if choice.Preview != "" {
						if s.previews == nil {
							s.previews = make(map[ChoiceID]bool)
						}
						mygoui.Button(c.Key("preview/"+string(choice.Value)), mygoui.ButtonOptions{
							Label:   locales.Catalog.Format(i18n.Locale(), "componentQuestionPreview", map[string]any{"option": choice.Label}),
							Variant: mygoui.ButtonGhost, Size: mygoui.ButtonSM, Wrap: true,
							OnClick: func() { s.previews[choice.Value] = !s.previews[choice.Value] },
						}).Expanded(s.previews[choice.Value]).FillWidth().MinWidth(0).Justify(ui.Start)
						if s.previews[choice.Value] {
							ui.Scroll(c.Key("preview-body/" + string(choice.Value))).FillWidth().MinWidth(0).
								MaxHeight(theme.Scale(c, 144)).Children(func() {
								ui.Text(c, choice.Preview).Selectable().FontSize(theme.Scale(c, 13))
							})
						}
					}
					if o.LetterShortcuts && i < 26 && !choice.Disabled && !o.Busy && !o.Disabled {
						key := ui.Key(int(ui.KeyA) + i)
						s.questionFocus.OnShortcut(c, 0, key, func() {
							if q.Multiple {
								value := s.multiple[choice.Value]
								*value = !*value
							} else {
								s.single = choice.Value
							}
							s.chose(q)
						})
					}
				}
			}
			var group ui.Element
			if q.Multiple {
				group = ui.Column(c.Key(string(q.Name))).Role(ui.RoleGroup).Children(choices)
			} else {
				group = ui.RadioGroup(c.Key(string(q.Name)), choices)
			}
			group.FillWidth().MinWidth(0).Gap(theme.Space8).Bind(&s.questionFocus).Label(q.Title)
			if q.AllowText {
				label := q.TextLabel
				if label == "" {
					label = locales.ComponentAnotherAnswer()
				}
				form.Textarea(c.Key("freeform"), &s.freeform, form.TextareaOptions{
					InputOptions: form.InputOptions{ControlOptions: form.ControlOptions{
						Label: label, Disabled: o.Busy || o.Disabled, Handle: &s.focus}, Placeholder: q.TextPlaceholder,
						OnChange: func(string) { s.typed(q) }},
					MinLines: 2, MaxLines: 5,
				})
			}
		}
		if o.MaxBodyHeight > 0 {
			ui.Scroll(c.Key("question-body/" + string(q.Name))).FillWidth().MinWidth(0).MinHeight(0).
				Label(locales.Catalog.Format(i18n.Locale(), "componentQuestionBody", nil)).
				MaxHeight(o.MaxBodyHeight).Children(func() {
				ui.Column(c).FillWidth().MinWidth(0).Gap(theme.Space12).Children(body)
			})
		} else {
			body()
		}
		ui.Column(c.Key("validation")).FillWidth().MinHeight(theme.Scale(c, 20)).Children(func() {
			if s.answerRequired {
				s.Error = locales.ComponentQuestionAnswerRequired()
			}
			if s.Error != "" {
				form.Error(c, s.Error)
			}
		})
		ui.Row(c).FillWidth().Wrap().Gap(theme.Space8).Justify(ui.End).Children(func() {
			if o.OnCancel != nil {
				mygoui.Button(c.Key("cancel"), mygoui.ButtonOptions{Label: locales.ComponentCancel(), Variant: mygoui.ButtonGhost,
					Disabled: o.Busy && !o.KeepCancelEnabled, OnClick: o.OnCancel})
			}
			mygoui.Button(c.Key("previous"), mygoui.ButtonOptions{Label: locales.ComponentPrevious(), Variant: mygoui.ButtonGhost,
				Disabled: o.Busy || o.Disabled || s.Index == 0, OnClick: func() { s.Previous(questions, o) }})
			if q.Optional {
				mygoui.Button(c.Key("skip"), mygoui.ButtonOptions{Label: locales.ComponentSkip(), Variant: mygoui.ButtonGhost,
					Disabled: o.Busy || o.Disabled, OnClick: func() { s.Skip(questions, o) }})
			}
			if s.Index < len(questions)-1 {
				mygoui.Button(c.Key("next"), mygoui.ButtonOptions{Label: locales.ComponentNext(), Disabled: o.Busy || o.Disabled,
					OnClick: func() { s.Next(questions, o) }})
				return
			}
			label := o.SubmitLabel
			if label == "" {
				label = locales.ComponentSubmitAnswers()
			}
			mygoui.Button(c.Key("submit"), mygoui.ButtonOptions{Label: label, Loading: o.Busy,
				Disabled: o.Busy || o.Disabled || s.Submitted || o.OnSubmit == nil, OnClick: func() { s.submit(questions, o) }})
		})
	})
	return outer
}
