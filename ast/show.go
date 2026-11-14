package ast

func Show(n any) string {
	switch x := n.(type) {
	case *IntLit:
		return ShowInt(x.Value)
	case *StrLit:
		return ShowStr(x.Value)
	case *Var:
		return ShowVar(x.Name)
	case *Con:
		return ShowCon(x.Name)
	case *Tuple:
		if len(x.Elems) == 0 {
			return "()"
		}
		s := "(tuple "
		for i, e := range x.Elems {
			if i > 0 {
				s += " "
			}
			s += Show(e)
		}
		return s + ")"
	case *List:
		if len(x.Elems) == 0 {
			return "(list)"
		}
		s := "(list "
		for i, e := range x.Elems {
			if i > 0 {
				s += " "
			}
			s += Show(e)
		}
		return s + ")"
	case *App:
		return "(app " + Show(x.Fn) + " " + Show(x.Arg) + ")"
	case *If:
		return "(if " + Show(x.Cond) + " " + Show(x.Then) + " " + Show(x.Else) + ")"
	case *Fn:
		arms := ""
		for i, arm := range x.Arms {
			if i > 0 {
				arms += " "
			}
			arms += "(" + Show(arm.Pat) + " " + Show(arm.Body) + ")"
		}
		return "(fn " + arms + ")"
	case *Let:
		decs := ""
		for i, d := range x.Decs {
			if i > 0 {
				decs += " "
			}
			switch v := d.(type) {
			case *ValDec:
				pat := Show(v.Pat)
				if pat == "_" {
					decs += Show(v.Expr)
				} else {
					decs += "(" + pat + " " + Show(v.Expr) + ")"
				}
			default:
				decs += Show(d)
			}
		}
		body := ""
		for i, e := range x.Body {
			if i > 0 {
				body += " "
			}
			body += Show(e)
		}
		return "(let " + decs + " " + body + ")"
	case *Seq:
		exprs := ""
		for i, e := range x.Exprs {
			if i > 0 {
				exprs += " "
			}
			exprs += Show(e)
		}
		return "(seq " + exprs + ")"
	case *AndAlso:
		return "(andalso " + Show(x.L) + " " + Show(x.R) + ")"
	case *OrElse:
		return "(orelse " + Show(x.L) + " " + Show(x.R) + ")"
	case *PWild:
		return "_"
	case *PVar:
		return ShowVar(x.Name)
	case *PTuple:
		if len(x.Elems) == 0 {
			return "()"
		}
		s := "(tuple "
		for i, p := range x.Elems {
			if i > 0 {
				s += " "
			}
			s += Show(p)
		}
		return s + ")"
	case *ValDec:
		return "(val " + Show(x.Pat) + " " + Show(x.Expr) + ")"
	case *FunDec:
		clauses := ""
		for i, c := range x.Clauses {
			if i > 0 {
				clauses += " "
			}
			pats := ""
			for j, p := range c.Pats {
				if j > 0 {
					pats += " "
				}
				pats += Show(p)
			}
			clauses += "(" + pats + " " + Show(c.Body) + ")"
		}
		return "(fun " + x.Name + " (" + clauses + "))"
	default:
		return "<?>"
	}
}

func ShowInt(v int64) string {
	if v < 0 {
		return "-" + ShowInt(-v)
	}
	s := []byte{}
	for v > 0 {
		s = append(s, byte('0'+v%10))
		v /= 10
	}
	if len(s) == 0 {
		return "0"
	}
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
	return string(s)
}

func ShowStr(s string) string {
	r := []byte{}
	for i := 0; i < len(s); i++ {
		c := s[i]
		switch c {
		case '"':
			r = append(r, '\\')
			fallthrough
		case '\\':
			r = append(r, '\\')
		default:
			r = append(r, c)
		}
	}
	return "\"" + string(r) + "\""
}

func ShowVar(name string) string {
	return name
}

func ShowCon(name string) string {
	return name
}
