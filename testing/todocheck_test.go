package testing

import (
	"testing"

	"github.com/preslavmihaylov/todocheck/testing/scenariobuilder"
)

func TestSingleLineMalformedTodos(t *testing.T) {
	err := scenariobuilder.NewScenario().
		WithBinary("../todocheck").
		WithBasepath("./scenarios/singleline_todos").
		WithConfig("./scenarios/configs/no_issue_tracker.yaml").
		ExpectTodoErr(
			scenariobuilder.NewTodoErr().
				WithType(scenariobuilder.TodoErrTypeMalformed).
				WithLocation("scenarios/singleline_todos/other.go", 3).
				ExpectLine("// TODO: This is a todo in a separate go file")).
		ExpectTodoErr(
			scenariobuilder.NewTodoErr().
				WithType(scenariobuilder.TodoErrTypeMalformed).
				WithLocation("scenarios/singleline_todos/main.go", 5).
				ExpectLine("// TODO: This is a malformed todo")).
		ExpectTodoErr(
			scenariobuilder.NewTodoErr().
				WithType(scenariobuilder.TodoErrTypeMalformed).
				WithLocation("scenarios/singleline_todos/main.go", 6).
				ExpectLine("// TODO: This is a malformed todo 2")).
		ExpectTodoErr(
			scenariobuilder.NewTodoErr().
				WithType(scenariobuilder.TodoErrTypeMalformed).
				WithLocation("scenarios/singleline_todos/main.go", 10).
				ExpectLine("func main() { // TODO: This is a todo comment at the end of a line")).
		ExpectTodoErr(
			scenariobuilder.NewTodoErr().
				WithType(scenariobuilder.TodoErrTypeMalformed).
				WithLocation("scenarios/singleline_todos/main.go", 15).
				ExpectLine("// TODO comment without colons")).
		ExpectTodoErr(
			scenariobuilder.NewTodoErr().
				WithType(scenariobuilder.TodoErrTypeMalformed).
				WithLocation("scenarios/singleline_todos/main.go", 17).
				ExpectLine("// This is a TODO comment at the middle of it")).
		Run()
	if err != nil {
		t.Errorf("%s", err)
	}
}

func TestMultiLineMalformedTodos(t *testing.T) {
	err := scenariobuilder.NewScenario().
		WithBinary("../todocheck").
		WithBasepath("./scenarios/multiline_todos").
		WithConfig("./scenarios/configs/no_issue_tracker.yaml").
		ExpectTodoErr(
			scenariobuilder.NewTodoErr().
				WithType(scenariobuilder.TodoErrTypeMalformed).
				WithLocation("scenarios/multiline_todos/main.go", 3).
				ExpectLine("/*").
				ExpectLine(" * This is a multiline").
				ExpectLine(" * TODO comment.").
				ExpectLine(" */")).
		ExpectTodoErr(
			scenariobuilder.NewTodoErr().
				WithType(scenariobuilder.TodoErrTypeMalformed).
				WithLocation("scenarios/multiline_todos/main.go", 8).
				ExpectLine("func main() { /*").
				ExpectLine("	 * This is a multiline TODO comment").
				ExpectLine("	 * spanning several lines").
				ExpectLine("	 */")).
		ExpectTodoErr(
			scenariobuilder.NewTodoErr().
				WithType(scenariobuilder.TodoErrTypeMalformed).
				WithLocation("scenarios/multiline_todos/main.go", 18).
				ExpectLine("/* This is a single-line multi-line TODO comment */")).
		Run()
	if err != nil {
		t.Errorf("%s", err)
	}
}