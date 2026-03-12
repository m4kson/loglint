package detector

import (
	"go/ast"
	"go/token"
	"go/types"

	"golang.org/x/tools/go/analysis"
)

type LoggerKind string

const (
	LoggerKindSlog   LoggerKind = "slog"
	LoggerKindZap    LoggerKind = "zap"
	LoggerKindStdlib LoggerKind = "stdlib"
)

type Call struct {
	Expr     *ast.CallExpr
	MsgLit   *ast.BasicLit
	MsgValue string
	MsgPos   token.Pos
	Kind     LoggerKind
}

type loggerMethod struct {
	importPath  string
	methodName  string
	msgArgIndex int
	kind        LoggerKind
}

var knownMethods = []loggerMethod{
	// ── log/slog ────────────────────────────────────────────────────────────
	{importPath: "log/slog", methodName: "Debug", msgArgIndex: 0, kind: LoggerKindSlog},
	{importPath: "log/slog", methodName: "Info", msgArgIndex: 0, kind: LoggerKindSlog},
	{importPath: "log/slog", methodName: "Warn", msgArgIndex: 0, kind: LoggerKindSlog},
	{importPath: "log/slog", methodName: "Error", msgArgIndex: 0, kind: LoggerKindSlog},
	{importPath: "log/slog", methodName: "DebugContext", msgArgIndex: 1, kind: LoggerKindSlog},
	{importPath: "log/slog", methodName: "InfoContext", msgArgIndex: 1, kind: LoggerKindSlog},
	{importPath: "log/slog", methodName: "WarnContext", msgArgIndex: 1, kind: LoggerKindSlog},
	{importPath: "log/slog", methodName: "ErrorContext", msgArgIndex: 1, kind: LoggerKindSlog},

	// ── go.uber.org/zap ─────────────────────────────────────────────────────
	{importPath: "go.uber.org/zap", methodName: "Debug", msgArgIndex: 0, kind: LoggerKindZap},
	{importPath: "go.uber.org/zap", methodName: "Info", msgArgIndex: 0, kind: LoggerKindZap},
	{importPath: "go.uber.org/zap", methodName: "Warn", msgArgIndex: 0, kind: LoggerKindZap},
	{importPath: "go.uber.org/zap", methodName: "Error", msgArgIndex: 0, kind: LoggerKindZap},
	{importPath: "go.uber.org/zap", methodName: "Fatal", msgArgIndex: 0, kind: LoggerKindZap},
	{importPath: "go.uber.org/zap", methodName: "Panic", msgArgIndex: 0, kind: LoggerKindZap},
	{importPath: "go.uber.org/zap", methodName: "DPanic", msgArgIndex: 0, kind: LoggerKindZap},
	// zap sugared logger — format methods (msg is the format string)
	{importPath: "go.uber.org/zap", methodName: "Debugf", msgArgIndex: 0, kind: LoggerKindZap},
	{importPath: "go.uber.org/zap", methodName: "Infof", msgArgIndex: 0, kind: LoggerKindZap},
	{importPath: "go.uber.org/zap", methodName: "Warnf", msgArgIndex: 0, kind: LoggerKindZap},
	{importPath: "go.uber.org/zap", methodName: "Errorf", msgArgIndex: 0, kind: LoggerKindZap},
	{importPath: "go.uber.org/zap", methodName: "Fatalf", msgArgIndex: 0, kind: LoggerKindZap},
	{importPath: "go.uber.org/zap", methodName: "Debugw", msgArgIndex: 0, kind: LoggerKindZap},
	{importPath: "go.uber.org/zap", methodName: "Infow", msgArgIndex: 0, kind: LoggerKindZap},
	{importPath: "go.uber.org/zap", methodName: "Warnw", msgArgIndex: 0, kind: LoggerKindZap},
	{importPath: "go.uber.org/zap", methodName: "Errorw", msgArgIndex: 0, kind: LoggerKindZap},
	{importPath: "go.uber.org/zap", methodName: "Fatalw", msgArgIndex: 0, kind: LoggerKindZap},

	// ── log (stdlib) ─────────────────────────────────────────────────────────
	{importPath: "log", methodName: "Print", msgArgIndex: 0, kind: LoggerKindStdlib},
	{importPath: "log", methodName: "Println", msgArgIndex: 0, kind: LoggerKindStdlib},
	{importPath: "log", methodName: "Printf", msgArgIndex: 0, kind: LoggerKindStdlib},
	{importPath: "log", methodName: "Fatal", msgArgIndex: 0, kind: LoggerKindStdlib},
	{importPath: "log", methodName: "Fatalf", msgArgIndex: 0, kind: LoggerKindStdlib},
	{importPath: "log", methodName: "Fatalln", msgArgIndex: 0, kind: LoggerKindStdlib},
	{importPath: "log", methodName: "Panic", msgArgIndex: 0, kind: LoggerKindStdlib},
	{importPath: "log", methodName: "Panicf", msgArgIndex: 0, kind: LoggerKindStdlib},
	{importPath: "log", methodName: "Panicln", msgArgIndex: 0, kind: LoggerKindStdlib},
}

func Detect(pass *analysis.Pass) []Call {
	index := buildIndex()

	var calls []Call

	for _, file := range pass.Files {
		ast.Inspect(file, func(n ast.Node) bool {
			callExpr, ok := n.(*ast.CallExpr)
			if !ok {
				return true
			}

			method, matched := resolveMethod(pass, callExpr, index)
			if !matched {
				return true
			}

			call, ok := extractCall(callExpr, method)
			if !ok {
				return true
			}

			calls = append(calls, call)
			return true
		})
	}

	return calls
}

type methodKey struct {
	importPath string
	name       string
}

func buildIndex() map[methodKey]loggerMethod {
	idx := make(map[methodKey]loggerMethod, len(knownMethods))
	for _, m := range knownMethods {
		idx[methodKey{importPath: m.importPath, name: m.methodName}] = m
	}
	return idx
}

func resolveMethod(
	pass *analysis.Pass,
	callExpr *ast.CallExpr,
	index map[methodKey]loggerMethod,
) (loggerMethod, bool) {
	switch fn := callExpr.Fun.(type) {
	case *ast.SelectorExpr:
		return resolveSelector(pass, fn, index)

	case *ast.Ident:
		return resolveIdent(pass, fn, index)
	}

	return loggerMethod{}, false
}

func resolveSelector(
	pass *analysis.Pass,
	sel *ast.SelectorExpr,
	index map[methodKey]loggerMethod,
) (loggerMethod, bool) {
	obj := pass.TypesInfo.ObjectOf(sel.Sel)
	if obj == nil {
		return loggerMethod{}, false
	}

	pkg := obj.Pkg()
	if pkg == nil {
		return loggerMethod{}, false
	}

	key := methodKey{importPath: pkg.Path(), name: sel.Sel.Name}
	m, ok := index[key]
	return m, ok
}

func resolveIdent(
	pass *analysis.Pass,
	ident *ast.Ident,
	index map[methodKey]loggerMethod,
) (loggerMethod, bool) {
	obj := pass.TypesInfo.ObjectOf(ident)
	if obj == nil {
		return loggerMethod{}, false
	}

	fn, ok := obj.(*types.Func)
	if !ok {
		return loggerMethod{}, false
	}

	pkg := fn.Pkg()
	if pkg == nil {
		return loggerMethod{}, false
	}

	key := methodKey{importPath: pkg.Path(), name: ident.Name}
	m, ok := index[key]
	return m, ok
}

func extractCall(callExpr *ast.CallExpr, method loggerMethod) (Call, bool) {
	if len(callExpr.Args) <= method.msgArgIndex {
		return Call{}, false
	}

	arg := callExpr.Args[method.msgArgIndex]

	lit, ok := arg.(*ast.BasicLit)
	if !ok {
		return Call{}, false
	}

	if lit.Kind != token.STRING {
		return Call{}, false
	}

	raw := lit.Value
	if len(raw) < 2 {
		return Call{}, false
	}
	unquoted := raw[1 : len(raw)-1]

	return Call{
		Expr:     callExpr,
		MsgLit:   lit,
		MsgValue: unquoted,
		MsgPos:   lit.Pos(),
		Kind:     method.kind,
	}, true
}
