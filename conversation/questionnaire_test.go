package conversation

import (
	"slices"
	"testing"

	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

func TestQuestionnairePreservesAnswersAndRejectsChangedChoices(t *testing.T) {
	questions := []Question{
		{Name: "direction", Title: "Direction", Choices: []Choice{{Value: "native", Label: "Native"}}},
		{Name: "details", Title: "Details", Multiple: true, Choices: []Choice{
			{Value: "source", Label: "Source"}, {Value: "status", Label: "Status"},
		}},
	}
	var state QuestionnaireState
	submits := 0
	o := QuestionnaireOptions{OnSubmit: func(answers map[QuestionID]Answer) error {
		submits++
		if !slices.Equal(answers["details"].Values, []ChoiceID{"source", "status"}) {
			t.Fatalf("multiple answers lost across navigation: %+v", answers)
		}
		a := answers["details"]
		a.Values[0] = "consumer edit"
		delete(answers, "direction")
		return nil
	}}
	tt := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		Questionnaire(c, &state, questions, o)
	}, 600, 700)
	for _, label := range []string{"Next", "Native", "Next", "Source", "Status", "Previous", "Next"} {
		if err := tt.Click(label); err != nil {
			t.Fatal(err)
		}
	}
	if state.Index != 1 || !slices.Equal(state.Answers["details"].Values, []ChoiceID{"source", "status"}) {
		t.Fatalf("previous discarded partial multiple answer: %+v", state.Answers)
	}
	questions[0].Choices[0].Disabled = true
	tt.Frame()
	if err := tt.Click("Submit answers"); err != nil {
		t.Fatal(err)
	}
	if submits != 0 || state.Index != 0 || state.Error == "" {
		t.Fatalf("disabled old answer was submitted: count=%d index=%d error=%q", submits, state.Index, state.Error)
	}
	questions[0].Choices[0].Disabled = false
	tt.Frame()
	if err := tt.Click("Next"); err != nil {
		t.Fatal(err)
	}
	if err := tt.Click("Submit answers"); err != nil {
		t.Fatal(err)
	}
	if submits != 1 || !state.Submitted || state.Answers["details"].Values[0] != "source" {
		t.Fatalf("consumer mutation escaped submitted snapshot: %+v, submits=%d", state.Answers, submits)
	}
}

func TestQuestionnaireFreeformDoesNotTriggerChoiceShortcut(t *testing.T) {
	var state QuestionnaireState
	questions := []Question{{Name: "answer", Title: "Answer", AllowText: true, TextLabel: "Written answer",
		Choices: []Choice{{Value: "a", Label: "First choice"}, {Value: "b", Label: "Second choice"}}}}
	var got Answer
	tt := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		Questionnaire(c, &state, questions, QuestionnaireOptions{LetterShortcuts: true,
			OnSubmit: func(answers map[QuestionID]Answer) error { got = answers["answer"]; return nil }})
	}, 600, 600)
	if err := tt.Click("Written answer"); err != nil {
		t.Fatal(err)
	}
	tt.Type("b")
	if err := tt.Click("Submit answers"); err != nil {
		t.Fatal(err)
	}
	if got.Text != "b" || len(got.Values) != 0 {
		t.Fatalf("typing in freeform activated a choice: %+v", got)
	}
}

func TestQuestionnaireFinalSkipAndBusyNavigation(t *testing.T) {
	state := QuestionnaireState{}
	questions := []Question{{Name: "optional", Title: "Optional details", Optional: true,
		Choices: []Choice{{Value: "x", Label: "Extra"}}}}
	var result Answer
	options := QuestionnaireOptions{OnSubmit: func(answers map[QuestionID]Answer) error {
		result = answers["optional"]
		return nil
	}}
	tt := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		Questionnaire(c, &state, questions, options)
	}, 500, 500)
	if err := tt.Click("Extra"); err != nil {
		t.Fatal(err)
	}
	if err := tt.Click("Skip"); err != nil {
		t.Fatal(err)
	}
	options.Busy = true
	tt.Frame()
	if state.Next(questions, options) || state.Skip(questions, options) {
		t.Fatal("busy questionnaire accepted a navigation mutation")
	}
	options.Busy = false
	tt.Frame()
	if err := tt.Click("Submit answers"); err != nil {
		t.Fatal(err)
	}
	if !result.Skipped || len(result.Values) != 0 || result.Text != "" {
		t.Fatalf("last optional skip revived its earlier answer: %+v", result)
	}
}

func TestQuestionnaireTextExclusiveIsOptInAndMultiKeepsBoth(t *testing.T) {
	for _, test := range []struct {
		name                string
		exclusive, multiple bool
	}{
		{name: "existing single behavior"},
		{name: "exclusive single", exclusive: true},
		{name: "exclusive multi", exclusive: true, multiple: true},
	} {
		t.Run(test.name, func(t *testing.T) {
			var state QuestionnaireState
			questions := []Question{{Name: "q", Title: "Question", AllowText: true,
				TextExclusive: test.exclusive, Multiple: test.multiple,
				Choices: []Choice{{Value: "a", Label: "First"}}, TextLabel: "Written answer"}}
			var got Answer
			tt := ui.NewTester(func(c *ui.Context) {
				theme.Apply(c)
				Questionnaire(c, &state, questions, QuestionnaireOptions{
					OnSubmit: func(answers map[QuestionID]Answer) error { got = answers["q"]; return nil },
				})
			}, 600, 700)
			if err := tt.Click("First"); err != nil {
				t.Fatal(err)
			}
			if err := tt.Click("Written answer"); err != nil {
				t.Fatal(err)
			}
			tt.Type("custom")
			if err := tt.Click("Submit answers"); err != nil {
				t.Fatal(err)
			}
			wantValues := 1
			if test.exclusive && !test.multiple {
				wantValues = 0
			}
			if got.Text != "custom" || len(got.Values) != wantValues {
				t.Fatalf("choice plus text semantics changed: %+v", got)
			}
		})
	}
}
