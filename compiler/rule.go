package compiler

type ParseRule struct {
	preffix    func()
	infix      func()
	precedence int
}

func (p *Parser) initRules() []ParseRule {
	return []ParseRule{
		LParentheses: {p.grouping, nil, PrecNone},
		RParentheses: {nil, nil, PrecNone},
		LBrace:       {nil, nil, PrecNone},
		RBrace:       {nil, nil, PrecNone},
		Comma:        {nil, nil, PrecNone},
		Dot:          {nil, nil, PrecNone},
		Minus:        {p.unary, p.grouping, PrecTerm},
		Plus:         {nil, p.binary, PrecTerm},
		Semicolon:    {nil, nil, PrecNone},
		Slash:        {nil, p.binary, PrecFactor},
		Star:         {nil, p.binary, PrecFactor},
		Not:          {nil, nil, PrecNone},
		NotEqual:     {nil, nil, PrecNone},
		Equal:        {nil, nil, PrecNone},
		EqualEqual:   {nil, nil, PrecNone},
		Greater:      {nil, nil, PrecNone},
		GreaterEqual: {nil, nil, PrecNone},
		Less:         {nil, nil, PrecNone},
		LessEqual:    {nil, nil, PrecNone},
		Identifier:   {nil, nil, PrecNone},
		String:       {nil, nil, PrecNone},
		Number:       {p.number, nil, PrecNone},
		And:          {nil, nil, PrecNone},
		Class:        {nil, nil, PrecNone},
		Else:         {nil, nil, PrecNone},
		False:        {nil, nil, PrecNone},
		For:          {nil, nil, PrecNone},
		Function:     {nil, nil, PrecNone},
		If:           {nil, nil, PrecNone},
		Null:         {nil, nil, PrecNone},
		Or:           {nil, nil, PrecNone},
		Print:        {nil, nil, PrecNone},
		Return:       {nil, nil, PrecNone},
		Super:        {nil, nil, PrecNone},
		This:         {nil, nil, PrecNone},
		True:         {nil, nil, PrecNone},
		Var:          {nil, nil, PrecNone},
		While:        {nil, nil, PrecNone},
		Error:        {nil, nil, PrecNone},
		Eof:          {nil, nil, PrecNone},
	}
}

func (p *Parser) getRule(tokenType int) *ParseRule {
	return &p.rules[tokenType]
}
