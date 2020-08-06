package main

import (
	"testing"

	"github.com/maxatome/go-testdeep"
)

func TestCreateCertRewriteRules(t *testing.T) {
	td := testdeep.NewT(t)

	// nil
	rules, err := createRewriteRules(nil)
	td.Nil(rules)
	td.CmpNoError(err)

	// one rule
	rules, err = createRewriteRules([]string{"[a]|b"})
	td.Cmp(len(rules), 1)
	td.CmpNoError(err)
	td.Cmp(rules[0]("asd"), "bsd")

	// two rules
	rules, err = createRewriteRules([]string{"[a]|b", "[s]|c"})
	td.Cmp(len(rules), 2)
	td.CmpNoError(err)
	td.Cmp(rules[0]("asd"), "bsd")
	td.Cmp(rules[1]("asd"), "acd")

	// error parse
	rules, err = createRewriteRules([]string{""})
	td.CmpError(err)
	rules, err = createRewriteRules([]string{"asd"})
	td.CmpError(err)
	rules, err = createRewriteRules([]string{"asd|asd|asd"})
	td.CmpError(err)

	// error compile
	rules, err = createRewriteRules([]string{"[asd|asd"})
	td.CmpError(err)
}
